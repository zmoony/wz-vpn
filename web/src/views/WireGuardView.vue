<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "../lib/element-plus";

import { createPeer, deletePeer, fetchPeers, togglePeer, type WireGuardPeer } from "../api/wg";
import PageHeader from "../components/PageHeader.vue";

const PeerFormDialog = defineAsyncComponent(() => import("../components/PeerFormDialog.vue"));
const PeerQrDialog = defineAsyncComponent(() => import("../components/PeerQrDialog.vue"));

const loading = ref(false);
const createLoading = ref(false);
const dialogOpen = ref(false);
const qrDialogOpen = ref(false);
const peers = ref<WireGuardPeer[]>([]);
const generatedConf = ref("");
const generatedQr = ref("");

const onlineCount = computed(() => peers.value.filter((item) => item.online).length);
const enabledCount = computed(() => peers.value.filter((item) => item.enabled).length);
const summaryCards = computed(() => [
  { label: "Peer 总数", value: String(peers.value.length), helper: "当前已登记的客户端数量" },
  { label: "在线 Peer", value: String(onlineCount.value), helper: "正在产生握手或流量的客户端" },
  { label: "启用中", value: String(enabledCount.value), helper: "仍然允许接入接口的 Peer" },
]);

async function load() {
  loading.value = true;
  try {
    const response = await fetchPeers();
    peers.value = response.items;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载 Peer 列表失败");
  } finally {
    loading.value = false;
  }
}

async function handleCreate(payload: { name: string; description: string }) {
  createLoading.value = true;
  try {
    const result = await createPeer(payload);
    generatedConf.value = result.clientConf;
    generatedQr.value = result.qrCode;
    qrDialogOpen.value = true;
    dialogOpen.value = false;
    ElMessage.success(`已创建 Peer：${result.peer.name}`);
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "创建 Peer 失败");
  } finally {
    createLoading.value = false;
  }
}

async function handleToggle(row: WireGuardPeer) {
  try {
    await togglePeer(row.id);
    ElMessage.success(`已${row.enabled ? "禁用" : "启用"} ${row.name}`);
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "切换状态失败");
  }
}

async function handleDelete(row: WireGuardPeer) {
  await ElMessageBox.confirm(`确定删除 Peer "${row.name}" 吗？`, "删除确认", { type: "warning" });
  try {
    await deletePeer(row.id);
    ElMessage.success(`已删除 ${row.name}`);
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "删除 Peer 失败");
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell wireguard-page">
    <section class="page-card wireguard-hero">
      <div class="page-card__body wireguard-hero__body">
        <PageHeader
          title="WireGuard Peer 管理"
          description="这里已经接入真实新增、启停、删除链路；运行态握手与流量会按宿主环境能力回填。更适合先用它建立客户端，再继续联动防火墙和反向代理能力。"
        >
          <div class="wireguard-header-actions">
            <el-button :loading="loading" @click="load">刷新</el-button>
            <el-button type="primary" @click="dialogOpen = true">新增 Peer</el-button>
          </div>
        </PageHeader>

        <div class="wireguard-summary">
          <div v-for="item in summaryCards" :key="item.label" class="surface-muted wireguard-summary__item">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
            <p>{{ item.helper }}</p>
          </div>
        </div>
      </div>
    </section>

    <section v-if="peers.length === 0" class="page-card">
      <div class="page-card__body">
        <div class="empty-panel">
          <strong>当前还没有 WireGuard Peer</strong>
          <p>你可以先创建第一个客户端，系统会返回客户端配置和二维码，便于手机或笔记本直接导入。</p>
          <el-button type="primary" @click="dialogOpen = true">立即新增 Peer</el-button>
        </div>
      </div>
    </section>

    <section v-else class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="Peer 列表"
          description="启用状态、最近握手和收发流量会帮助你判断当前连接是否活跃；删除前建议先禁用观察。"
        />
        <el-table :data="peers">
          <el-table-column prop="name" label="名称" min-width="160" />
          <el-table-column prop="clientIpv4" label="客户端 IPv4" width="150" />
          <el-table-column prop="endpoint" label="Endpoint" min-width="200" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? "启用" : "禁用" }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="握手" min-width="180">
            <template #default="{ row }">
              {{ row.lastHandshakeAt || "暂无" }}
            </template>
          </el-table-column>
          <el-table-column label="流量" min-width="200">
            <template #default="{ row }">
              Rx {{ row.rxBytes }} / Tx {{ row.txBytes }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <el-space>
                <el-button link type="primary" @click="handleToggle(row)">
                  {{ row.enabled ? "禁用" : "启用" }}
                </el-button>
                <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <PeerFormDialog v-model="dialogOpen" :loading="createLoading" @submit="handleCreate" />
    <PeerQrDialog v-model="qrDialogOpen" :conf="generatedConf" :qr-code="generatedQr" />
  </div>
</template>

<style scoped>
.wireguard-hero {
  background:
    radial-gradient(circle at 88% 10%, rgba(245, 184, 65, 0.18), transparent 18%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(246, 250, 249, 0.97));
}

.wireguard-hero__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.wireguard-header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.wireguard-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.wireguard-summary__item {
  padding: 16px 18px;
}

.wireguard-summary__item span {
  display: block;
  color: var(--text-muted);
  font-size: 13px;
}

.wireguard-summary__item strong {
  display: block;
  margin-top: 8px;
  font-size: 26px;
}

.wireguard-summary__item p {
  margin: 10px 0 0;
  color: var(--text-secondary);
  line-height: 1.65;
}

@media (max-width: 1024px) {
  .wireguard-summary {
    grid-template-columns: 1fr;
  }
}
</style>
