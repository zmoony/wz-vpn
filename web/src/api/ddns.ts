import client from "./client";

export type DDNSConfig = {
  id: number;
  provider: string;
  accessKeyId: string;
  domain: string;
  subdomain: string;
  enabled: boolean;
  checkIntervalSeconds: number;
  lastKnownIpv6: string;
  lastStatus: string;
  lastError: string;
  lastSyncedAt?: string;
};

export type DDNSPayload = {
  provider: string;
  accessKeyId: string;
  accessKeySecret: string;
  domain: string;
  subdomain: string;
  enabled: boolean;
  checkIntervalSeconds: number;
};

export type DDNSRuntimeStatus = {
  lastKnownIpv6: string;
  lastStatus: string;
  lastError: string;
  lastSyncedAt?: string;
};

export async function fetchDDNSConfigs() {
  const response = await client.get("/ddns");
  return response.data as { items: DDNSConfig[]; runtimeStatus: DDNSRuntimeStatus };
}

export async function createDDNSConfig(payload: DDNSPayload) {
  const response = await client.post("/ddns", payload);
  return response.data as DDNSConfig;
}

export async function updateDDNSConfig(id: number, payload: DDNSPayload) {
  const response = await client.put(`/ddns/${id}`, payload);
  return response.data as DDNSConfig;
}

export async function syncDDNSConfigs() {
  const response = await client.post("/ddns/sync");
  return response.data as { message: string };
}
