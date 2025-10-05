package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/animeshs34/transaction_routine/internal/respository"
	"github.com/animeshs34/transaction_routine/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	// Accounts
	mux.HandleFunc("/accounts", h.accountsRoot) // POST
	mux.HandleFunc("/accounts/", h.accountsOne) // GET /accounts/{id}

	// Transactions
	mux.HandleFunc("/transactions", h.transactionsRoot) // POST

	// Health - this is probing endpoints
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}

func (h *Handler) accountsRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createAccount(w, r)
	default:
		methodNotAllowed(w, http.MethodPost)
	}
}

func (h *Handler) accountsOne(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /accounts/{id}
		idStr := strings.TrimPrefix(r.URL.Path, "/accounts/")
		if idStr == "" || strings.Contains(idStr, "/") {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid account id")
			return
		}
		acc, err := h.svc.GetAccount(id)
		if err != nil {
			if errors.Is(err, respository.ErrAccountNotFound) {
				writeError(w, http.StatusNotFound, "account not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "could not get account")
			return
		}
		writeJSON(w, http.StatusOK, acc)
	default:
		methodNotAllowed(w, http.MethodGet)
	}
}

type createAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	acc, err := h.svc.CreateAccount(req.DocumentNumber)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDocument) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create account")
		return
	}
	writeJSON(w, http.StatusCreated, acc)
}

type createTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
	EventDate       *string `json:"event_date,omitempty"` // optional; RFC3339
}

func (h *Handler) transactionsRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req createTransactionRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		var t *time.Time
		if req.EventDate != nil && *req.EventDate != "" {
			parsed, err := time.Parse(time.RFC3339Nano, *req.EventDate)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid event_date; must be RFC3339")
				return
			}
			t = &parsed
		}

		tx, err := h.svc.CreateTransaction(req.AccountID, req.OperationTypeID, req.Amount, t)
		if err != nil {
			switch {
			case errors.Is(err, respository.ErrAccountNotFound):
				writeError(w, http.StatusNotFound, "account not found")
			case errors.Is(err, service.ErrInvalidOperationType):
				writeError(w, http.StatusBadRequest, "invalid operation_type_id")
			case errors.Is(err, service.ErrInvalidAmount):
				writeError(w, http.StatusBadRequest, "amount must be greater than zero")
			default:
				writeError(w, http.StatusInternalServerError, "could not create transaction")
			}
			return
		}
		writeJSON(w, http.StatusCreated, tx)
	default:
		methodNotAllowed(w, http.MethodPost)
	}
}

// Helpers

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func methodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
