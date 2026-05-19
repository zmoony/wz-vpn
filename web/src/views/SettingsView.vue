<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "../lib/element-plus";

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
    <section class="page-card settings-hero">
      <div class="page-card__body settings-hero__body">
        <PageHeader
          title="系统设置"
          description="这里集中维护当前会影响 WireGuard 主链路的关键参数，同时把运行时环境路径收成只读摘要，避免你在修改配置时来回切换上下文。"
        >
          <el-space>
            <el-button :loading="loading" @click="load">刷新设置</el-button>
            <el-button type="primary" :loading="saving" @click="save">保存设置</el-button>
          </el-space>
        </PageHeader>

        <div class="settings-hero__tips">
          <div class="surface-muted settings-hero__tip">
            <strong>会实时影响</strong>
            <p>接口名、VPN 子网、Endpoint 和公钥会影响后续 Peer 创建与部分运行时操作。</p>
          </div>
          <div class="surface-muted settings-hero__tip">
            <strong>仍由环境变量控制</strong>
            <p>证书目录、Nginx 路径、DDNS 配置路径等进程级配置，当前仍以启动时 env 为准。</p>
          </div>
        </div>
      </div>
    </section>

    <section class="settings-grid">
      <article class="page-card">
        <div class="page-card__body">
          <h3 class="section-title">WireGuard 运行参数</h3>
          <p class="section-subtitle">建议先填好接口、地址和公钥，再继续配置 Peer、反代或防火墙规则。</p>

          <el-form label-position="top" class="settings-form">
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
              <el-input v-model="form.publicKey" type="textarea" :rows="5" placeholder="WireGuard server public key" />
            </el-form-item>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <h3 class="section-title">运行环境摘要</h3>
          <p class="section-subtitle">这部分是当前进程的只读环境快照，适合排查路径和服务启动参数是否正确。</p>

          <div class="runtime-grid">
            <div class="surface-muted runtime-item">
              <span>运行环境</span>
              <strong>{{ runtime?.appEnv || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>监听地址</span>
              <strong>{{ runtime?.listenAddr || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>数据目录</span>
              <strong>{{ runtime?.dataDir || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>数据库</span>
              <strong>{{ runtime?.databasePath || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>前端产物</span>
              <strong>{{ runtime?.webDistDir || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>Nginx 站点目录</span>
              <strong>{{ runtime?.nginxSitesDir || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>DDNS 配置文件</span>
              <strong>{{ runtime?.ddnsGoConfigPath || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>证书目录</span>
              <strong>{{ runtime?.certsInstallRoot || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>防火墙规则文件</span>
              <strong>{{ runtime?.nftablesRulesPath || "-" }}</strong>
            </div>
            <div class="surface-muted runtime-item">
              <span>待确认秒数</span>
              <strong>{{ runtime?.firewallPendingSeconds ?? "-" }}</strong>
            </div>
          </div>

          <el-alert
            class="settings-alert"
            title="当前设置边界"
            type="info"
            :closable="false"
            description="这里先保存会影响 WireGuard 和部分防火墙渲染的关键配置。证书、Nginx、DDNS 等路径型配置仍以启动进程时的环境变量为准。"
          />
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.settings-hero {
  background:
    radial-gradient(circle at 88% 12%, rgba(245, 184, 65, 0.18), transparent 20%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(247, 250, 249, 0.96));
}

.settings-hero__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-hero__tips {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.settings-hero__tip {
  padding: 16px 18px;
}

.settings-hero__tip strong {
  display: block;
  color: var(--text);
}

.settings-hero__tip p {
  margin: 10px 0 0;
  color: var(--text-secondary);
  line-height: 1.7;
}

.settings-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(340px, 0.92fr);
  gap: 20px;
}

.settings-form {
  margin-top: 20px;
}

.runtime-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 20px;
}

.runtime-item {
  min-height: 92px;
  padding: 16px;
}

.runtime-item span {
  display: block;
  color: var(--text-muted);
  font-size: 13px;
}

.runtime-item strong {
  display: block;
  margin-top: 10px;
  color: var(--text);
  line-height: 1.55;
  word-break: break-all;
}

.settings-alert {
  margin-top: 18px;
}

@media (max-width: 1180px) {
  .settings-grid,
  .settings-hero__tips,
  .runtime-grid {
    grid-template-columns: 1fr;
  }
}
</style>
