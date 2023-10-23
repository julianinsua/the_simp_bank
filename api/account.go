package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/julianinsua/the_simp_bank.git/internal/database"
)

/*
Account creation body
*/
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"Currency" binding:"required,oneof=USD EUR"`
}

/*
Account creation handler
*/
func (s Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := database.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0.0,
		Currency: req.Currency,
	}

	acc, err := s.store.CreateAccount(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

/*
Get account by id url params
*/
type GetAccountRequest struct {
	ID string `uri:"id" binding:"required"`
}

/*
Get account by id handler
*/
func (s Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	acc, err := s.store.GetAccount(ctx, accID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

/*
Get Account List Query parameters
*/

type GetAccountListRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5,max=10"`
}

/*
Get account list handler
*/

func (s Server) getAccountList(ctx *gin.Context) {
	var req GetAccountListRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accs, err := s.store.GetAccountsList(ctx, database.GetAccountsListParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accs)
}
