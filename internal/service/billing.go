package service

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"tenantly/internal/config"
	"tenantly/internal/model"
	"tenantly/internal/store"
)

type Billing struct {
	repo    *store.Billing
	tenants *store.Tenants
}

func NewBilling(repo *store.Billing, tenants *store.Tenants) *Billing {
	return &Billing{repo: repo, tenants: tenants}
}

func (s *Billing) Summary(ctx context.Context, tenantRef string) (model.BillingSummary, error) {
	return s.repo.Summary(ctx, tenantRef)
}

func (s *Billing) Invoice(ctx context.Context, tenantID int64) (model.Invoice, error) {
	t, err := s.tenants.Get(ctx, tenantID)
	if err != nil {
		return model.Invoice{}, ErrNotFound
	}
	price, err := s.pricing(ctx, t.Plan)
	if err != nil {
		return model.Invoice{}, err
	}
	seats, err := s.repo.Seats(ctx, tenantID)
	if err != nil {
		return model.Invoice{}, err
	}
	return model.Invoice{Plan: price.Code, Seats: seats, TotalCent: price.UnitCents + seats*price.SeatCents}, nil
}

func (s *Billing) pricing(ctx context.Context, planCode string) (model.PlanPrice, error) {
	return s.repo.PriceFor(ctx, planCode)
}

func (s *Billing) Quote(ctx context.Context, plan string) (model.PlanPrice, error) {
	return s.repo.PriceFor(ctx, plan)
}

func (s *Billing) ChangePlan(ctx context.Context, tenantID int64, plan config.Plan) (config.Plan, error) {
	seats, err := s.repo.Seats(ctx, tenantID)
	if err != nil {
		return config.Plan{}, err
	}
	if seats > plan.SeatLimit {
		return config.Plan{}, ErrSeatLimit
	}
	return plan, s.repo.SetPlan(ctx, tenantID, plan.Code)
}

func (s *Billing) Preview(ctx context.Context, tenantID int64, kind string) (string, string, error) {
	t, err := s.tenants.Get(ctx, tenantID)
	if err != nil {
		return "", "", ErrNotFound
	}
	seats, err := s.repo.Seats(ctx, tenantID)
	if err != nil {
		return "", "", err
	}
	data := map[string]any{"Period": time.Now().UTC().Format("2006-01"), "Seats": seats}

	footerTpl, err := template.New("footer").Parse(t.InvoiceFooter)
	if err != nil {
		return "", "", ErrInvalid
	}
	bodySrc, ok := config.InvoiceTemplates[kind]
	if !ok {
		bodySrc = config.InvoiceTemplates["standard"]
	}
	bodyTpl, err := template.New("body").Parse(bodySrc)
	if err != nil {
		return "", "", err
	}
	var head, foot bytes.Buffer
	if err := bodyTpl.Execute(&head, data); err != nil {
		return "", "", err
	}
	if err := footerTpl.Execute(&foot, data); err != nil {
		return "", "", ErrInvalid
	}
	return head.String(), foot.String(), nil
}
