import { createRouter, createWebHistory } from "vue-router";

import AppLayout from "../layouts/AppLayout.vue";
import { useAuthStore } from "../stores/auth";

const LoginView = () => import("../views/LoginView.vue");
const DashboardView = () => import("../views/DashboardView.vue");
const WireGuardView = () => import("../views/WireGuardView.vue");
const ProxyView = () => import("../views/ProxyView.vue");
const DDNSView = () => import("../views/DDNSView.vue");
const CertsView = () => import("../views/CertsView.vue");
const FirewallView = () => import("../views/FirewallView.vue");
const SettingsView = () => import("../views/SettingsView.vue");

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/login", name: "login", component: LoginView },
    {
      path: "/",
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        { path: "", name: "dashboard", component: DashboardView },
        { path: "wireguard", name: "wireguard", component: WireGuardView },
        { path: "proxy", name: "proxy", component: ProxyView },
        { path: "ddns", name: "ddns", component: DDNSView },
        { path: "certs", name: "certs", component: CertsView },
        { path: "firewall", name: "firewall", component: FirewallView },
        { path: "settings", name: "settings", component: SettingsView },
      ],
    },
  ],
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth) {
    if (!auth.initialized) {
      await auth.restore();
    }

    if (!auth.isAuthenticated) {
      return { name: "login" };
    }
  }

  if (to.name === "login" && auth.isAuthenticated) {
    return { name: "dashboard" };
  }

  return true;
});

export default router;
