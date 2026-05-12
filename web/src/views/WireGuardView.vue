<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

import { createPeer, deletePeer, fetchPeers, togglePeer, type WireGuardPeer } from "../api/wg";
import PageHeader from "../components/PageHeader.vue";
import PeerFormDialog from "../components/PeerFormDialog.vue";
import PeerQrDialog from "../components/PeerQrDialog.vue";

const loading = ref(false);
const createLoading = ref(false);
const dialogOpen = ref(false);
const qrDialogOpen = ref(false);
const peers = ref<WireGuardPeer[]>([]);
const generatedConf = ref("");
const generatedQr = ref("");

const onlineCount = computed(() => peers.value.filter((item) => item.online).length);

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
  <div class="page-shell">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="WireGuard Peer 管理"
          description="这轮已经接入真实新增、启停、删除链路；运行态握手与流量会按宿主环境能力回填。"
        >
          <div class="header-actions">
            <el-tag type="success">在线 {{ onlineCount }}</el-tag>
            <el-button :loading="loading" @click="load">刷新</el-button>
            <el-button type="primary" @click="dialogOpen = true">新增 Peer</el-button>
          </div>
        </PageHeader>
      </div>
    </section>

    <section class="page-card">
      <div class="page-card__body">
        <el-table :data="peers">
          <el-table-column prop="name" label="名称" min-width="160" />
          <el-table-column prop="clientIpv4" label="客户端 IPv4" width="140" />
          <el-table-column prop="endpoint" label="Endpoint" min-width="180" />
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
          <el-table-column label="流量" min-width="180">
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
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
