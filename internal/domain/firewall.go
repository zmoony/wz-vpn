package domain

import "time"

type FirewallRule struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Kind           string    `json:"kind"`
	TemplateKey    string    `json:"templateKey"`
	Protocol       string    `json:"protocol"`
	Port           int       `json:"port"`
	PortRangeStart int       `json:"portRangeStart"`
	PortRangeEnd   int       `json:"portRangeEnd"`
	SourceCIDR     string    `json:"sourceCidr"`
	Action         string    `json:"action"`
	Enabled        bool      `json:"enabled"`
	Priority       int       `json:"priority"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type FirewallPendingState struct {
	Pending       bool       `json:"pending"`
	BackupPath    string     `json:"backupPath"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	AppliedAt     *time.Time `json:"appliedAt,omitempty"`
	OriginalPath  string     `json:"originalPath"`
}

type FirewallForwardConfig struct {
	Enabled    bool   `json:"enabled"`
	WGInterface string `json:"wgInterface"`
	LanCIDR    string `json:"lanCidr"`
}

type FirewallForwardRule struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	SourceCIDR      string    `json:"sourceCidr"`
	DestinationCIDR string    `json:"destinationCidr"`
	Protocol        string    `json:"protocol"`
	DestinationPort int       `json:"destinationPort"`
	Enabled         bool      `json:"enabled"`
	Priority        int       `json:"priority"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
