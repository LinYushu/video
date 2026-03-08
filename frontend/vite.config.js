import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // 监听所有地址，暴露到局域网
    port: 5173       // 保持默认端口不变
  }
})
