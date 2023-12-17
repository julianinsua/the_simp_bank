package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/julianinsua/the_simp_bank/internal/database"
)

/*
Account creation body
*/
type transferRequest struct {
	FromAccountID uuid.UUID `json:"FromAccountId" binding:"required"`
	ToAccountID   uuid.UUID `json:"ToAccountId" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Currency      string    `json:"currency" binding:"required,currency"`
}

/*
Account creation handler
*/
func (s Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !s.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	params := database.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s Server) validAccount(ctx *gin.Context, accId uuid.UUID, currency string) bool {
	acc, err := s.store.GetAccount(ctx, accId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if acc.Currency != currency {
		err = fmt.Errorf("Account [%s] currency mismatch: has %s vs. %s", accId, acc.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
