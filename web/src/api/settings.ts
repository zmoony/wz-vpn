import client from "./client";

export type WireGuardSettings = {
  interface: string;
  subnetV4: string;
  serverV4: string;
  dns: string;
  endpoint: string;
  publicKey: string;
};

export type RuntimeSettings = {
  appEnv: string;
  listenAddr: string;
  dataDir: string;
  databasePath: string;
  webDistDir: string;
  nginxSitesDir: string;
  ddnsGoConfigPath: string;
  certsInstallRoot: string;
  nftablesRulesPath: string;
  firewallPendingSeconds: number;
};

export type SettingsView = {
  wireGuard: WireGuardSettings;
  runtime: RuntimeSettings;
};

export async function fetchSettings() {
  const response = await client.get("/settings");
  return response.data as SettingsView;
}

export async function updateSettings(payload: { wireGuard: WireGuardSettings }) {
  const response = await client.put("/settings", payload);
  return response.data as SettingsView;
}
