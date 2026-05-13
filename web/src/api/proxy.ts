import client from "./client";

export type ProxyHost = {
  id: number;
  name: string;
  serverName: string;
  upstreamUrl: string;
  certificateRootDomain: string;
  certificateCertPath: string;
  certificateKeyPath: string;
  enabled: boolean;
  description: string;
  lastApplyStatus: string;
  lastApplyError: string;
  lastAppliedAt?: string;
};

export type ProxyPayload = {
  name: string;
  serverName: string;
  upstreamUrl: string;
  certificateRootDomain: string;
  enabled: boolean;
  description: string;
};

export async function fetchProxyHosts() {
  const response = await client.get("/proxy");
  return response.data as { items: ProxyHost[] };
}

export async function createProxyHost(payload: ProxyPayload) {
  const response = await client.post("/proxy", payload);
  return response.data as ProxyHost;
}

export async function updateProxyHost(id: number, payload: ProxyPayload) {
  const response = await client.put(`/proxy/${id}`, payload);
  return response.data as ProxyHost;
}

export async function deleteProxyHost(id: number) {
  await client.delete(`/proxy/${id}`);
}
