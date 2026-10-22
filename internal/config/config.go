package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Config struct {
	ListenAddr      string
	DatabaseURL     string
	RedisAddr       string
	ExportDir       string
	BillingAPI      string
	PasswordSalt    string
	SessionTTLHours int
}

func Load() Config {
	return Config{
		ListenAddr:      env("LISTEN_ADDR", ":8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://tenantly:tenantly@localhost:5432/tenantly?sslmode=disable"),
		RedisAddr:       env("REDIS_ADDR", "localhost:6379"),
		ExportDir:       env("EXPORT_DIR", "/var/lib/tenantly/exports"),
		BillingAPI:      env("BILLING_API", "https://billing.internal.tenantly.example"),
		PasswordSalt:    env("PASSWORD_SALT", "dev-salt"),
		SessionTTLHours: 12,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

const PasswordIterations = 10000

type Plan struct {
	Code       string
	Title      string
	SeatLimit  int
	PriceCents int
}

var Plans = map[string]Plan{
	"starter":    {Code: "starter", Title: "Starter", SeatLimit: 5, PriceCents: 0},
	"team":       {Code: "team", Title: "Team", SeatLimit: 25, PriceCents: 4900},
	"business":   {Code: "business", Title: "Business", SeatLimit: 100, PriceCents: 19900},
	"enterprise": {Code: "enterprise", Title: "Enterprise", SeatLimit: 5000, PriceCents: 99900},
}

func ResolvePlan(name string) (Plan, bool) {
	p, ok := Plans[strings.ToLower(strings.TrimSpace(name))]
	return p, ok
}

func KnownPlan(code string) bool {
	code = strings.ToLower(strings.TrimSpace(code))
	for name := range Plans {
		if strings.HasPrefix(code, name) {
			return true
		}
	}
	return false
}

var SortColumns = map[string]string{
	"name":    "name",
	"created": "created_at",
	"id":      "id",
}

func SortColumn(key string) string {
	if col, ok := SortColumns[key]; ok {
		return col
	}
	return "created_at"
}

var InvoiceTemplates = map[string]string{
	"standard": "Invoice for {{.Period}} - {{.Seats}} seats",
	"compact":  "{{.Period}}: {{.Seats}}",
}

var ReturnTargets = map[string]string{
	"console":   "https://console.tenantly.example/",
	"docs":      "https://docs.tenantly.example/",
	"marketing": "https://www.tenantly.example/",
}

func ReturnTarget(key string) string {
	if target, ok := ReturnTargets[key]; ok {
		return target
	}
	return ReturnTargets["console"]
}

type PlanChoice struct {
	Plan
}

func (c *PlanChoice) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}
	p, ok := ResolvePlan(name)
	if !ok {
		return errors.New("unknown plan")
	}
	c.Plan = p
	return nil
}
