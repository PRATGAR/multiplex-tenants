package model

import "time"

type Session struct {
	UserID   int64  `json:"user_id"`
	TenantID int64  `json:"tenant_id"`
	Role     string `json:"role"`
}

type Tenant struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Plan          string    `json:"plan"`
	InvoiceFooter string    `json:"invoice_footer"`
	CreatedAt     time.Time `json:"created_at"`
}

type User struct {
	ID           int64  `json:"id"`
	TenantID     int64  `json:"tenant_id"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

type Project struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Invitation struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type InvitePayload struct {
	TenantID int64  `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Domain struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	VerifyToken string `json:"verify_token"`
	Status      string `json:"status"`
}

type Export struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Rows      int       `json:"rows"`
	CreatedAt time.Time `json:"created_at"`
}

type PlanPrice struct {
	Code      string `json:"code"`
	UnitCents int    `json:"unit_cents"`
	SeatCents int    `json:"seat_cents"`
}

type BillingSummary struct {
	Seats        int    `json:"seats"`
	BalanceCents int64  `json:"balance_cents"`
	Currency     string `json:"currency"`
}

type AuditEntry struct {
	Action string    `json:"action"`
	At     time.Time `json:"at"`
}

type Invoice struct {
	Plan      string `json:"plan"`
	Seats     int    `json:"seats"`
	TotalCent int    `json:"total_cents"`
}
