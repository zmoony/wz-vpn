package domain

import "time"

type ProxyHost struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	ServerName          string     `json:"serverName"`
	UpstreamURL         string     `json:"upstreamUrl"`
	CertificateRootDomain string   `json:"certificateRootDomain"`
	CertificateCertPath string     `json:"certificateCertPath"`
	CertificateKeyPath  string     `json:"certificateKeyPath"`
	Enabled             bool       `json:"enabled"`
	Description         string     `json:"description"`
	LastApplyStatus     string     `json:"lastApplyStatus"`
	LastApplyError      string     `json:"lastApplyError"`
	LastAppliedAt       *time.Time `json:"lastAppliedAt,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}
