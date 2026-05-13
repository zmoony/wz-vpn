import { createRouter, createWebHistory } from "vue-router";

import AppLayout from "../layouts/AppLayout.vue";
import { useAuthStore } from "../stores/auth";
import CertsView from "../views/CertsView.vue";
import DashboardView from "../views/DashboardView.vue";
import DDNSView from "../views/DDNSView.vue";
import FirewallView from "../views/FirewallView.vue";
import LoginView from "../views/LoginView.vue";
import ProxyView from "../views/ProxyView.vue";
import SettingsView from "../views/SettingsView.vue";
import WireGuardView from "../views/WireGuardView.vue";

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
