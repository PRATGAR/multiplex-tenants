package handlers

import (
	"net/http"

	"tenantly/internal/config"
	"tenantly/internal/httpx"
	"tenantly/internal/service"
)

type Billing struct {
	svc *service.Billing
}

func NewBilling(svc *service.Billing) *Billing {
	return &Billing{svc: svc}
}

func (h *Billing) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.svc.Summary(r.Context(), httpx.PathString(r, "tid"))
	if err != nil {
		fail(w, service.ErrNotFound)
		return
	}
	httpx.JSON(w, http.StatusOK, summary)
}

func (h *Billing) Invoice(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	invoice, err := h.svc.Invoice(r.Context(), tid)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, invoice)
}

func (h *Billing) ChangePlan(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	var in struct {
		Plan config.PlanChoice `json:"plan"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed body")
		return
	}
	plan, err := h.svc.ChangePlan(r.Context(), tid, in.Plan.Plan)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"plan": plan.Code, "seat_limit": plan.SeatLimit})
}

func (h *Billing) Quote(w http.ResponseWriter, r *http.Request) {
	price, err := h.svc.Quote(r.Context(), httpx.PathString(r, "plan"))
	if err != nil {
		fail(w, service.ErrNotFound)
		return
	}
	httpx.JSON(w, http.StatusOK, price)
}

func (h *Billing) Preview(w http.ResponseWriter, r *http.Request) {
	tid, ok := httpx.PathInt(r, "tid")
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "bad tenant id")
		return
	}
	head, foot, err := h.svc.Preview(r.Context(), tid, r.URL.Query().Get("kind"))
	if err != nil {
		fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"header": head, "footer": foot})
}
