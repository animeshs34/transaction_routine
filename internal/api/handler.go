package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/animesh/transaction_routine/internal/service"
	customErrors "github.com/animesh/transaction_routine/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	accountService     service.AccountService
	transactionService service.TransactionService
}

func NewHandler(accSvc service.AccountService, txSvc service.TransactionService) *Handler {
	return &Handler{
		accountService:     accSvc,
		transactionService: txSvc,
	}
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccount(c *gin.Context) {
	idStr := c.Param("accountId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, customErrors.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" binding:"required"`
	OperationTypeID int     `json:"operation_type_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tx, err := h.transactionService.CreateTransaction(
		c.Request.Context(),
		req.AccountID,
		req.OperationTypeID,
		req.Amount,
	)

	if err != nil {
		if errors.Is(err, customErrors.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, customErrors.ErrOperationType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, tx)
}
