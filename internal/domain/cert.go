package domain

import "time"

type CertificateConfig struct {
	ID                 int64      `json:"id"`
	RootDomain         string     `json:"rootDomain"`
	Provider           string     `json:"provider"`
	AccessKeyID        string     `json:"accessKeyId"`
	AccessKeySecretEnc string     `json:"-"`
	InstallDir         string     `json:"installDir"`
	EnabledAutoRenew   bool       `json:"enabledAutoRenew"`
	FullchainPath      string     `json:"fullchainPath"`
	PrivateKeyPath     string     `json:"privateKeyPath"`
	LastIssueStatus    string     `json:"lastIssueStatus"`
	LastIssueError     string     `json:"lastIssueError"`
	LastIssuedAt       *time.Time `json:"lastIssuedAt,omitempty"`
	NotAfter           *time.Time `json:"notAfter,omitempty"`
	DaysRemaining      int        `json:"daysRemaining"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}
