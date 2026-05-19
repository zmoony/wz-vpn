<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "../lib/element-plus";

import PageHeader from "../components/PageHeader.vue";
import { fetchConntrackEntries, type ConntrackEntry } from "../api/conntrack";
import {
  applyFirewallRules,
  confirmFirewallRules,
  createFirewallForwardRule,
  createFirewallRule,
  deleteFirewallForwardRule,
  deleteFirewallRule,
  fetchFirewallForwardConfig,
  fetchFirewallForwardRules,
  fetchFirewallPendingState,
  fetchFirewallRules,
  previewFirewallRules,
  updateFirewallForwardConfig,
  updateFirewallForwardRule,
  updateFirewallRule,
  type FirewallForwardConfig,
  type FirewallForwardRule,
  type FirewallForwardRulePayload,
  type FirewallPendingState,
  type FirewallRule,
  type FirewallRulePayload,
} from "../api/firewall";

const loading = ref(false);
const saving = ref(false);
const applying = ref(false);
const confirming = ref(false);
const editingId = ref<number | null>(null);
const items = ref<FirewallRule[]>([]);
const defaultsInitialized = ref(false);
const previewText = ref("");
const pendingState = ref<FirewallPendingState | null>(null);

const forwardSaving = ref(false);
const forwardEditingId = ref<number | null>(null);
const forwardItems = ref<FirewallForwardRule[]>([]);
const legacyForwardSaving = ref(false);
const legacyForwardConfig = reactive<FirewallForwardConfig>({
  enabled: false,
  wgInterface: "wg0",
  lanCidr: "",
});

const conntrackLoading = ref(false);
const conntrackSourceIP = ref("");
const conntrackItems = ref<ConntrackEntry[]>([]);
const conntrackAvailable = ref(true);
const conntrackMessage = ref("");

const form = reactive<FirewallRulePayload>({
  name: "",
  kind: "template",
  templateKey: "ssh",
  protocol: "tcp",
  port: 0,
  portRangeStart: 0,
  portRangeEnd: 0,
  sourceCidr: "",
  enabled: true,
  priority: 0,
  description: "",
});

const forwardForm = reactive<FirewallForwardRulePayload>({
  name: "",
  sourceCidr: "10.66.66.0/24",
  destinationCidr: "",
  protocol: "any",
  destinationPort: 0,
  enabled: true,
  priority: 0,
  description: "",
});

const actionLabel = computed(() => (editingId.value ? "更新 input 规则" : "新增 input 规则"));
const isCustom = computed(() => form.kind === "custom");
const forwardActionLabel = computed(() => (forwardEditingId.value ? "更新 forward 规则" : "新增 forward 规则"));
const usesLegacyForwardFallback = computed(() => forwardItems.value.length === 0 && legacyForwardConfig.enabled);
const summaryChips = computed(() => [
  { label: "Input 规则", value: String(items.value.length) },
  { label: "Forward 规则", value: String(forwardItems.value.length) },
  { label: "待确认状态", value: pendingState.value?.pending ? "进行中" : "无" },
]);

function resetForm() {
  editingId.value = null;
  form.name = "";
  form.kind = "template";
  form.templateKey = "ssh";
  form.protocol = "tcp";
  form.port = 0;
  form.portRangeStart = 0;
  form.portRangeEnd = 0;
  form.sourceCidr = "";
  form.enabled = true;
  form.priority = 0;
  form.description = "";
}

function resetForwardForm() {
  forwardEditingId.value = null;
  forwardForm.name = "";
  forwardForm.sourceCidr = "10.66.66.0/24";
  forwardForm.destinationCidr = "";
  forwardForm.protocol = "any";
  forwardForm.destinationPort = 0;
  forwardForm.enabled = true;
  forwardForm.priority = 0;
  forwardForm.description = "";
}

async function load() {
  loading.value = true;
  try {
    const [rules, pending, forwardRules, legacyForward] = await Promise.all([
      fetchFirewallRules(),
      fetchFirewallPendingState(),
      fetchFirewallForwardRules(),
      fetchFirewallForwardConfig(),
    ]);
    items.value = rules.items;
    defaultsInitialized.value = rules.initialized;
    pendingState.value = pending.pendingState;
    forwardItems.value = forwardRules.items;
    legacyForwardConfig.enabled = legacyForward.enabled;
    legacyForwardConfig.wgInterface = legacyForward.wgInterface;
    legacyForwardConfig.lanCidr = legacyForward.lanCidr;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载防火墙配置失败");
  } finally {
    loading.value = false;
  }
}

