// main.js
import { createApp } from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.css';
import { createRouter, createWebHistory } from 'vue-router';

// ルートの定義
const routes = [
    { path: '/', component: App } // 他のルート定義もここに追加
];

// ルーターのインスタンスを作成
const router = createRouter({
    history: createWebHistory(),
    routes,
});

// アプリケーションの作成とルーターのマウント
createApp(App).use(router).mount('#app');