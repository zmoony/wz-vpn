package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/api/middleware"
	"github.com/zmoony/pi-gateway/internal/config"
	"github.com/zmoony/pi-gateway/internal/service"
)

type Dependencies struct {
	Config        config.Config
	AuthService   service.AuthService
	SystemService service.SystemStatsService
	WGService     service.WireGuardService
	ProxyService  service.ProxyService
	DDNSService   service.DDNSService
	CertService   service.CertificateService
	FirewallService *service.FirewallService
}

func NewRouter(deps Dependencies) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery())

	engine.GET("/api/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := AuthHandler{Service: deps.AuthService, Config: deps.Config}
	sysHandler := SystemHandler{Service: deps.SystemService}
	wgHandler := WireGuardHandler{Service: deps.WGService}
	proxyHandler := ProxyHandler{Service: deps.ProxyService}
	ddnsHandler := DDNSHandler{Service: deps.DDNSService}
	certHandler := CertificateHandler{Service: deps.CertService}
	firewallHandler := FirewallHandler{Service: deps.FirewallService}
	moduleHandler := ModuleStubHandler{}

	apiGroup := engine.Group("/api")
	{
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/me", middleware.Auth(deps.AuthService), authHandler.Me)
		}

		protected := apiGroup.Group("/")
		protected.Use(middleware.Auth(deps.AuthService))
		{
			protected.GET("/sys/stats", sysHandler.Get)
			protected.GET("/wg/peers", wgHandler.List)
			protected.POST("/wg/peers", wgHandler.Create)
			protected.POST("/wg/peers/:id/toggle", wgHandler.Toggle)
			protected.DELETE("/wg/peers/:id", wgHandler.Delete)

			protected.GET("/proxy", proxyHandler.List)
			protected.POST("/proxy", proxyHandler.Create)
			protected.PUT("/proxy/:id", proxyHandler.Update)
			protected.DELETE("/proxy/:id", proxyHandler.Delete)

			protected.GET("/ddns", ddnsHandler.List)
			protected.POST("/ddns", ddnsHandler.Create)
			protected.PUT("/ddns/:id", ddnsHandler.Update)
			protected.POST("/ddns/sync", ddnsHandler.Sync)
			protected.GET("/certs", certHandler.List)
			protected.POST("/certs", certHandler.Create)
			protected.PUT("/certs/:id", certHandler.Update)
			protected.POST("/certs/:id/renew", certHandler.Renew)
			protected.GET("/firewall/rules", firewallHandler.List)
			protected.POST("/firewall/rules", firewallHandler.Create)
			protected.PUT("/firewall/rules/:id", firewallHandler.Update)
			protected.DELETE("/firewall/rules/:id", firewallHandler.Delete)
			protected.POST("/firewall/preview", firewallHandler.Preview)
			protected.POST("/firewall/apply", firewallHandler.Apply)
			protected.POST("/firewall/confirm", firewallHandler.Confirm)
			protected.GET("/firewall/pending", firewallHandler.Pending)
			protected.GET("/settings", moduleHandler.List("settings"))
		}
	}

	registerWeb(engine, deps.Config.WebDistDir)
	return engine
}

func registerWeb(engine *gin.Engine, distDir string) {
	if stat, err := os.Stat(distDir); err == nil && stat.IsDir() {
		engine.Static("/assets", distDir+"/assets")
		engine.GET("/", func(c *gin.Context) {
			c.File(distDir + "/index.html")
		})
		engine.NoRoute(func(c *gin.Context) {
			if c.Request.Method != http.MethodGet || len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.File(distDir + "/index.html")
		})
		return
	}

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Pi-Gateway backend is running. Build the frontend into web/dist to serve the SPA.",
		})
	})
}
