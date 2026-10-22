package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tenantly/internal/config"
	"tenantly/internal/handlers"
	"tenantly/internal/middleware"
	"tenantly/internal/pathutil"
	"tenantly/internal/validate"
)

type Handlers struct {
	Auth        *handlers.Auth
	Tenants     *handlers.Tenants
	Projects    *handlers.Projects
	Users       *handlers.Users
	Invitations *handlers.Invitations
	Billing     *handlers.Billing
	Domains     *handlers.Domains
	Exports     *handlers.Exports
}

func NewRouter(h Handlers, guard *middleware.Guard, exportDir string) http.Handler {
	root := mux.NewRouter()
	api := root.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.RequestLog)

	api.HandleFunc("/auth/login", h.Auth.Login).Methods(http.MethodPost)
	api.Handle("/auth/callback",
		middleware.QueryRewrite("next", validate.LocalPathOrRoot)(http.HandlerFunc(h.Auth.Callback))).Methods(http.MethodGet)
	api.Handle("/auth/logout",
		middleware.QueryRewrite("return", config.ReturnTarget)(http.HandlerFunc(h.Auth.Logout))).Methods(http.MethodGet)
	api.HandleFunc("/tenants", h.Tenants.Create).Methods(http.MethodPost)
	api.HandleFunc("/invitations/accept", h.Invitations.Accept).Methods(http.MethodPost)
	api.Handle("/invitations/{token}/decline",
		middleware.PathCheck("token", validate.InviteToken)(http.HandlerFunc(h.Invitations.Decline))).Methods(http.MethodPost)

	authed := api.PathPrefix("").Subrouter()
	authed.Use(guard.Authenticate)
	authed.HandleFunc("/me/projects", h.Projects.Mine).Methods(http.MethodGet)
	authed.Handle("/projects",
		middleware.QueryCheck("q", validate.SearchTerm)(
			middleware.QueryRewrite("sort", config.SortColumn)(http.HandlerFunc(h.Projects.Search)))).Methods(http.MethodGet)
	authed.HandleFunc("/tenants/{tid}/projects", h.Projects.ForTenant).Methods(http.MethodGet)
	authed.HandleFunc("/tenants/{tid:[0-9]+}/billing/summary", h.Billing.Summary).Methods(http.MethodGet)
	authed.Handle("/billing/plans/{plan}/price",
		middleware.PathCheck("plan", config.KnownPlan)(http.HandlerFunc(h.Billing.Quote))).Methods(http.MethodGet)
	authed.Handle("/exports/download",
		middleware.QueryRewrite("file", func(rel string) string {
			return pathutil.SafeJoin(exportDir, rel)
		})(http.HandlerFunc(h.Exports.Download))).Methods(http.MethodGet)
	authed.Handle("/exports/{eid:[0-9]+}/compress",
		middleware.RequireRole("admin", "owner")(
			middleware.QueryCheck("level", validate.GzipLevel)(http.HandlerFunc(h.Exports.Compress)))).Methods(http.MethodPost)

	scoped := authed.PathPrefix("/tenants/{tid:[0-9]+}").Subrouter()
	scoped.Use(middleware.TenantScope)
	scoped.Handle("",
		middleware.RequireRole("admin", "owner")(http.HandlerFunc(h.Tenants.Update))).Methods(http.MethodPatch)
	scoped.HandleFunc("/welcome", h.Tenants.Welcome).Methods(http.MethodGet)
	scoped.HandleFunc("/users/{uid:[0-9]+}/activity", h.Users.Activity).Methods(http.MethodGet)
	scoped.Handle("/invitations",
		middleware.RequireRole("admin", "owner")(http.HandlerFunc(h.Invitations.Create))).Methods(http.MethodPost)
	scoped.HandleFunc("/billing/invoice", h.Billing.Invoice).Methods(http.MethodGet)
	scoped.Handle("/plan",
		middleware.RequireRole("owner")(http.HandlerFunc(h.Billing.ChangePlan))).Methods(http.MethodPut)
	scoped.HandleFunc("/invoices/preview", h.Billing.Preview).Methods(http.MethodGet)
	scoped.Handle("/domains",
		middleware.RequireRole("admin", "owner")(http.HandlerFunc(h.Domains.Add))).Methods(http.MethodPost)
	scoped.Handle("/domains/{did:[0-9]+}/verify",
		middleware.RequireRole("admin", "owner")(http.HandlerFunc(h.Domains.Verify))).Methods(http.MethodPost)
	scoped.Handle("/exports",
		middleware.RequireRole("admin", "owner")(http.HandlerFunc(h.Exports.Create))).Methods(http.MethodPost)
	scoped.Handle("/exports/{eid:[0-9]+}/file",
		middleware.CSVAttachment(http.HandlerFunc(h.Exports.File))).Methods(http.MethodGet)

	return root
}
