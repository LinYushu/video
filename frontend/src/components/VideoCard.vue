<template>
    <div class="video-card" @click="$emit('click')">
        <div class="poster-container">
            <img class="poster pc-poster" :src="video.poster" :alt="video.title">
            
            <img class="poster mobile-poster" :src="video.fanart" :alt="video.title">
        </div>
        <div class="info">
            <h3>{{ video.title }}</h3>
        </div>
    </div>
</template>

<script>
export default {
    name: 'VideoCard',
    props: {
        video: {
            type: Object,
            required: true
        }
    }
}
</script>

<style scoped>
.video-card {
    cursor: pointer;
    transition: all 0.3s ease;
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    width: 100%;
}

.video-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

/* --- 手机端默认样式 --- */
.poster-container {
    position: relative;
    width: 100%;
    /* 手机端使用横版比例，适配 fanart 图片 */
    aspect-ratio: 16 / 9; 
    overflow: hidden;
}

/* 关键：手机端隐藏 PC 版图片，显示手机版图片 */
.pc-poster {
    display: none;
}
.mobile-poster {
    display: block;
}

.poster {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    /* 保证图片铺满容器并保持比例，多余部分会被裁剪 */
    object-fit: cover; 
    /* 可选：居中显示图片，防止重点内容被裁掉 */
    object-position: center; 
}

.info {
    padding: 12px;
    background: white;
}

h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 500;
    color: #333;
    
    /* 手机端 2 行截断 */
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2; 
    line-clamp: 2;
    overflow: hidden;
    text-overflow: ellipsis;
    
    line-height: 1.5; 
    max-height: 3em; 
    word-break: break-all;
}

/* --- 大于 768px（PC端）时的样式覆盖 --- */
@media (min-width: 768px) {
    .poster-container {
        aspect-ratio: auto; 
        /* PC 端恢复 9:16 的竖向海报比例 */
        padding-top: 137.78%; 
    }

    /* 关键：PC端隐藏手机版图片，显示PC版图片 */
    .pc-poster {
        display: block;
    }
    .mobile-poster {
        display: none;
    }

    h3 {
        /* PC端 4 行截断 */
        -webkit-line-clamp: 4; 
        line-clamp: 4;
        max-height: 6em; 
    }
}
</style>