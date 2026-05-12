package domain

import "time"

type DDNSConfig struct {
	ID                   int64      `json:"id"`
	Provider             string     `json:"provider"`
	AccessKeyID          string     `json:"accessKeyId"`
	AccessKeySecretEnc   string     `json:"-"`
	Domain               string     `json:"domain"`
	Subdomain            string     `json:"subdomain"`
	Enabled              bool       `json:"enabled"`
	CheckIntervalSeconds int        `json:"checkIntervalSeconds"`
	LastKnownIPv6        string     `json:"lastKnownIpv6"`
	LastStatus           string     `json:"lastStatus"`
	LastError            string     `json:"lastError"`
	LastSyncedAt         *time.Time `json:"lastSyncedAt,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
}

type DDNSRuntimeStatus struct {
	LastKnownIPv6 string     `json:"lastKnownIpv6"`
	LastStatus    string     `json:"lastStatus"`
	LastError     string     `json:"lastError"`
	LastSyncedAt  *time.Time `json:"lastSyncedAt,omitempty"`
}
