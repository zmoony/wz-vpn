<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";

import PageHeader from "../components/PageHeader.vue";
import { fetchSettings, updateSettings, type RuntimeSettings, type WireGuardSettings } from "../api/settings";

const loading = ref(false);
const saving = ref(false);
const runtime = ref<RuntimeSettings | null>(null);

const form = reactive<WireGuardSettings>({
  interface: "wg0",
  subnetV4: "10.66.66.0/24",
  serverV4: "10.66.66.1/24",
  dns: "192.168.1.1",
  endpoint: "vpn.example.com:51820",
  publicKey: "",
});

async function load() {
  loading.value = true;
  try {
    const response = await fetchSettings();
    runtime.value = response.runtime;
    form.interface = response.wireGuard.interface;
    form.subnetV4 = response.wireGuard.subnetV4;
    form.serverV4 = response.wireGuard.serverV4;
    form.dns = response.wireGuard.dns;
    form.endpoint = response.wireGuard.endpoint;
    form.publicKey = response.wireGuard.publicKey;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载系统设置失败");
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    const response = await updateSettings({ wireGuard: { ...form } });
    runtime.value = response.runtime;
    ElMessage.success("系统设置已保存");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存系统设置失败");
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell settings-page">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="系统设置"
          description="这轮先把 WireGuard 关键参数接成真实可保存配置；修改后会影响后续 Peer 创建、运行时接口和部分转发渲染。"
        >
          <el-space>
            <el-button :loading="loading" @click="load">刷新设置</el-button>
            <el-button type="primary" :loading="saving" @click="save">保存设置</el-button>
          </el-space>
        </PageHeader>
      </div>
    </section>

    <section class="settings-grid">
      <article class="page-card">
        <div class="page-card__body">
          <h3>WireGuard 运行参数</h3>
          <el-form label-position="top">
            <el-form-item label="接口名">
              <el-input v-model="form.interface" placeholder="wg0" />
            </el-form-item>
            <el-form-item label="VPN 子网">
              <el-input v-model="form.subnetV4" placeholder="10.66.66.0/24" />
            </el-form-item>
            <el-form-item label="服务端地址">
              <el-input v-model="form.serverV4" placeholder="10.66.66.1/24" />
            </el-form-item>
            <el-form-item label="客户端 DNS">
              <el-input v-model="form.dns" placeholder="192.168.1.1" />
            </el-form-item>
            <el-form-item label="公网 Endpoint">
              <el-input v-model="form.endpoint" placeholder="vpn.example.com:51820" />
            </el-form-item>
            <el-form-item label="服务端公钥">
              <el-input v-model="form.publicKey" type="textarea" :rows="4" placeholder="WireGuard server public key" />
            </el-form-item>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <h3>运行环境摘要</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="运行环境">{{ runtime?.appEnv || "-" }}</el-descriptions-item>
            <el-descriptions-item label="监听地址">{{ runtime?.listenAddr || "-" }}</el-descriptions-item>
            <el-descriptions-item label="数据目录">{{ runtime?.dataDir || "-" }}</el-descriptions-item>
            <el-descriptions-item label="数据库">{{ runtime?.databasePath || "-" }}</el-descriptions-item>
            <el-descriptions-item label="前端产物">{{ runtime?.webDistDir || "-" }}</el-descriptions-item>
            <el-descriptions-item label="Nginx 站点目录">{{ runtime?.nginxSitesDir || "-" }}</el-descriptions-item>
            <el-descriptions-item label="DDNS 配置文件">{{ runtime?.ddnsGoConfigPath || "-" }}</el-descriptions-item>
            <el-descriptions-item label="证书目录">{{ runtime?.certsInstallRoot || "-" }}</el-descriptions-item>
            <el-descriptions-item label="防火墙规则文件">{{ runtime?.nftablesRulesPath || "-" }}</el-descriptions-item>
            <el-descriptions-item label="防火墙待确认秒数">{{ runtime?.firewallPendingSeconds ?? "-" }}</el-descriptions-item>
          </el-descriptions>

          <el-alert
            class="settings-alert"
            title="说明"
            type="info"
            :closable="false"
            description="这里先保存会影响 WireGuard 和部分防火墙渲染的关键配置。环境路径类配置仍以当前进程启动时的 env 为准。"
          />
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.settings-grid {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 20px;
}

.settings-alert {
  margin-top: 18px;
}

@media (max-width: 1200px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>
