package http

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type ChargeHandler struct {
	failFirst int64
	requests  atomic.Int64
}

func NewChargeHandler(failFirst int64) *ChargeHandler {
	return &ChargeHandler{failFirst: failFirst}
}

type chargeRequest struct {
	PaymentID  string `json:"payment_id"`
	Amount     int64  `json:"amount_minor"`
	Currency   string `json:"currency"`
}

func (h *ChargeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request chargeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if request.PaymentID == "" || request.Amount <= 0 || request.Currency == "" {
		http.Error(w, "payment_id, positive amount_minor, and currency are required", http.StatusBadRequest)
		return
	}

	number := h.requests.Add(1)
	if number <= h.failFirst {
		w.Header().Set("Retry-After", "1")
		http.Error(w, "simulated temporary gateway failure", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"gateway_reference": "gw-" + request.PaymentID,
		"status":            "approved",
	})
}