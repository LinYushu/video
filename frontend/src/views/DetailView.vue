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

            <Gallery :images="video.fanarts" />
        </div>
    </div>
</template>

<script>
import Gallery from '../components/Gallery.vue'
import videosApi from '../api/videos'

export default {
    components: { Gallery },
    // 接收两个参数：演员名称和视频ID
    props: ['actressName', 'id'], 
    data() {
        return {
            video: null
        }
    },
    async created() {
        // 请求详情时传入两个参数
        this.video = await videosApi.getVideoDetail(this.actressName, this.id)
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

/* 完善视频播放器的样式 */
.video-player {
  margin-top: 2rem;
  margin-bottom: 2rem;
  border-radius: 12px;
  overflow: hidden; /* 保证视频的四个角也能贴合圆角 */
  box-shadow: 0 5px 25px rgba(0, 0, 0, 0.1);
  background-color: #000; /* 视频未加载时显示黑色背景 */
  display: flex;
  justify-content: center;
}

/* 针对 video 标签本身的约束 */
.video-player video {
  width: 100%;
  max-height: 75vh; /* 限制视频最大高度，防止大屏幕上占据过多空间 */
  display: block;
  outline: none;
}
</style>