<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "../lib/element-plus";

import { useAuthStore } from "../stores/auth";

const router = useRouter();
const auth = useAuthStore();
const loading = ref(false);
const form = reactive({
  username: "admin",
  password: "admin123456",
});

async function submit() {
  loading.value = true;
  try {
    await auth.signIn(form.username, form.password);
    await router.push({ name: "dashboard" });
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "登录失败");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login">
    <section class="login__hero">
      <div class="login__copy">
        <span class="login__eyebrow">漓江青 · 网关控制台</span>
        <h1>把家庭网络运维，收进一块更清晰的山水面板里。</h1>
        <p>
          你可以在同一套界面里管理 WireGuard、反向代理、证书、DDNS 和防火墙，
          既保留系统编排的控制力，也降低首次配置时的认知压力。
        </p>
      </div>

      <div class="login__highlights">
        <article class="surface-muted login__highlight">
          <strong>配置集中</strong>
          <p>统一入口管理 VPN、证书和防火墙，不再分散在多个命令和配置文件之间。</p>
        </article>
        <article class="surface-muted login__highlight">
          <strong>风险可见</strong>
          <p>高风险操作有预览、确认和回滚链路，让变更更可控。</p>
        </article>
        <article class="surface-muted login__highlight">
          <strong>风格统一</strong>
          <p>以自然配色和轻量层级组织复杂信息，提升可读性而不过度装饰。</p>
        </article>
      </div>
    </section>

    <section class="login__panel page-card">
      <div class="page-card__body">
        <div class="login__panel-head">
          <small>Pi-Gateway</small>
          <h2>进入控制台</h2>
          <p>使用管理员账号登录，继续你的网关配置与运维工作。</p>
        </div>

        <el-form label-position="top" @submit.prevent="submit">
          <el-form-item label="用户名">
            <el-input v-model="form.username" placeholder="请输入管理员用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
          <el-button type="primary" size="large" class="login__submit" :loading="loading" @click="submit">
            进入控制台
          </el-button>
        </el-form>
      </div>
    </section>
  </div>
</template>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(360px, 420px);
  gap: 28px;
  align-items: stretch;
  padding: 28px;
}

.login__hero {
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 28px;
  padding: 42px;
  border-radius: 26px;
  background:
    radial-gradient(circle at 10% 0%, rgba(245, 184, 65, 0.28), transparent 30%),
    radial-gradient(circle at 85% 15%, rgba(22, 152, 142, 0.22), transparent 34%),
    linear-gradient(160deg, rgba(255, 255, 255, 0.92) 0%, rgba(244, 248, 246, 0.95) 46%, rgba(237, 243, 241, 0.98) 100%);
  border: 1px solid rgba(139, 110, 78, 0.12);
  box-shadow: var(--shadow-card);
}

.login__hero::after {
  content: "";
  position: absolute;
  right: -60px;
  bottom: -40px;
  width: 260px;
  height: 260px;
  border-radius: 44% 56% 65% 35% / 48% 42% 58% 52%;
  background: linear-gradient(180deg, rgba(22, 152, 142, 0.1), rgba(139, 110, 78, 0.1));
}

.login__copy {
  position: relative;
  z-index: 1;
  max-width: 680px;
}

.login__eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--earth);
  font-size: 13px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.login__eyebrow::before {
  content: "";
  width: 28px;
  height: 2px;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--brand), var(--accent));
}

.login__copy h1 {
  max-width: 740px;
  margin: 18px 0 16px;
  font-size: 48px;
  line-height: 1.12;
}

.login__copy p {
  max-width: 640px;
  margin: 0;
  color: var(--text-secondary);
  font-size: 18px;
  line-height: 1.85;
}

.login__highlights {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.login__highlight {
  padding: 18px;
}

.login__highlight strong {
  display: block;
  color: var(--text);
  font-size: 18px;
}

.login__highlight p {
  margin: 10px 0 0;
  color: var(--text-secondary);
  line-height: 1.75;
}

.login__panel {
  align-self: center;
  background: rgba(255, 255, 255, 0.96);
}

.login__panel-head small {
  color: var(--brand);
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.login__panel-head h2 {
  margin: 14px 0 8px;
  font-size: 32px;
  line-height: 1.2;
}

.login__panel-head p {
  margin: 0 0 24px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.login__submit {
  width: 100%;
  margin-top: 10px;
}

@media (max-width: 1180px) {
  .login {
    grid-template-columns: 1fr;
  }

  .login__highlights {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .login {
    padding: 18px;
  }

  .login__hero {
    padding: 28px 22px;
  }

  .login__copy h1 {
    font-size: 36px;
  }
}
</style>
