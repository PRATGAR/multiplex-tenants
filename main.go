package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"tenantly/internal/config"
	"tenantly/internal/exporter"
	"tenantly/internal/handlers"
	"tenantly/internal/middleware"
	"tenantly/internal/server"
	"tenantly/internal/service"
	"tenantly/internal/store"
	"tenantly/internal/verifier"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	tenantsRepo := store.NewTenants(db)
	usersRepo := store.NewUsers(db)
	projectsRepo := store.NewProjects(db)
	invitesRepo := store.NewInvitations(db)
	domainsRepo := store.NewDomains(db)
	exportsRepo := store.NewExports(db)
	billingRepo := store.NewBilling(db)
	auditRepo := store.NewAudit(db)
	files := exporter.New(cfg.ExportDir)

	ttl := time.Duration(cfg.SessionTTLHours) * time.Hour
	authSvc := service.NewAuth(usersRepo, rdb, cfg.PasswordSalt, ttl)
	tenantsSvc := service.NewTenants(tenantsRepo, usersRepo, files, cfg.PasswordSalt)
	usersSvc := service.NewUsers(usersRepo, auditRepo)
	projectsSvc := service.NewProjects(projectsRepo)
	invitesSvc := service.NewInvitations(invitesRepo, usersRepo, auditRepo, rdb, cfg.PasswordSalt)
	billingSvc := service.NewBilling(billingRepo, tenantsRepo)
	domainsSvc := service.NewDomains(domainsRepo, verifier.New())
	exportsSvc := service.NewExports(tenantsRepo, usersRepo, exportsRepo, files)

	router := server.NewRouter(server.Handlers{
		Auth:        handlers.NewAuth(authSvc),
		Tenants:     handlers.NewTenants(tenantsSvc),
		Projects:    handlers.NewProjects(projectsSvc),
		Users:       handlers.NewUsers(usersSvc),
		Invitations: handlers.NewInvitations(invitesSvc),
		Billing:     handlers.NewBilling(billingSvc),
		Domains:     handlers.NewDomains(domainsSvc),
		Exports:     handlers.NewExports(exportsSvc),
	}, middleware.NewGuard(rdb), cfg.ExportDir)

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on %s", cfg.ListenAddr)
	log.Fatal(srv.ListenAndServe())
}
