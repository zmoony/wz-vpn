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
      <div class="brand">
        <small>Raspberry Pi Gateway</small>
        <strong>Pi-Gateway</strong>
      </div>
      <nav class="nav">
        <RouterLink
          v-for="item in navigation"
          :key="item.to"
          :to="item.to"
          class="nav__item"
          :class="{ 'nav__item--active': route.path === item.to }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </aside>

    <div class="layout__content">
      <header class="layout__header">
        <div>
          <small>当前模块</small>
          <h1>{{ title }}</h1>
        </div>
        <div class="layout__header-actions">
          <span>{{ auth.user?.username }}</span>
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
  min-height: 100vh;
  display: grid;
  grid-template-columns: 280px 1fr;
}

.layout__sidebar {
  padding: 28px 22px;
  background: linear-gradient(180deg, #081120 0%, #0e1c31 100%);
  color: #f8fbff;
}

.brand {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 28px;
}

.brand small {
  color: rgba(248, 251, 255, 0.7);
}

.brand strong {
  font-size: 28px;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav__item {
  padding: 12px 14px;
  border-radius: 14px;
  color: rgba(248, 251, 255, 0.76);
}

.nav__item--active,
.nav__item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
}

.layout__content {
  display: flex;
  flex-direction: column;
}

.layout__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 24px 28px 0;
}

.layout__header small {
  color: var(--muted);
}

.layout__header h1 {
  margin: 6px 0 0;
  font-size: 28px;
}

.layout__header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.ghost-button {
  padding: 10px 14px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
}

.layout__main {
  padding: 28px;
}

@media (max-width: 1080px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .layout__sidebar {
    display: none;
  }
}
</style>
