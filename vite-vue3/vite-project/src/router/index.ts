import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'


const routes: Array<RouteRecordRaw> = [
    {
        //路由初始指向
        path: '/',
        name: 'HelloWorld',
        component: () => import('../components/HelloWorld.vue'),
    },
    {
        // 登录
        path: '/login',
        name: 'Login',
        component: () => import('../components/Login.vue'),
    },
    {
        // 注册
        path: '/register',
        name: 'Register',
        component: () => import('../components/Register.vue'),
    },
    {
        // 主页
        path: '/index',
        name: 'Index',
        component: () => import('../components/Index.vue'),
    },
    {
        // 测试
        path: '/test',
        name: 'Test',
        component: () => import('../components/TT.vue'),
    }
    ,
    {
        // 测试
        path: '/test2',
        name: 'Test2',
        component: () => import('../components/TT2.vue'),
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
