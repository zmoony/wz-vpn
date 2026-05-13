<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

import PageHeader from "../components/PageHeader.vue";
import {
  applyFirewallRules,
  confirmFirewallRules,
  createFirewallRule,
  deleteFirewallRule,
  fetchFirewallPendingState,
  fetchFirewallRules,
  previewFirewallRules,
  updateFirewallRule,
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
const previewText = ref("");
const pendingState = ref<FirewallPendingState | null>(null);

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

const actionLabel = computed(() => (editingId.value ? "更新规则" : "新增规则"));
const isCustom = computed(() => form.kind === "custom");

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

async function load() {
  loading.value = true;
  try {
    const [rules, pending] = await Promise.all([fetchFirewallRules(), fetchFirewallPendingState()]);
    items.value = rules.items;
    pendingState.value = pending.pendingState;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载防火墙规则失败");
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

async function submit() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateFirewallRule(editingId.value, { ...form });
      ElMessage.success("防火墙规则已更新");
    } else {
      await createFirewallRule({ ...form });
      ElMessage.success("防火墙规则已创建");
    }
    resetForm();
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存防火墙规则失败");
  } finally {
    saving.value = false;
  }
}

async function remove(item: FirewallRule) {
  await ElMessageBox.confirm(`确定删除规则 "${item.name}" 吗？`, "删除确认", { type: "warning" });
  try {
    await deleteFirewallRule(item.id);
    ElMessage.success("规则已删除");
    if (editingId.value === item.id) {
      resetForm();
    }
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "删除规则失败");
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

onMounted(async () => {
  await load();
  await preview();
});
</script>

<template>
  <div class="page-shell firewall-page">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="防火墙"
          description="这轮先只管理 nftables 的 input 链，并带真实 apply、确认和 30 秒自动回滚。"
        >
          <el-space>
            <el-button :loading="loading" @click="load">刷新规则</el-button>
            <el-button @click="preview">预览</el-button>
            <el-button type="warning" :loading="applying" @click="apply">应用</el-button>
            <el-button
              v-if="pendingState?.pending"
              type="success"
              :loading="confirming"
              @click="confirm"
            >
              确认生效
            </el-button>
          </el-space>
        </PageHeader>
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
          <h3>{{ actionLabel }}</h3>
          <el-form label-position="top">
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
            <el-space>
              <el-button type="primary" :loading="saving" @click="submit">{{ actionLabel }}</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-space>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
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

          <el-divider />

          <h3>规则预览</h3>
          <pre class="preview">{{ previewText || "点击“预览”查看最终 nftables 规则文本。" }}</pre>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.firewall-grid {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 20px;
}

.pending-banner {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: #8a5b00;
  background: #fff8e6;
}

.range-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.preview {
  margin: 0;
  min-height: 240px;
  padding: 16px;
  overflow: auto;
  white-space: pre-wrap;
  border-radius: 14px;
  background: #0f172a;
  color: #edf5ff;
}

@media (max-width: 1200px) {
  .firewall-grid {
    grid-template-columns: 1fr;
  }
}
</style>
