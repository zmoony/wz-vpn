<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

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
    <section class="login__panel page-card">
      <div class="page-card__body">
        <small class="eyebrow">Pi-Gateway</small>
        <h1>家庭网关统一控制台</h1>
        <p>先把登录和 WireGuard 管理链路打通，后续再扩展 DDNS、证书、反代和防火墙。</p>

        <el-form label-position="top" @submit.prevent="submit">
          <el-form-item label="用户名">
            <el-input v-model="form.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.password" type="password" show-password />
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
  place-items: center;
  padding: 24px;
}

.login__panel {
  width: min(100%, 460px);
}

.eyebrow {
  color: var(--brand);
  font-weight: 600;
}

h1 {
  margin: 12px 0 8px;
  font-size: 32px;
}

p {
  margin: 0 0 24px;
  color: var(--muted);
  line-height: 1.6;
}

.login__submit {
  width: 100%;
  margin-top: 6px;
}
</style>
