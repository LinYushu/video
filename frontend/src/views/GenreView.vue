<template>
    <div class="container">
        <h1>标签: {{ genreName }}</h1>
        <div class="video-grid">
            <VideoCard v-for="video in videos" :key="video.id" :video="video" @click="navigateToDetail(video)" />
        </div>
    </div>
</template>

<script>
import VideoCard from '../components/VideoCard.vue'
import videosApi from '../api/videos'

export default {
    name: 'GenreView',
    components: { VideoCard },
    props: ['genreName'],
    data() {
        return {
            videos: []
        }
    },
    async created() {
        this.videos = await videosApi.getVideoListByGenre(this.genreName)
    },
    methods: {
        navigateToDetail(video) {
            // 依赖我们在后端新加的 video.actress 字段
            this.$router.push({ 
                name: 'detail', 
                params: { actressName: video.actress, id: video.id } 
            })
        }
    }
}
</script>

<style scoped>
/* 样式与 ActressView.vue 保持一致 */
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