function startEdit(item: FirewallRule) {
  editingId.value = item.id;
  form.name = item.name;
  form.kind = item.kind;
  form.templateKey = item.templateKey;
  form.protocol = item.protocol;
  form.port = item.port;
  form.portRangeStart = item.portRangeStart;
  form.portRangeEnd = item.portRangeEnd;
  form.sourceCidr = item.sourceCidr;
  form.enabled = item.enabled;
  form.priority = item.priority;
  form.description = item.description;
}

function startForwardEdit(item: FirewallForwardRule) {
  forwardEditingId.value = item.id;
  forwardForm.name = item.name;
  forwardForm.sourceCidr = item.sourceCidr;
  forwardForm.destinationCidr = item.destinationCidr;
  forwardForm.protocol = item.protocol;
  forwardForm.destinationPort = item.destinationPort;
  forwardForm.enabled = item.enabled;
  forwardForm.priority = item.priority;
  forwardForm.description = item.description;
}

async function submit() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateFirewallRule(editingId.value, { ...form });
      ElMessage.success("input 规则已更新");
    } else {
      await createFirewallRule({ ...form });
      ElMessage.success("input 规则已创建");
    }
    resetForm();
    await load();
    await preview();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存 input 规则失败");
  } finally {
    saving.value = false;
  }
}

async function submitForward() {
  forwardSaving.value = true;
  try {
    if (forwardEditingId.value) {
      await updateFirewallForwardRule(forwardEditingId.value, { ...forwardForm });
      ElMessage.success("forward 规则已更新");
    } else {
      await createFirewallForwardRule({ ...forwardForm });
      ElMessage.success("forward 规则已创建");
    }
    resetForwardForm();
    await load();
    await preview();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存 forward 规则失败");
  } finally {
    forwardSaving.value = false;
  }
}

async function remove(item: FirewallRule) {
  await ElMessageBox.confirm(`确定删除 input 规则 "${item.name}" 吗？`, "删除确认", { type: "warning" });
  try {
    await deleteFirewallRule(item.id);
    ElMessage.success("规则已删除");
    if (editingId.value === item.id) {
      resetForm();
    }
    await load();
    await preview();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "删除规则失败");
  }
}

async function removeForward(item: FirewallForwardRule) {
  await ElMessageBox.confirm(`确定删除 forward 规则 "${item.name}" 吗？`, "删除确认", { type: "warning" });
  try {
    await deleteFirewallForwardRule(item.id);
    ElMessage.success("forward 规则已删除");
    if (forwardEditingId.value === item.id) {
      resetForwardForm();
    }
    await load();
    await preview();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "删除 forward 规则失败");
  }
}

async function preview() {
  try {
    const response = await previewFirewallRules();
    previewText.value = response.preview;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "预览规则失败");
  }
}

async function apply() {
  applying.value = true;
  try {
    const response = await applyFirewallRules();
    pendingState.value = response.pendingState;
    previewText.value = response.preview;
    ElMessage.success("规则已应用，等待 30 秒内确认");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "应用规则失败");
  } finally {
    applying.value = false;
  }
}

async function confirm() {
  confirming.value = true;
  try {
    await confirmFirewallRules();
    pendingState.value = null;
    ElMessage.success("规则已确认生效");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "确认规则失败");
  } finally {
    confirming.value = false;
  }
}

async function saveLegacyForward() {
  legacyForwardSaving.value = true;
  try {
    const updated = await updateFirewallForwardConfig({
      enabled: legacyForwardConfig.enabled,
      lanCidr: legacyForwardConfig.lanCidr,
    });
    legacyForwardConfig.enabled = updated.enabled;
    legacyForwardConfig.wgInterface = updated.wgInterface;
    legacyForwardConfig.lanCidr = updated.lanCidr;
    ElMessage.success("兼容 forward 配置已更新");
    await preview();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存兼容 forward 配置失败");
  } finally {
    legacyForwardSaving.value = false;
  }
}

