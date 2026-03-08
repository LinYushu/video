<template>
    <div class="container">
        <h1>我的小姐姐</h1>
        <div class="actress-grid">
            <ActressCard 
                v-for="actress in actresses" 
                :key="actress.id" 
                :actress="actress"
                @click="navigateToActress(actress.id)"
            />
        </div>
    </div>
</template>

<script>
import ActressCard from '../components/ActressCard.vue'
import videosApi from '../api/videos'

export default {
    name: 'HomeView',
    components: { ActressCard },
    data() {
        return {
            actresses: [],
            scrollPosition: 0
        }
    },
    async created() {
        if (!this.actresses.length) {
            this.actresses = await videosApi.getActressList()
        }
    },
    activated() {
        window.scrollTo(0, this.scrollPosition)
    },
    beforeRouteLeave(to, from, next) {
        this.scrollPosition = window.scrollY
        next()
    },
    methods: {
        navigateToActress(actressName) {
            this.$router.push({ name: 'actress', params: { actressName } })
        }
    }
}
</script>

<style scoped>
/* 手机端默认减小内边距，解决左侧大量留白问题 */
.container {
  padding: 10px 15px; 
}

h1 {
  color: var(--text-color, #333);
  margin-bottom: 1.5rem;
  font-weight: 600;
  position: relative;
  display: inline-block;
  margin-top: 10px;
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

.actress-grid {
  display: grid;
  /* 强制一行显示 4 个，平分宽度 */
  grid-template-columns: repeat(4, 1fr);
  gap: 10px; /* 适当缩小手机端的间距 */
  padding: 10px 0;
  justify-items: center; /* 保证卡片在各自的网格中居中 */
}

/* 当屏幕宽度大于 768px（平板和电脑）时覆盖上方样式 */
@media (min-width: 768px) {
  .container {
    padding: 2rem; /* PC端恢复原本的内边距 */
  }
  .actress-grid {
    grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
    gap: 35px;
  }
}
</style>