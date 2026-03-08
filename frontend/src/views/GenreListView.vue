<template>
    <div class="container">
        <h1>探索标签</h1>
        <div class="genre-grid">
            <div 
                class="genre-card" 
                v-for="genre in genres" 
                :key="genre.name" 
                @click="goToGenre(genre.name)"
            >
                <span class="genre-name">{{ genre.name }}</span>
                <span class="genre-count">({{ genre.count }}部)</span>
            </div>
        </div>
    </div>
</template>

<script>
import videosApi from '../api/videos'

export default {
    name: 'GenreListView',
    data() {
        return {
            genres: []
        }
    },
    async created() {
        this.genres = await videosApi.getGenreList()
    },
    methods: {
        goToGenre(name) {
            this.$router.push({ name: 'genre', params: { genreName: name } })
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

.genre-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr); /* 强制 4 列等宽 */
    gap: 6px; /* 缩小手机端卡片间距 */
    padding: 10px 0;
}

.genre-card {
    background: white;
    border-radius: 10px;
    padding: 10px 5px; /* 缩小左右内边距，给文字留出更多空间 */
    text-align: center;
    cursor: pointer;
    box-shadow: 0 5px 15px rgba(255, 107, 139, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 70px; /* 缩小手机端卡片高度 */
    overflow: hidden; /* 防止内容撑破卡片 */
}

.genre-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 25px rgba(255, 107, 139, 0.2);
}

.genre-name {
    font-size: 0.8rem; /* 手机端字体调小一点 */
    color: var(--text-color, #5a3a4a);
    font-weight: bold;
    margin-bottom: 5px;

    width: 100%;             /* 必须占满容器宽度 */
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2; /* 限制最多显示两行 */
    white-space: normal;
    overflow: hidden;        /* 隐藏超出部分 */
    text-overflow: ellipsis; /* 显示省略号 */
    padding: 0 2px;
    box-sizing: border-box;
    line-height: 1.2; /* 控制行高，避免两行文字太拥挤 */
}

.genre-count {
    font-size: 0.8rem;
    color: var(--secondary-color, #ff6b8b);
    font-weight: 500;
}

/* PC/平板端媒体查询：屏幕大于 768px 时恢复大卡片布局 */
@media (min-width: 768px) {
    .genre-grid {
        /* PC 端自动填充，每个最小 140px宽 */
        grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
        gap: 15px;
    }
    
    .genre-card {
        height: 100px;
        padding: 20px 10px;
    }

    .genre-name {
        font-size: 1.0rem; /* 恢复大字体 */
    }

    .genre-count {
        font-size: 0.9rem;
    }
}
</style>