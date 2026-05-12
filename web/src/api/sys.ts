import client from "./client";

export type SystemStats = {
  cpuPercent: number;
  load1: number;
  load5: number;
  load15: number;
  memoryPercent: number;
  memoryUsedMb: number;
  memoryTotalMb: number;
  temperatureC: number;
  diskPercent: number;
  diskUsedGb: number;
  diskTotalGb: number;
  network: Array<{ name: string; rxBytes: number; txBytes: number }>;
};

export async function fetchSystemStats() {
  const response = await client.get("/sys/stats");
  return response.data as SystemStats;
}
