import client from "./client";

export type WireGuardPeer = {
  id: number;
  name: string;
  clientIpv4: string;
  publicKey: string;
  dns: string;
  allowedIps: string;
  endpoint: string;
  persistentKeepalive: number;
  enabled: boolean;
  description: string;
  lastHandshakeAt?: string;
  rxBytes: number;
  txBytes: number;
  online: boolean;
};

export async function fetchPeers() {
  const response = await client.get("/wg/peers");
  return response.data as { items: WireGuardPeer[] };
}

export async function createPeer(payload: { name: string; description: string }) {
  const response = await client.post("/wg/peers", payload);
  return response.data as {
    peer: WireGuardPeer;
    clientConf: string;
    qrCode: string;
  };
}

export async function togglePeer(id: number) {
  const response = await client.post(`/wg/peers/${id}/toggle`);
  return response.data as WireGuardPeer;
}

export async function deletePeer(id: number) {
  await client.delete(`/wg/peers/${id}`);
}
