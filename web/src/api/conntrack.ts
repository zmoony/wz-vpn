import client from "./client";

export type ConntrackEntry = {
  protocol: string;
  sourceIp: string;
  sourcePort: string;
  destinationIp: string;
  destinationPort: string;
  state: string;
};

export async function fetchConntrackEntries(sourceIp = "") {
  const response = await client.get("/firewall/conntrack", {
    params: sourceIp ? { source_ip: sourceIp } : {},
  });
  return response.data as { available: boolean; message: string; items: ConntrackEntry[] };
}
