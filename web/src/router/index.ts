/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import Connect from '@/views/Connect.vue'
import Dashboard from '@/views/Dashboard.vue'
import Device from '@/views/Device.vue'
import DeviceAI from '@/views/DeviceAI.vue'
import DeviceHistoric from '@/views/DeviceHistoric.vue'
import DeviceNerd from '@/views/DeviceNerd.vue'
import DeviceSettings from '@/views/DeviceSettings.vue'
import DeviceStats from '@/views/DeviceStats.vue'
import DeviceStatus from '@/views/DeviceStatus.vue'
import Login from '@/views/Login.vue'
import Settings from '@/views/Settings.vue'
import { createRouter, createWebHistory } from 'vue-router/auto'

const extraMenuItems = [
  { title: 'Status', icon: 'mdi-information-slab-circle', to: { name: 'DeviceStatus' } },
  { title: 'A.I.', icon: 'mdi-robot-happy', to: { name: 'DeviceAI' } },
  { title: 'Histórico', icon: 'mdi-history', to: { name: 'DeviceHistoric' } },
  { title: 'Estatística', icon: 'mdi-chart-box-outline', to: { name: 'DeviceNerd' } }, // Corrected name
  { title: 'Configurações', icon: 'mdi-list-status', to: { name: 'DeviceSettings' } },
];

const routes = [
  { path: "/", name: "Login", component: Login },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: Dashboard,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/connect",
    name: "Connect",
    component: Connect,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/device/:id",
    name: "Device",
    component: Device,
    icon: "mdi-cellphone",
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
  {
    path: "/device/:id/status",
    name: "DeviceStatus",
    component: DeviceStatus,
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
  {
    path: "/device/:id/artificial-inteligence",
    name: "DeviceAI",
    component: DeviceAI,
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
  {
    path: "/device/:id/historic",
    name: "DeviceHistoric",
    component: DeviceHistoric,
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
  {
    path: "/device/:id/nerd",
    name: "DeviceNerd",
    component: DeviceStats,
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
  {
    path: "/device/:id/settings",
    name: "DeviceSettings",
    component: DeviceSettings,
    meta: {
      requiresAuth: true,
      extraMenuItems: extraMenuItems,
    },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !isAuthenticated()) {
    next({ name: 'Login' });
  } else {
    next();
  }
});

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (!localStorage.getItem('vuetify:dynamic-reload')) {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    } else {
      console.error('Dynamic import error, reloading page did not fix it', err)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

function isAuthenticated() {
  return localStorage.getItem('isLogged') === 'true'
}

export default router
