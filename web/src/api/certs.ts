import client from "./client";

export type CertificateConfig = {
  id: number;
  rootDomain: string;
  provider: string;
  accessKeyId: string;
  installDir: string;
  enabledAutoRenew: boolean;
  fullchainPath: string;
  privateKeyPath: string;
  lastIssueStatus: string;
  lastIssueError: string;
  lastIssuedAt?: string;
  notAfter?: string;
  daysRemaining: number;
};

export type CertificatePayload = {
  rootDomain: string;
  provider: string;
  accessKeyId: string;
  accessKeySecret: string;
  installDir: string;
  enabledAutoRenew: boolean;
};

export async function fetchCertificates() {
  const response = await client.get("/certs");
  return response.data as { items: CertificateConfig[] };
}

export async function createCertificate(payload: CertificatePayload) {
  const response = await client.post("/certs", payload);
  return response.data as CertificateConfig;
}

export async function updateCertificate(id: number, payload: CertificatePayload) {
  const response = await client.put(`/certs/${id}`, payload);
  return response.data as CertificateConfig;
}

export async function renewCertificate(id: number) {
  const response = await client.post(`/certs/${id}/renew`);
  return response.data as CertificateConfig;
}
