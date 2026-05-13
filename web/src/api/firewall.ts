import client from "./client";

export type FirewallRule = {
  id: number;
  name: string;
  kind: string;
  templateKey: string;
  protocol: string;
  port: number;
  portRangeStart: number;
  portRangeEnd: number;
  sourceCidr: string;
  action: string;
  enabled: boolean;
  priority: number;
  description: string;
};

export type FirewallRulePayload = Omit<FirewallRule, "id" | "action">;

export type FirewallPendingState = {
  pending: boolean;
  backupPath: string;
  expiresAt?: string;
  appliedAt?: string;
  originalPath: string;
};

export async function fetchFirewallRules() {
  const response = await client.get("/firewall/rules");
  return response.data as { items: FirewallRule[] };
}

export async function createFirewallRule(payload: FirewallRulePayload) {
  const response = await client.post("/firewall/rules", payload);
  return response.data as FirewallRule;
}

export async function updateFirewallRule(id: number, payload: FirewallRulePayload) {
  const response = await client.put(`/firewall/rules/${id}`, payload);
  return response.data as FirewallRule;
}

export async function deleteFirewallRule(id: number) {
  await client.delete(`/firewall/rules/${id}`);
}

export async function previewFirewallRules() {
  const response = await client.post("/firewall/preview");
  return response.data as { preview: string };
}

export async function applyFirewallRules() {
  const response = await client.post("/firewall/apply");
  return response.data as { pendingState: FirewallPendingState; preview: string };
}

export async function confirmFirewallRules() {
  const response = await client.post("/firewall/confirm");
  return response.data as { message: string };
}

export async function fetchFirewallPendingState() {
  const response = await client.get("/firewall/pending");
  return response.data as { pendingState: FirewallPendingState | null };
}
