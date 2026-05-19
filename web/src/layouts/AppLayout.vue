<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useAuthStore } from "../stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

const navigation = [
  { label: "总览", to: "/" },
  { label: "WireGuard", to: "/wireguard" },
  { label: "反向代理", to: "/proxy" },
  { label: "DDNS", to: "/ddns" },
  { label: "证书", to: "/certs" },
  { label: "防火墙", to: "/firewall" },
  { label: "设置", to: "/settings" },
];

const title = computed(() => {
  const current = navigation.find((item) => item.to === route.path);
  return current?.label ?? "Pi-Gateway";
});

async function handleLogout() {
  await auth.signOut();
  await router.push({ name: "login" });
}
</script>

<template>
  <div class="layout">
    <aside class="layout__sidebar">
      <div class="layout__brand surface-muted">
        <small>Raspberry Pi Gateway</small>
        <strong>Pi-Gateway</strong>
        <p>面向家庭网络的统一控制台，把 VPN、证书、反代和防火墙放进同一片山水界面里。</p>
      </div>

      <nav class="nav">
        <RouterLink
          v-for="item in navigation"
          :key="item.to"
          :to="item.to"
          class="nav__item"
          :class="{ 'nav__item--active': route.path === item.to }"
        >
          <span class="nav__item-label">{{ item.label }}</span>
          <small>{{ route.path === item.to ? "当前页面" : "进入" }}</small>
        </RouterLink>
      </nav>

      <div class="layout__sidebar-note surface-muted">
        <strong>运维提示</strong>
        <p>推荐先在“设置”页补齐 WireGuard 参数，再逐步启用防火墙、证书和反向代理能力。</p>
      </div>
    </aside>

    <div class="layout__content">
      <header class="layout__header">
        <div class="layout__header-main">
          <span class="layout__crumb">当前模块</span>
          <div class="layout__headline">
            <h1>{{ title }}</h1>
            <p>延续自然配色和轻量信息层级，让高风险配置也能更清楚地完成。</p>
          </div>
        </div>

        <div class="layout__header-actions">
          <div class="layout__user-chip surface-muted">
            <span class="layout__user-dot" />
            <div>
              <strong>{{ auth.user?.username }}</strong>
              <small>已登录</small>
            </div>
          </div>
          <button class="ghost-button" type="button" @click="handleLogout">退出登录</button>
        </div>
      </header>

      <main class="layout__main">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.layout {
  height: 100vh;
  display: grid;
  grid-template-columns: 296px 1fr;
  overflow: hidden;
}

.layout__sidebar {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 24px;
  height: 100vh;
  overflow-y: auto;
  background:
    linear-gradient(180deg, rgba(18, 66, 63, 0.96), rgba(13, 45, 43, 0.98)),
    linear-gradient(180deg, #0d2d2b 0%, #173634 100%);
  color: #f7f7f2;
}

.layout__brand,
.layout__sidebar-note {
  padding: 18px;
  color: rgba(255, 255, 255, 0.92);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.12);
}

.layout__brand small,
.layout__sidebar-note p {
  color: rgba(247, 247, 242, 0.72);
}

.layout__brand strong {
  display: block;
  margin-top: 6px;
  font-size: 30px;
  line-height: 1.15;
}

.layout__brand p,
.layout__sidebar-note p {
  margin: 10px 0 0;
  line-height: 1.7;
}

.layout__sidebar-note strong {
  color: #f5b841;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.nav__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1px solid transparent;
  border-radius: 16px;
  color: rgba(247, 247, 242, 0.78);
  background: rgba(255, 255, 255, 0.03);
}

.nav__item small {
  color: rgba(247, 247, 242, 0.54);
}

.nav__item--active,
.nav__item:hover {
  color: #ffffff;
  background: linear-gradient(135deg, rgba(22, 152, 142, 0.28), rgba(245, 184, 65, 0.14));
  border-color: rgba(245, 184, 65, 0.18);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.06);
}

.nav__item-label {
  font-weight: 600;
}

.layout__content {
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: 100vh;
  min-height: 0;
  overflow: hidden;
}

.layout__header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 18px;
  padding: 26px var(--space-page) 0;
  background:
    linear-gradient(180deg, rgba(251, 250, 247, 0.95) 0%, rgba(245, 245, 245, 0.88) 72%, rgba(245, 245, 245, 0) 100%);
  backdrop-filter: blur(12px);
}

.layout__header-main {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.layout__crumb {
  color: var(--earth);
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.layout__headline h1 {
  margin: 0;
  font-size: 34px;
  line-height: 1.15;
}

.layout__headline p {
  margin: 10px 0 0;
  color: var(--text-secondary);
  line-height: 1.7;
}

.layout__header-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.layout__user-chip {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
}

.layout__user-chip strong {
  display: block;
}

.layout__user-chip small {
  color: var(--text-muted);
}

.layout__user-dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: linear-gradient(180deg, var(--success), var(--brand));
  box-shadow: 0 0 0 6px rgba(34, 197, 94, 0.08);
}

.ghost-button {
  min-height: 44px;
  padding: 10px 16px;
  border: 1px solid rgba(22, 152, 142, 0.18);
  border-radius: 12px;
  color: var(--brand-deep);
  background: rgba(255, 255, 255, 0.88);
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.layout__main {
  flex: 1;
  min-height: 0;
  padding: 24px var(--space-page) 30px;
  overflow-y: auto;
}

@media (max-width: 1120px) {
  .layout {
    height: auto;
    grid-template-columns: 1fr;
    overflow: visible;
  }

  .layout__sidebar {
    gap: 12px;
    padding: 18px var(--space-page);
    height: auto;
    overflow: visible;
  }

  .nav {
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 4px;
  }

  .nav__item {
    min-width: 150px;
  }

  .layout__content {
    height: auto;
    overflow: visible;
  }

  .layout__header {
    position: static;
    background: transparent;
    backdrop-filter: none;
  }

  .layout__main {
    overflow: visible;
  }
}

@media (max-width: 720px) {
  .layout__header {
    flex-direction: column;
  }

  .layout__header-actions {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
