import axios from 'axios'

const API_BASE = 'http://192.168.1.100:31471'

export default {
    async getActressList() {
        const response = await axios.get(`${API_BASE}/api/actresses`)
        return response.data.map(actress => ({
            ...actress,
            // 严谨起见，也可以加上判空，不过头像一般都有
            poster: actress.poster ? `${API_BASE}${actress.poster}` : ''
        }))
    },

    async getVideoListByActress(actressName) {
        const response = await axios.get(`${API_BASE}/api/actress/${encodeURIComponent(actressName)}`)
        return response.data.map(video => ({
            ...video,
            // 原来的 poster 处理
            poster: video.poster ? `${API_BASE}${video.poster}` : '',

            fanart: video.fanart ? `${API_BASE}${video.fanart}` : ''
        }))
    },

    async getVideoDetail(actressName, id) {
        const response = await axios.get(`${API_BASE}/api/videos/${encodeURIComponent(actressName)}/${id}`)
        const data = response.data

        return {
            ...data,
            poster: data.fanarts?.length ? `${API_BASE}${data.fanarts[0]}` : null,
            videoFile: data.videoFile ? `${API_BASE}${data.videoFile}` : null,
            fanarts: data.fanarts?.map(img => `${API_BASE}${img}`) || []
        }
    },

    async addVideo(id) {
        console.log(id);
        const pattern = /^([A-Z0-9]+)-\d+$/;

        if (!id.trim()) {
            alert('请输入视频内容', true);
            return;
        }

        if (!pattern.test(id)) {
            alert('格式错误：请输入(字母或数字)-数字的组合，如 ABC-123', true);
            return;
        }

        // 友情提示：这里的 this.isAdding 和 this.inputContent 是你在 Vue 组件里的写法
        // 放在这个纯 JS 文件里，this 指向的是默认导出的对象，而不是你的 Vue 实例哦。
        // 如果你要控制页面的 loading 状态，建议把这两个状态的修改移回你的 Vue 组件内部。
        this.isAdding = true;
        try {
            const response = await axios.get(`${API_BASE}/api/addvideo/${encodeURIComponent(id)}`, {
                headers: {
                    'Authorization': 'Bearer IBHUSDBWQHJEJOBDSW'
                }
            });
            console.log(response.data);

            if (response.status >= 200 && response.status < 300) {
                this.inputContent = ''; // 同上，建议移回 Vue 组件
                // 使用原生alert替代ElMessage
                alert('视频添加成功');
            } else {
                alert(response.data.message || '添加视频失败');
            }
        } catch (error) {
            console.error('添加视频出错:', error);
            alert(`添加视频失败: ${error.message}`);
        } finally {
            this.isAdding = false; // 同上，建议移回 Vue 组件
        }
    }
}