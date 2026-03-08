<template>
    <div class="actress-card" @click="$emit('click')">
        <div class="avatar-wrapper">
            <img :src="actress.poster" :alt="actress.title" class="avatar-img" />
        </div>
        <div class="actress-name">{{ actress.title }}</div>
    </div>
</template>

<script>
export default {
    name: 'ActressCard',
    props: {
        actress: {
            type: Object,
            required: true
        }
    }
}
</script>

<style scoped>
.actress-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s ease;
  width: 100%;
}

.actress-card:hover {
  transform: translateY(-5px);
}

.avatar-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  /* 强制设定为圆形，并加上 overflow: hidden 裁剪多余部分 */
  border-radius: 50% !important; 
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 8px;
  border: 2px solid #fff;
  transition: box-shadow 0.2s ease;
  /* 确保内部图片不会撑破圆角 */
  transform: translateZ(0); 
}

.actress-card:hover .avatar-wrapper {
  box-shadow: 0 8px 16px rgba(255, 107, 139, 0.25);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.actress-name {
  font-size: 0.8rem;
  color: var(--text-color, #333);
  font-weight: 500;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  padding: 0 2px;
}

@media (min-width: 768px) {
  .avatar-wrapper {
    margin-bottom: 12px;
    border: 3px solid #fff;
    /* 在 PC 端再次强制声明圆形，防止被未知样式覆盖 */
    border-radius: 50% !important; 
  }
  .actress-name {
    font-size: 0.95rem;
  }
}
</style>