async function loadConntrack() {
  conntrackLoading.value = true;
  try {
    const response = await fetchConntrackEntries(conntrackSourceIP.value);
    conntrackAvailable.value = response.available;
    conntrackMessage.value = response.message;
    conntrackItems.value = response.items;
  } catch (error) {
    conntrackAvailable.value = false;
    conntrackItems.value = [];
    conntrackMessage.value = error instanceof Error ? error.message : "加载当前连接失败";
  } finally {
    conntrackLoading.value = false;
  }
}

onMounted(async () => {
  await load();
  await preview();
  await loadConntrack();
});
</script>

<template>
  <div class="page-shell firewall-page">
    <section class="page-card firewall-hero">
      <div class="page-card__body firewall-hero__body">
        <PageHeader
          title="防火墙"
          description="当前支持 input 规则、独立 forward 规则、规则预览、apply / 确认 / 自动回滚，以及当前连接查看。第一次进入时如果还没有基础规则，系统会先帮你建好一组安全起点。"
        >
          <el-space wrap>
            <el-button :loading="loading" @click="load">刷新规则</el-button>
            <el-button @click="preview">预览</el-button>
            <el-button type="warning" :loading="applying" @click="apply">应用</el-button>
            <el-button v-if="pendingState?.pending" type="success" :loading="confirming" @click="confirm">
              确认生效
            </el-button>
          </el-space>
        </PageHeader>

        <div class="firewall-hero__chips">
          <div v-for="chip in summaryChips" :key="chip.label" class="surface-muted firewall-chip">
            <span>{{ chip.label }}</span>
            <strong>{{ chip.value }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section v-if="defaultsInitialized" class="page-card">
      <div class="page-card__body">
        <el-alert
          type="success"
          :closable="false"
          title="已自动初始化基础防火墙规则"
          description="系统已经写入 SSH、HTTP、HTTPS、WireGuard 四条基础模板规则。它们目前只存在于规则列表里，尚未自动下发；建议你先检查和预览，再决定是否点击“应用”。"
        />
      </div>
    </section>

    <section v-if="pendingState?.pending" class="page-card">
      <div class="page-card__body pending-banner">
        <strong>存在待确认的防火墙变更</strong>
        <span>如果 30 秒内不确认，系统会自动回滚到备份规则。</span>
      </div>
    </section>

    <section class="firewall-grid">
      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="Input 规则"
            description="公网入站规则适合管理 SSH、HTTP、HTTPS、WireGuard 等端口；模板规则用来快速起步，自定义规则处理更细的来源限制。"
          />
          <el-form label-position="top" class="firewall-form">
            <el-form-item label="规则名称">
              <el-input v-model="form.name" placeholder="如 allow-ssh-office" />
            </el-form-item>
            <el-form-item label="规则类型">
              <el-select v-model="form.kind">
                <el-option label="模板规则" value="template" />
                <el-option label="自定义规则" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="form.kind === 'template'" label="模板">
              <el-select v-model="form.templateKey">
                <el-option label="SSH" value="ssh" />
                <el-option label="HTTP" value="http" />
                <el-option label="HTTPS" value="https" />
                <el-option label="WireGuard" value="wireguard" />
              </el-select>
            </el-form-item>
            <template v-else>
              <el-form-item label="协议">
                <el-select v-model="form.protocol">
                  <el-option label="TCP" value="tcp" />
                  <el-option label="UDP" value="udp" />
                  <el-option label="ICMP" value="icmp" />
                  <el-option label="ICMPv6" value="icmpv6" />
                </el-select>
              </el-form-item>
              <el-form-item v-if="isCustom && (form.protocol === 'tcp' || form.protocol === 'udp')" label="单端口">
                <el-input-number v-model="form.port" :min="0" :max="65535" />
              </el-form-item>
              <el-form-item v-if="isCustom && (form.protocol === 'tcp' || form.protocol === 'udp')" label="端口范围（可选）">
                <div class="range-row">
                  <el-input-number v-model="form.portRangeStart" :min="0" :max="65535" />
                  <span>到</span>
                  <el-input-number v-model="form.portRangeEnd" :min="0" :max="65535" />
                </div>
              </el-form-item>
            </template>
            <el-form-item label="来源 IP/CIDR">
              <el-input v-model="form.sourceCidr" placeholder="可留空，或如 192.168.1.0/24" />
            </el-form-item>
            <el-form-item label="优先级">
              <el-input-number v-model="form.priority" :min="0" :max="999" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="form.description" type="textarea" :rows="3" placeholder="可选备注" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="form.enabled" active-text="启用这条规则" />
            </el-form-item>
            <el-space wrap>
              <el-button type="primary" :loading="saving" @click="submit">{{ actionLabel }}</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-space>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="Input 列表"
            description="基础模板和自定义规则都会在这里按优先级展示；高风险场景建议先限制来源 IP，再决定是否 apply。"
          />
          <el-table :data="items">
            <el-table-column prop="priority" label="优先级" width="90" />
            <el-table-column prop="name" label="名称" min-width="160" />
            <el-table-column label="类型" width="120">
              <template #default="{ row }">
                {{ row.kind === "template" ? `模板:${row.templateKey}` : `自定义:${row.protocol}` }}
              </template>
            </el-table-column>
            <el-table-column label="端口" width="120">
              <template #default="{ row }">
                <span v-if="row.port">{{ row.port }}</span>
                <span v-else-if="row.portRangeStart && row.portRangeEnd">{{ row.portRangeStart }}-{{ row.portRangeEnd }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="sourceCidr" label="来源" min-width="160" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? "启用" : "停用" }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-space>
                  <el-button link type="primary" @click="startEdit(row)">编辑</el-button>
                  <el-button link type="danger" @click="remove(row)">删除</el-button>
                </el-space>
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-if="items.length === 0"
            description="当前还没有 input 规则。页面首次加载时会自动补基础模板；如果这里仍为空，请刷新或检查后端日志。"
          />
        </div>
      </article>
    </section>

    <section class="firewall-grid">
      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="Forward 规则"
            description="优先服务 WireGuard 到内网的转发放行场景。先定义来源与目标网段，再细化到协议和目标端口。"
          />
          <el-form label-position="top" class="firewall-form">
            <el-form-item label="规则名称">
              <el-input v-model="forwardForm.name" placeholder="如 wg-lan-https" />
            </el-form-item>
            <el-form-item label="来源 CIDR">
              <el-input v-model="forwardForm.sourceCidr" placeholder="如 10.66.66.0/24" />
            </el-form-item>
            <el-form-item label="目标 CIDR">
              <el-input v-model="forwardForm.destinationCidr" placeholder="如 192.168.1.0/24" />
            </el-form-item>
            <el-form-item label="协议">
              <el-select v-model="forwardForm.protocol">
                <el-option label="ANY" value="any" />
                <el-option label="TCP" value="tcp" />
                <el-option label="UDP" value="udp" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="forwardForm.protocol === 'tcp' || forwardForm.protocol === 'udp'" label="目标端口">
              <el-input-number v-model="forwardForm.destinationPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item label="优先级">
              <el-input-number v-model="forwardForm.priority" :min="0" :max="999" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="forwardForm.description" type="textarea" :rows="3" placeholder="可选备注" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="forwardForm.enabled" active-text="启用这条规则" />
            </el-form-item>
            <el-space wrap>
              <el-button type="primary" :loading="forwardSaving" @click="submitForward">{{ forwardActionLabel }}</el-button>
              <el-button @click="resetForwardForm">重置</el-button>
            </el-space>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="Forward 列表"
            description="如果你已经建立独立 forward 规则，系统会优先使用这里的显式规则，而不是旧的单网段兼容配置。"
          />
          <el-table :data="forwardItems">
            <el-table-column prop="priority" label="优先级" width="90" />
            <el-table-column prop="name" label="名称" min-width="160" />
            <el-table-column prop="sourceCidr" label="来源 CIDR" min-width="160" />
            <el-table-column prop="destinationCidr" label="目标 CIDR" min-width="160" />
            <el-table-column prop="protocol" label="协议" width="100" />
            <el-table-column prop="destinationPort" label="目标端口" width="110" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? "启用" : "停用" }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-space>
                  <el-button link type="primary" @click="startForwardEdit(row)">编辑</el-button>
                  <el-button link type="danger" @click="removeForward(row)">删除</el-button>
                </el-space>
              </template>
            </el-table-column>
          </el-table>

          <el-alert
            v-if="usesLegacyForwardFallback"
            class="legacy-alert"
            title="当前还没有独立 forward 规则，预览和应用时会继续兼容旧的 wg0 -> 内网单网段配置。"
            type="info"
            :closable="false"
          />
        </div>
      </article>
    </section>

    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="兼容 Forward 配置"
          description="保留旧的单网段兼容配置；只有在你还没有创建独立 forward 规则时，它才会参与最终渲染。"
        >
          <el-space>
            <el-tag type="info">接口 {{ legacyForwardConfig.wgInterface || "wg0" }}</el-tag>
            <el-button type="primary" :loading="legacyForwardSaving" @click="saveLegacyForward">保存兼容配置</el-button>
          </el-space>
        </PageHeader>
        <div class="forward-row">
          <el-switch v-model="legacyForwardConfig.enabled" active-text="启用旧的 wg0 -> 内网兼容配置" />
          <el-input v-model="legacyForwardConfig.lanCidr" placeholder="如 192.168.1.0/24" />
        </div>
      </div>
    </section>

    <section class="firewall-preview-grid">
      <article class="page-card">
        <div class="page-card__body">
          <h3 class="section-title">规则预览</h3>
          <p class="section-subtitle">这是最终会写入 nftables 的规则文本，建议每次 apply 前先快速浏览一次。</p>
          <pre class="preview">{{ previewText || "点击“预览”查看最终 nftables 规则文本。" }}</pre>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="当前连接"
            description="只读展示 conntrack 明细，可按来源 IP 过滤，便于判断 WireGuard 客户端和其它来源访问是否符合预期。"
          >
            <div class="conntrack-toolbar">
              <el-input v-model="conntrackSourceIP" placeholder="按来源 IP 过滤，例如 10.66.66.2" clearable />
              <el-button :loading="conntrackLoading" @click="loadConntrack">刷新连接</el-button>
            </div>
          </PageHeader>
          <el-alert
            v-if="!conntrackAvailable"
            class="conntrack-alert"
            type="warning"
            :closable="false"
            :title="conntrackMessage || '当前宿主机未提供 conntrack 命令，连接明细暂不可用。'"
          />
          <el-table v-else :data="conntrackItems">
            <el-table-column prop="protocol" label="协议" width="100" />
            <el-table-column prop="sourceIp" label="来源 IP" min-width="150" />
            <el-table-column prop="sourcePort" label="来源端口" width="110" />
            <el-table-column prop="destinationIp" label="目标 IP" min-width="150" />
            <el-table-column prop="destinationPort" label="目标端口" width="110" />
            <el-table-column prop="state" label="状态" width="140" />
          </el-table>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.firewall-hero {
  background:
    radial-gradient(circle at 90% 10%, rgba(245, 184, 65, 0.2), transparent 20%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(247, 250, 249, 0.97));
}

