package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv            string
	ListenAddr        string
	DataDir           string
	DatabasePath      string
	WebDistDir        string
	JWTSecret         string
	SessionTTLHours   int
	DefaultAdminUser  string
	DefaultAdminPass  string
	WireGuardIface    string
	WireGuardSubnetV4 string
	WireGuardServerV4 string
	WireGuardDNS      string
	WireGuardEndpoint string
	WireGuardPubKey   string
	NginxSitesDir     string
	DDNSGoConfigPath  string
	DDNSGoStatusPath  string
	DDNSGoReloadCmd   string
	ACMEShPath        string
	CertsInstallRoot  string
	NginxReloadCmd    string
	NftablesRulesPath string
	FirewallStatePath string
	FirewallBackupsDir string
	FirewallPendingSeconds int
}

func Load() Config {
	dataDir := getenv("PI_GATEWAY_DATA_DIR", "data")

	return Config{
		AppEnv:            getenv("PI_GATEWAY_ENV", "development"),
		ListenAddr:        getenv("PI_GATEWAY_LISTEN_ADDR", ":8443"),
		DataDir:           dataDir,
		DatabasePath:      getenv("PI_GATEWAY_DB_PATH", dataDir+"/app.db"),
		WebDistDir:        getenv("PI_GATEWAY_WEB_DIST", "web/dist"),
		JWTSecret:         getenv("PI_GATEWAY_JWT_SECRET", "change-me-before-production"),
		SessionTTLHours:   getenvInt("PI_GATEWAY_SESSION_TTL_HOURS", 168),
		DefaultAdminUser:  getenv("PI_GATEWAY_ADMIN_USER", "admin"),
		DefaultAdminPass:  getenv("PI_GATEWAY_ADMIN_PASS", "admin123456"),
		WireGuardIface:    getenv("PI_GATEWAY_WG_IFACE", "wg0"),
		WireGuardSubnetV4: getenv("PI_GATEWAY_WG_SUBNET_V4", "10.66.66.0/24"),
		WireGuardServerV4: getenv("PI_GATEWAY_WG_SERVER_V4", "10.66.66.1/24"),
		WireGuardDNS:      getenv("PI_GATEWAY_WG_DNS", "192.168.1.1"),
		WireGuardEndpoint: getenv("PI_GATEWAY_WG_ENDPOINT", "vpn.example.com:51820"),
		WireGuardPubKey:   getenv("PI_GATEWAY_WG_PUBLIC_KEY", "replace-with-server-public-key"),
		NginxSitesDir:     getenv("PI_GATEWAY_NGINX_SITES_DIR", dataDir+"/nginx/sites-enabled"),
		DDNSGoConfigPath:  getenv("PI_GATEWAY_DDNSGO_CONFIG_PATH", dataDir+"/ddns-go/config.json"),
		DDNSGoStatusPath:  getenv("PI_GATEWAY_DDNSGO_STATUS_PATH", dataDir+"/ddns-go/status.json"),
		DDNSGoReloadCmd:   getenv("PI_GATEWAY_DDNSGO_RELOAD_COMMAND", ""),
		ACMEShPath:        getenv("PI_GATEWAY_ACME_SH_PATH", "~/.acme.sh/acme.sh"),
		CertsInstallRoot:  getenv("PI_GATEWAY_CERTS_INSTALL_ROOT", dataDir+"/certs"),
		NginxReloadCmd:    getenv("PI_GATEWAY_NGINX_RELOAD_COMMAND", "nginx -s reload"),
		NftablesRulesPath: getenv("PI_GATEWAY_NFTABLES_RULES_PATH", dataDir+"/nftables/pi-gateway.nft"),
		FirewallStatePath: getenv("PI_GATEWAY_FIREWALL_STATE_PATH", dataDir+"/nftables/pending-apply.json"),
		FirewallBackupsDir: getenv("PI_GATEWAY_FIREWALL_BACKUPS_DIR", dataDir+"/nftables/backups"),
		FirewallPendingSeconds: getenvInt("PI_GATEWAY_FIREWALL_PENDING_SECONDS", 30),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
