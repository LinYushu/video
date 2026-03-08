<template>
    <div class="detail-container" v-if="video">
        <div class="header">
            <h1>{{ video.title }}</h1>
            <p class="release-date">发行日期: {{ video.releaseDate }}</p>
        </div>
        
        <div class="content">
            <div class="poster">
                <img :src="video.poster" :alt="video.title">
            </div>

            <div class="video-player" v-if="video.videoFile">
                <video 
                    controls 
                    preload="metadata" 
                    controlsList="nodownload"
                    :src="video.videoFile"
                >
                    您的浏览器不支持 HTML5 视频播放。
                </video>
            </div>

            <div class="meta-section">
                <div class="actress-info" @click="goToActress">
                    <img 
                        :src="video.actressAvatar" 
                        :alt="video.actress" 
                        class="actress-avatar" 
                        @error="handleAvatarError"
                    >
                    <span class="actress-name">{{ video.actress }}</span>
                </div>
                
                <div class="genres-list" v-if="video.genres && video.genres.length">
                    <span 
                        class="genre-tag" 
                        v-for="genre in video.genres" 
                        :key="genre"
                        @click="goToGenre(genre)"
                    >
                        {{ genre }}
                    </span>
                </div>
            </div>
            <Gallery :images="video.fanarts" />
        </div>
    </div>
</template>

<script>
import Gallery from '../components/Gallery.vue'
import videosApi from '../api/videos'

export default {
    components: { Gallery },
    props: ['actressName', 'id'], 
    data() {
        return {
            video: null
        }
    },
    async created() {
        this.video = await videosApi.getVideoDetail(this.actressName, this.id)
    },
    methods: {
        // 新增：跳转到演员主页
        goToActress() {
            this.$router.push({ 
                name: 'actress', 
                params: { actressName: this.video.actress } 
            })
        },
        // 新增：跳转到标签分类页
        goToGenre(genreName) {
            this.$router.push({ 
                name: 'genre', 
                params: { genreName: genreName } 
            })
        },
        // 新增：如果某个演员没有头像图片，隐藏图片破损的 icon
        handleAvatarError(e) {
            e.target.style.display = 'none';
        }
    }
}
</script>

<style scoped>
.detail-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  background: white;
  border-radius: 16px;
  box-shadow: 0 5px 30px rgba(255, 107, 139, 0.1);
}

.header h1 {
  color: var(--text-color, #333);
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.release-date {
  color: #ff6b8b;
  font-size: 1rem;
  margin-bottom: 2rem;
}

.poster img {
  width: 100%;
  max-height: 500px;
  object-fit: contain;
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
}

.video-player {
  margin-top: 2rem;
  margin-bottom: 2rem;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 5px 25px rgba(0, 0, 0, 0.1);
  background-color: #000;
  display: flex;
  justify-content: center;
}

.video-player video {
  width: 100%;
  max-height: 75vh;
  display: block;
  outline: none;
}

/* ================= 新增：元数据展示区样式 ================= */
/* ================= 修改：元数据展示区样式 (极简 & 一行同行展示) ================= */
.meta-section {
    margin: 1.5rem 0;
    display: flex;
    flex-direction: row; /* 水平排布，不换行 */
    align-items: center; /* 垂直居中对齐 */
    flex-wrap: wrap; /* 如果屏幕太窄标签太多，允许标签自然折行 */
    gap: 15px; /* 头像名字组合 与 标签列表 之间的间距 */
    padding: 10px 0; /* 去除了繁重的背景和阴影，保留极简风格 */
}

.actress-info {
    display: flex;
    align-items: center;
    gap: 10px; /* 头像与名字之间的距离 */
    cursor: pointer;
    transition: opacity 0.2s ease;
}

.actress-info:hover {
    opacity: 0.7; /* 简单的悬停透明度反馈，去掉复杂的阴影放大效果 */
}

.actress-avatar {
    width: 55px; /* 头像缩小 */
    height: 55px;
    border-radius: 50%;
    object-fit: cover;
    border: 1px solid var(--secondary-color, #ff6b8b); /* 边框变细 */
}

.actress-name {
    font-size: 1.1rem;
    font-weight: bold;
    color: var(--text-color, #5a3a4a);
    white-space: nowrap; /* 保证名字不折行 */
}

.genres-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px; /* 【修改点】减小标签与标签之间的空隙 */
    align-items: center;
}

.genre-tag {
    /* 【修改点】极简风格的标签：浅色背景、小字体、小内边距 */
    background-color: #fff0f3; 
    color: var(--secondary-color, #ff6b8b);
    border: 1px solid var(--accent-color, #ffcdd8);
    padding: 4px 10px; /* 减小标签占用的空间 */
    border-radius: 12px; /* 圆角变小 */
    font-size: 0.85rem; /* 字体变小 */
    cursor: pointer;
    transition: all 0.2s ease;
}

.genre-tag:hover {
    background-color: var(--secondary-color, #ff6b8b);
    color: white;
}
/* ================= 新增结束 ================= */
</style>