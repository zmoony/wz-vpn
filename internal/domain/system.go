package domain

type SystemStats struct {
	CPUPercent    float64           `json:"cpuPercent"`
	Load1         float64           `json:"load1"`
	Load5         float64           `json:"load5"`
	Load15        float64           `json:"load15"`
	MemoryPercent float64           `json:"memoryPercent"`
	MemoryUsedMB  uint64            `json:"memoryUsedMb"`
	MemoryTotalMB uint64            `json:"memoryTotalMb"`
	TemperatureC  float64           `json:"temperatureC"`
	DiskPercent   float64           `json:"diskPercent"`
	DiskUsedGB    float64           `json:"diskUsedGb"`
	DiskTotalGB   float64           `json:"diskTotalGb"`
	Network       []NetworkCounters `json:"network"`
}

type NetworkCounters struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rxBytes"`
	TxBytes uint64 `json:"txBytes"`
}
