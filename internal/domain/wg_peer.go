package domain

import "time"

type WireGuardPeer struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	ClientIPv4          string     `json:"clientIpv4"`
	ClientIPv6          string     `json:"clientIpv6"`
	PublicKey           string     `json:"publicKey"`
	PrivateKeyEncrypted string     `json:"-"`
	PresharedKeyEnc     string     `json:"-"`
	DNS                 string     `json:"dns"`
	AllowedIPs          string     `json:"allowedIps"`
	Endpoint            string     `json:"endpoint"`
	PersistentKeepalive int        `json:"persistentKeepalive"`
	Enabled             bool       `json:"enabled"`
	Description         string     `json:"description"`
	LastHandshakeAt     *time.Time `json:"lastHandshakeAt,omitempty"`
	RxBytes             int64      `json:"rxBytes"`
	TxBytes             int64      `json:"txBytes"`
	Online              bool       `json:"online"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

type WireGuardPeerSpec struct {
	Name                string `json:"name"`
	ClientIPv4          string `json:"clientIpv4"`
	ClientIPv6          string `json:"clientIpv6"`
	PublicKey           string `json:"publicKey"`
	PrivateKey          string `json:"privateKey,omitempty"`
	PresharedKey        string `json:"presharedKey,omitempty"`
	DNS                 string `json:"dns"`
	AllowedIPs          string `json:"allowedIps"`
	Endpoint            string `json:"endpoint"`
	PersistentKeepalive int    `json:"persistentKeepalive"`
	Enabled             bool   `json:"enabled"`
	Description         string `json:"description"`
}

type WireGuardRuntimePeer struct {
	PublicKey       string
	LastHandshakeAt *time.Time
	RxBytes         int64
	TxBytes         int64
	Online          bool
}