.firewall-hero__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.firewall-hero__chips {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.firewall-chip {
  padding: 16px 18px;
}

.firewall-chip span {
  display: block;
  color: var(--text-muted);
  font-size: 13px;
}

.firewall-chip strong {
  display: block;
  margin-top: 8px;
  font-size: 24px;
}

.pending-banner {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: #8a5b00;
  background: #fff8e6;
}

.firewall-grid {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 20px;
}

.firewall-preview-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.1fr);
  gap: 20px;
}

.firewall-form {
  margin-top: 18px;
}

.forward-row {
  display: grid;
  grid-template-columns: 360px 1fr;
  gap: 16px;
  align-items: center;
  margin-top: 18px;
}

.conntrack-toolbar {
  display: grid;
  grid-template-columns: minmax(280px, 420px) auto;
  gap: 12px;
  align-items: center;
}

.range-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.legacy-alert,
.conntrack-alert {
  margin-top: 16px;
}

.preview {
  margin: 18px 0 0;
  min-height: 320px;
  padding: 18px;
  overflow: auto;
  white-space: pre-wrap;
  border: 1px solid rgba(22, 152, 142, 0.16);
  border-radius: 16px;
  background: linear-gradient(180deg, #13312f 0%, #0c201f 100%);
  color: #eef7f5;
}

@media (max-width: 1180px) {
  .firewall-grid,
  .firewall-preview-grid,
  .firewall-hero__chips {
    grid-template-columns: 1fr;
  }

  .forward-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .conntrack-toolbar {
    grid-template-columns: 1fr;
  }
}
</style>
