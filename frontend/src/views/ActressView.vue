<template>
    <div class="container">
        <h1>{{ actressName }} 的作品</h1>
        
        <div class="filter-section" v-if="uniqueGenres.length > 0">
            <span 
                class="genre-filter-tag" 
                :class="{ active: selectedGenre === '' }"
                @click="selectedGenre = ''"
            >
                全部作品
            </span>
            <span 
                class="genre-filter-tag" 
                v-for="genre in uniqueGenres" 
                :key="genre"
                :class="{ active: selectedGenre === genre }"
                @click="selectedGenre = genre"
            >
                {{ genre }}
            </span>
        </div>

        <div class="video-grid">
            <VideoCard 
                v-for="video in filteredVideos" 
                :key="video.id" 
                :video="video" 
                @click="navigateToDetail(video.id)" 
            />
        </div>
        
        <div v-if="filteredVideos.length === 0" class="no-results">
            <p>没有找到相关作品</p>
        </div>
    </div>
</template>

<script>
import VideoCard from '../components/VideoCard.vue'
import videosApi from '../api/videos'

export default {
    name: 'ActressView',
    components: { VideoCard },
    props: ['actressName'],
    data() {
        return {
            videos: [],
            selectedGenre: '' // 当前选中的过滤标签，为空字符串时表示不筛选
        }
    },
    computed: {
        // 自动提取并去重该演员所有作品的标签
        uniqueGenres() {
            const genresSet = new Set();
            this.videos.forEach(video => {
                if (video.genres && video.genres.length > 0) {
                    video.genres.forEach(g => genresSet.add(g));
                }
            });
            // 转换为数组并排序，让标签展示更稳定
            return Array.from(genresSet).sort();
        },
        // 根据选中标签过滤视频列表
        filteredVideos() {
            if (!this.selectedGenre) {
                return this.videos; // 未选择标签，展示全部
            }
            return this.videos.filter(video => 
                video.genres && video.genres.includes(this.selectedGenre)
            );
        }
    },
    async created() {
        this.videos = await videosApi.getVideoListByActress(this.actressName)
    },
    methods: {
        navigateToDetail(id) {
            this.$router.push({ 
                name: 'detail', 
                params: { actressName: this.actressName, id: id } 
            })
        }
    }
}
</script>

<style scoped>
.container { 
    padding: 15px; 
}

h1 { 
    color: var(--text-color, #333); 
    margin-bottom: 1.5rem; 
    font-weight: 600; 
    position: relative; 
    display: inline-block; 
}

h1::after { 
    content: ''; 
    position: absolute; 
    bottom: -8px; 
    left: 0; 
    width: 50px; 
    height: 3px; 
    background: var(--secondary-color, #ff6b8b); 
}

/* ================= 新增：标签筛选区样式 ================= */
.filter-section {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 20px;
    align-items: center;
}

.genre-filter-tag {
    background-color: #fff0f3; 
    color: var(--secondary-color, #ff6b8b);
    border: 1px solid var(--accent-color, #ffcdd8);
    padding: 5px 12px;
    border-radius: 15px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    user-select: none; /* 防止频繁点击时选中文本 */
}

.genre-filter-tag:hover {
    background-color: #ffe0e8;
}

/* 激活(被选中)时的样式 */
.genre-filter-tag.active {
    background-color: var(--secondary-color, #ff6b8b);
    color: white;
    border-color: var(--secondary-color, #ff6b8b);
    font-weight: bold;
}

.no-results {
    text-align: center;
    color: #999;
    padding: 40px 0;
    font-size: 1.1rem;
}
/* ===================================================== */

.video-grid {
    display: grid;
    grid-template-columns: 1fr; 
    gap: 20px; 
    padding: 10px 0;
}

@media (min-width: 768px) {
    .container { 
        padding: 2rem; 
    }
    .video-grid {
        grid-template-columns: repeat(4, 1fr);
        gap: 25px; 
    }
}
</style>