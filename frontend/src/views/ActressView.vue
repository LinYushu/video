<template>
    <div class="container">
        <h1>{{ actressName }} 的作品</h1>
        <div class="video-grid">
            <VideoCard v-for="video in videos" :key="video.id" :video="video" @click="navigateToDetail(video.id)" />
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
            videos: []
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

/* 手机端优先：一行显示一个视频 */
.video-grid {
    display: grid;
    grid-template-columns: 1fr; /* 强制一行一个 */
    gap: 20px; 
    padding: 10px 0;
}

/* 当屏幕宽度大于 768px（平板和电脑）时覆盖上方样式 */
@media (min-width: 768px) {
    .container { 
        padding: 2rem; 
    }
    .video-grid {
        /* PC端：强制一行显示 4 个 */
        grid-template-columns: repeat(4, 1fr);
        gap: 25px; 
    }
}
</style>