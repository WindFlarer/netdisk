<template>
    <div>
        <form id="myform" class="login" @submit.prevent="login">
            <p>login</p>
            <input type="text" placeholder="用户名" name="username" v-model="userName">
            <input type="password" placeholder="密码" name="password" v-model="password">
            <input type="submit" class="btn" value="登录">
            <input type="button" class="btn" value="注册" v-on:click="goToRegister">
        </form>

    </div>
</template>

<!-- -------------------------------------------------------------------------------------------- -->
<script  setup lang="ts">

import { ref, onBeforeUnmount, onBeforeMount } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

// 定义响应式数据
const userName = ref('');
const password = ref('');


// 获取路由实例
const router = useRouter();

// 定义方法
async function login() {
    try {
        const response = await axios.post('http://localhost:8080/user/login', {
            userName: userName.value,
            password: password.value,
        })


        if (response.data.code === 0) {
            // 保存JWT令牌到本地存储
            localStorage.setItem('token', response.data.token)
            window.alert('登陆成功');

            // 跳转到主页面
            router.push('/');
        } else {
            window.alert('用户名或密码错误');
        }
    } catch (error) {
        console.error(error);
        window.alert('网络错误，请稍后重试');
    }
}

// 跳转注册页面
function goToRegister() {
    router.push('/register');
}

// 避免影响其他组件
onBeforeMount(() => {
    document.querySelector("body")?.classList.add("loginBgc")
});

onBeforeUnmount(() => {
    document.querySelector("body")?.classList.remove("loginBgc")
});

</script>

<!-- -------------------------------------------------------------------------------------------- -->
<style>
/*无法选中文字*/
* {
    user-select: none;
}

/*设置背景*/
.loginBgc {
    background-image: url("../assets/wallpaper.jpg");
    background-attachment: fixed;
    background-size: cover;
}

/*设置Login居中, 背景为whitesmoke, 大小为400px正方形, 圆角25px, 文字居中, 内边距上下5,左右40,内边距不算整体长度*/
.login {
    /*设置图片居中*/
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    /*背景颜色*/
    background-color: whitesmoke;
    /*设置长宽*/
    width: 400px;
    height: 400px;
    /*设置圆角*/
    border-radius: 25px;
    /*文字居中*/
    text-align: center;
    /*设置内边距*/
    padding: 5px 40px;
    /*内边距不算在整体距离中*/
    box-sizing: border-box;
}

/*设置login文字*/
p {
    font-size: 42px;
    font-weight: 600;
    color: black;
}

/*宽度与父级相同, 高度为48px 颜色与父级相同, 下面与下一个元素间隔10px, 有一个下划线*/
input {
    /*空间 大小相关*/
    width: 100%;
    height: 48px;
    font-size: 22px;
    margin-bottom: 10px;
    /*颜色相关*/
    outline: none;
    background-color: whitesmoke;
    border: none;
    border-bottom: 2px solid silver;
    color: black;
}

.btn {

    width: 38%;
    height: 48px;
    margin-top: 40px;
    margin-left: 15px;
    margin-right: 15px;
    font-size: 28px;
    font-weight: 600;


    background-color: #59c2c5;
    border-radius: 8px;
    color: white;
}

.btn:hover {
    background-color: #59c2a0;
    cursor: pointer;
}
</style>