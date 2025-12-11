import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import ResumeBuilderView from '@/views/ResumeBuilderView.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'resume-builder',
    component: ResumeBuilderView
  }
];

const router = createRouter({
  history: createWebHistory('/'),
  routes
});

export default router;
