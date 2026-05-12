package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/api"
	"github.com/zmoony/pi-gateway/internal/config"
	"github.com/zmoony/pi-gateway/internal/platform"
	"github.com/zmoony/pi-gateway/internal/platform/ddnsgo"
	"github.com/zmoony/pi-gateway/internal/platform/nginx"
	"github.com/zmoony/pi-gateway/internal/platform/systemstats"
	"github.com/zmoony/pi-gateway/internal/platform/wireguard"
	"github.com/zmoony/pi-gateway/internal/service"
	"github.com/zmoony/pi-gateway/internal/store"
)

type Application struct {
	Config config.Config
	DB     *sql.DB
	Engine *gin.Engine
	Server *http.Server
}

func New() (*Application, error) {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := store.Open(cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	userStore := store.SQLiteUserStore{DB: db}
	peerStore := store.SQLiteWireGuardPeerStore{DB: db}
	proxyStore := store.SQLiteProxyHostStore{DB: db}
	ddnsStore := store.SQLiteDDNSConfigStore{DB: db}

	authService := service.AuthService{
		Users:          userStore,
		JWTSecret:      cfg.JWTSecret,
		SessionTTL:     time.Duration(cfg.SessionTTLHours) * time.Hour,
		DefaultUser:    cfg.DefaultAdminUser,
		DefaultPassRaw: cfg.DefaultAdminPass,
	}
	if err := authService.EnsureDefaultAdmin(context.Background()); err != nil {
		return nil, fmt.Errorf("ensure default admin: %w", err)
	}

	commandRunner := platform.ExecCommandRunner{}
	wireGuardManager := wireguard.SystemManager{Runner: commandRunner}
	qrGenerator := wireguard.Base64Generator{Runner: commandRunner}
	nginxManager := nginx.SystemManager{
		Runner:   commandRunner,
		SitesDir: cfg.NginxSitesDir,
	}
	ddnsManager := ddnsgo.SystemManager{
		ConfigPath:    cfg.DDNSGoConfigPath,
		StatusPath:    cfg.DDNSGoStatusPath,
		ReloadCommand: cfg.DDNSGoReloadCmd,
		Runner:        commandRunner,
	}

	systemService := service.SystemStatsService{
		Collector: systemstats.StaticCollector{},
	}

	wgService := service.WireGuardService{
		Peers:          peerStore,
		Manager:        wireGuardManager,
		QRGenerator:    qrGenerator,
		InterfaceName:  cfg.WireGuardIface,
		ServerEndpoint: cfg.WireGuardEndpoint,
		ServerPublicKey: cfg.WireGuardPubKey,
		DefaultDNS:     cfg.WireGuardDNS,
		SubnetV4:       cfg.WireGuardSubnetV4,
		ServerCIDRV4:   cfg.WireGuardServerV4,
	}
	proxyService := service.ProxyService{
		Hosts:   proxyStore,
		Manager: nginxManager,
	}
	ddnsService := service.DDNSService{
		Configs: ddnsStore,
		Manager: ddnsManager,
	}

	engine := api.NewRouter(api.Dependencies{
		Config:         cfg,
		AuthService:    authService,
		SystemService:  systemService,
		WGService:      wgService,
		ProxyService:   proxyService,
		DDNSService:    ddnsService,
	})

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &Application{
		Config: cfg,
		DB:     db,
		Engine: engine,
		Server: server,
	}, nil
}

func (a *Application) Run() error {
	return a.Server.ListenAndServe()
}
