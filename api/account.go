package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/token"
	"github.com/lib/pq"
)

/*
Account creation body
*/
type createAccountRequest struct {
	Currency string `json:"Currency" binding:"required,currency"`
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PASETOPayload)

	params := database.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0.0,
		Currency: req.Currency,
	}

	acc, err := s.store.CreateAccount(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {

			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PASETOPayload)
	if acc.Owner != authPayload.Username {
		err := fmt.Errorf("User %s declared in token is unauthorized to access account %s", authPayload.Username, accID.String())
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.PASETOPayload)
	accs, err := s.store.GetAccountsList(ctx, database.GetAccountsListParams{
		Owner:  authPayload.Username,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accs)
}
