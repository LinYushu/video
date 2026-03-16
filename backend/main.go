package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	// "strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 服务器配置
const (
	basePath   = "F://sVideo"
	serverPort = ":31471"
	apiKey     = "IBHUSDBWQHJEJOBDSW"
)

// 全局缓存
var (
	actressListCache []VideoItem
	videoListCache map[string][]VideoItem
	genreListCache     []GenreItem
	genreVideoMapCache map[string][]VideoItem
	cacheMutex     sync.RWMutex
	logger         = log.New(os.Stdout, "[MissAV] ", log.LstdFlags|log.Lshortfile)
)

// VideoItem 表示视频列表项
type VideoItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Poster string `json:"poster"`
	Fanart string `json:"fanart"`
	Actress string `json:"actress"`
	Genres  []string `json:"genres"`
}

// VideoDetail 视频详细信息
type VideoDetail struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"releaseDate"`
	Fanarts     []string `json:"fanarts"`
	VideoFile   string   `json:"videoFile,omitempty"`
	Actress		string   `json:"actress"`
	Genres      []string `json:"genres"`
}

// NfoFile NFO文件结构
type NfoFile struct {
	XMLName     xml.Name `xml:"movie"`
	Title       string   `xml:"title"`
	ReleaseDate string   `xml:"releasedate"`
	Premiered   string   `xml:"premiered"`
	Genres      []string `xml:"genre"`
}

type GenreItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允许的域名（生产环境应替换为实际前端域名）
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许的HTTP方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		// 允许的请求头
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 允许携带Cookie（如果需要）
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 预检请求直接返回200
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	logger.Println("Starting MissAV server...")

	// 初始化缓存
	if err := buildCache(); err != nil {
		logger.Fatalf("Failed to build initial cache: %v", err)
	}

	// 启动定时缓存更新
	go startCacheUpdater(30 * time.Minute)

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/actresses", listActressesHandler)
	mux.HandleFunc("/api/actress/", listVideosByActressHandler)
	mux.HandleFunc("/api/videos/", videoDetailHandler)
	mux.HandleFunc("/api/addvideo/", addVideoHandler)
	mux.HandleFunc("/file/", imageHandler)
	mux.HandleFunc("/api/genres", listGenresHandler)
	mux.HandleFunc("/api/genre/", listVideosByGenreHandler)

	// 包装CORS中间件
	handler := enableCORS(mux)

	// 确保端口号是干净的（如 "8080" 而不是 ":8080"）
	port := strings.TrimPrefix(serverPort, ":")

	// 方案1：显式创建IPv6监听器（兼容IPv4）
	listener4, _ := net.Listen("tcp4", "0.0.0.0:"+port)
	listener6, _ := net.Listen("tcp6", "[::]:"+port)
	go http.Serve(listener4, handler)
	log.Fatal(http.Serve(listener6, handler))

	logger.Printf("Server started on port %s (IPv6, with IPv4 compatibility)", port)
}

// startCacheUpdater 定时更新缓存
func startCacheUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		logger.Println("Starting scheduled cache update...")
		if err := buildCache(); err != nil {
			logger.Printf("Cache update failed: %v", err)
		} else {
			logger.Println("Cache updated successfully")
		}
	}
}

func buildCache() error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	startTime := time.Now()
	logger.Println("Building cache...")

	actressDirs, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("read directory %s failed: %w", basePath, err)
	}

	newActressList := []VideoItem{}
	newVideoList := make(map[string][]VideoItem)
	newGenreVideoMap := make(map[string][]VideoItem)

	var videoCount int
	for _, actressDir := range actressDirs {
		if !actressDir.IsDir() {
			continue
		}
		actressName := actressDir.Name()

		// 排除 thumb 文件夹本身
		if strings.ToLower(actressName) == "thumb" {
			continue
		}

		actressPoster := fmt.Sprintf("/file/thumb/%s.jpg", actressName)

		videoDirs, err := os.ReadDir(filepath.Join(basePath, actressName))
		if err != nil {
			logger.Printf("Error reading actress dir %s: %v", actressName, err)
			continue
		}

		var actressVideos []VideoItem
		for _, videoDir := range videoDirs {
			if !videoDir.IsDir() {
				continue
			}
			videoID := videoDir.Name()
			videoPath := filepath.Join(basePath, actressName, videoID)
			
			// 视频封面
			posterPath := filepath.Join(videoPath, videoID+"-poster.jpg")
			if _, err := os.Stat(posterPath); err != nil {
				continue // 没有封面的视频跳过
			}
			title := videoID // 默认 fallback
			fanartUrl := ""  // 新增：用于存储找到的横版图片路径
			
			if files, err := os.ReadDir(videoPath); err == nil {
				for _, f := range files {
					if !f.IsDir() {
						name := f.Name()
						// 1. 找视频文件做标题
						if strings.Contains(name, videoID) {
							ext := strings.ToLower(filepath.Ext(name))
							if ext == ".mp4" || ext == ".mkv" || ext == ".avi" || ext == ".wmv" {
								title = strings.TrimSuffix(name, filepath.Ext(name))
							}
						}
						// 2. 找横版剧照/封面 (寻找包含 -fanart 并且是 .jpg 的文件)
						if strings.HasPrefix(name, videoID+"-fanart") && strings.HasSuffix(strings.ToLower(name), ".jpg") {
							// 只取找到的第一张图作为列表页的横版封面
							if fanartUrl == "" {
								fanartUrl = fmt.Sprintf("/file/%s/%s/%s", actressName, videoID, name)
							}
						}
					}
				}
			}

			vItem := VideoItem{
				ID:      videoID,
				Title:   title,
				Poster:  fmt.Sprintf("/file/%s/%s/%s-poster.jpg", actressName, videoID, videoID),
				Fanart:  fanartUrl,
				Actress: actressName,
			}

			nfoPath := filepath.Join(videoPath, videoID+".nfo")
			var itemGenres []string // 临时存放当前视频的标签
			
			if file, err := os.Open(nfoPath); err == nil {
				decoder := xml.NewDecoder(file)
				decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
					return input, nil
				}
				var nfo NfoFile
				if err := decoder.Decode(&nfo); err == nil {
					seenGenres := make(map[string]bool)
					for i, g := range nfo.Genres {
						g = strings.TrimSpace(g)
						// 过滤第一个且与番号相同的标签
						if i == 0 && strings.EqualFold(g, videoID) {
							continue
						}
						if g != "" && !seenGenres[g] {
							seenGenres[g] = true
							itemGenres = append(itemGenres, g)
						}
					}
				}
				file.Close()
			}
			
			// 将标签保存到当前视频项中
			vItem.Genres = itemGenres
			actressVideos = append(actressVideos, vItem)

			// 将带着标签的完整 vItem 追加到全局标签分类缓存中
			for _, g := range itemGenres {
				newGenreVideoMap[g] = append(newGenreVideoMap[g], vItem)
			}


			videoCount++
		}

		if len(actressVideos) > 0 {
			newVideoList[actressName] = actressVideos
			newActressList = append(newActressList, VideoItem{
				ID:     actressName,
				Title:  actressName,
				Poster: actressPoster,
			})
		}
	}

	newGenreList := make([]GenreItem, 0, len(newGenreVideoMap))
	for g, vids := range newGenreVideoMap {
		newGenreList = append(newGenreList, GenreItem{Name: g, Count: len(vids)})
	}
    // 降序排列
	sort.Slice(newGenreList, func(i, j int) bool {
		return newGenreList[i].Count > newGenreList[j].Count
	})

	actressListCache = newActressList
	videoListCache = newVideoList
	genreListCache = newGenreList               // 更新缓存
	genreVideoMapCache = newGenreVideoMap       // 更新缓存

	logger.Printf("Cache built successfully. Actresses: %d, Genres: %d, Videos: %d, Duration: %v",
		len(actressListCache), len(genreListCache), videoCount, time.Since(startTime))
	return nil
}

// parseTitleAndDate 解析NFO文件获取标题和日期
func parseTitleAndDate(videoID string) (title, releaseDate string, err error) {
	nfoPath := filepath.Join(basePath, videoID, videoID+".nfo")

	// 使用os.Open确保能处理BOM头
	file, err := os.Open(nfoPath)
	if err != nil {
		return "", "", fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	// 创建UTF-8解码器（自动处理BOM头）
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		// 强制使用UTF-8，忽略文件声明的编码
		return input, nil
	}

	var nfo NfoFile
	if err := decoder.Decode(&nfo); err != nil {
		return "", "", fmt.Errorf("xml decode failed: %w", err)
	}

	// 确定发布日期
	date := nfo.ReleaseDate
	if date == "" {
		date = nfo.Premiered
	}

	// 确保标题不为空
	if nfo.Title == "" {
		nfo.Title = videoID
	}

	return nfo.Title, date, nil
}

func listActressesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(actressListCache)
}

// listVideosHandler 获取视频列表
func listVideosByActressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	actressName := strings.TrimPrefix(r.URL.Path, "/api/actress/")
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	videos, ok := videoListCache[actressName]
	if !ok {
		httpError(w, "Actress not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(videos)
}

// videoDetailHandler 获取视频详情
func videoDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/videos/"), "/")
	if len(parts) != 2 {
		httpError(w, "Invalid path format", http.StatusBadRequest)
		return
	}
	actressName, videoID := parts[0], parts[1]

	detail := VideoDetail{ID: videoID, Actress: actressName}
	videoPath := filepath.Join(basePath, actressName, videoID)

	// 手动解析 NFO 获取日期和标签
	nfoPath := filepath.Join(videoPath, videoID+".nfo")
	detail.ReleaseDate = "Unknown"
	var genres []string

	if file, err := os.Open(nfoPath); err == nil {
		decoder := xml.NewDecoder(file)
		decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			return input, nil
		}
		var nfo NfoFile
		if err := decoder.Decode(&nfo); err == nil {
			// 1. 获取发行日期
			if nfo.ReleaseDate != "" {
				detail.ReleaseDate = nfo.ReleaseDate
			} else if nfo.Premiered != "" {
				detail.ReleaseDate = nfo.Premiered
			}
			
			// 2. 提取标签，去重并过滤掉番号
			seenGenres := make(map[string]bool)
			for i, g := range nfo.Genres {
				g = strings.TrimSpace(g)
				// 忽略第一个且名字与番号相同的标签
				if i == 0 && strings.EqualFold(g, videoID) {
					continue
				}
				if g != "" && !seenGenres[g] {
					seenGenres[g] = true
					genres = append(genres, g)
				}
			}
		}
		file.Close()
	}
	detail.Genres = genres // 将整理好的标签数组赋值给详情对象
	title := videoID
	actualVideoFile := ""
	if files, err := ioutil.ReadDir(videoPath); err == nil {
		for _, f := range files {
			if !f.IsDir() && strings.Contains(f.Name(), videoID) {
				ext := strings.ToLower(filepath.Ext(f.Name()))
				if ext == ".mp4" || ext == ".mkv" || ext == ".avi" || ext == ".wmv" {
					title = strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
					actualVideoFile = f.Name() // 记录带后缀的完整文件名用于播放
					break
				}
			}
		}
	}
	detail.Title = title

	if files, err := ioutil.ReadDir(videoPath); err == nil {
		var fanarts []string
		for _, file := range files {
			name := file.Name()
			if !file.IsDir() && strings.HasPrefix(name, videoID+"-fanart") && strings.HasSuffix(name, ".jpg") {
				fanarts = append(fanarts, fmt.Sprintf("/file/%s/%s/%s", actressName, videoID, name))
			}
		}
		sort.Strings(fanarts)
		detail.Fanarts = fanarts
	}

	// 使用刚刚找到的真实文件名拼接视频播放路径
	if actualVideoFile != "" {
		detail.VideoFile = fmt.Sprintf("/file/%s/%s/%s", actressName, videoID, actualVideoFile)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(detail)
}

// imageHandler 处理图片请求
func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取 /file/ 后面的相对路径
	relativePath := strings.TrimPrefix(r.URL.Path, "/file/")
	if relativePath == "" {
		httpError(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	imagePath := filepath.Join(basePath, relativePath)

	// 安全检查：防止路径穿越 (例如请求 /file/../../../etc/passwd)
	cleanBasePath := filepath.Clean(basePath)
	cleanImagePath := filepath.Clean(imagePath)
	if !strings.HasPrefix(cleanImagePath, cleanBasePath+string(os.PathSeparator)) {
		httpError(w, "Forbidden path", http.StatusForbidden)
		return
	}

	fileInfo, err := os.Stat(cleanImagePath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else if err != nil || fileInfo.IsDir() {
		httpError(w, "Internal server error or invalid file", http.StatusInternalServerError)
		return
	}

	// 设置Content-Type
	switch filepath.Ext(cleanImagePath) {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".mp4":
		w.Header().Set("Content-Type", "video/mp4")
	}

	http.ServeFile(w, r, cleanImagePath)
}

// downloadQueueHandler 新增下载队列
func addVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从请求头获取API密钥
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// 验证API密钥
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != apiKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// 读取请求体
	videoID := strings.TrimPrefix(r.URL.Path, "/api/addvideo/")
	if videoID == "" {
		httpError(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// 将请求体转换为字符串，强制大写
	id := strings.ToUpper(string(videoID))
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	logger.Printf("Received ID: %s\n", id)

	// 检查sqlite里面是否有这个车牌号
	db, err := sql.Open("sqlite3", "../db/downloaded.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	exists, err := checkStringExists(db, id)
	if err != nil {
		log.Fatal(err)
	}

	response := fmt.Sprintf("%s already downloaded", id)
	// 执行Python脚本
	if !exists {
		response = fmt.Sprintf("Add %s to download queue", id)
		go func() {
			cmd := exec.Command("sh", "-c", fmt.Sprintf("cd .. && python3 main.py %s", id))
			err := cmd.Run()
			if err != nil {
				logger.Printf("command exec failed: %v", err)
			} else {
				logger.Printf("command exec succ!")
			}
		}()
	}
	logger.Println(response)

	// 设置响应内容类型
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

// listGenresHandler 获取所有标签
func listGenresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(genreListCache)
}

// listVideosByGenreHandler 获取某个标签下的所有视频
func listVideosByGenreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	genreName := strings.TrimPrefix(r.URL.Path, "/api/genre/")
	// 处理中文字符的 URL 编码
	if decoded, err := url.QueryUnescape(genreName); err == nil {
		genreName = decoded
	}

	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	videos, ok := genreVideoMapCache[genreName]
	if !ok {
		// 即使没有找到，返回空数组
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode([]VideoItem{})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(videos)
}

func checkStringExists(db *sql.DB, target string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM MissAV WHERE bvid = ? LIMIT 1)"
	err := db.QueryRow(query, target).Scan(&exists)
	return exists, err
}

// httpError 统一的HTTP错误响应
func httpError(w http.ResponseWriter, message string, code int) {
	logger.Printf("HTTP Error %d: %s", code, message)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
