import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import DetailView from '../views/DetailView.vue'
import ActressView from '../views/ActressView.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: HomeView
    },
    {
        path: '/actress/:actressName',
        name: 'actress',
        component: ActressView,
        props: true
    },
    {
        path: '/video/:actressName/:id',
        name: 'detail',
        component: DetailView,
        props: true
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router