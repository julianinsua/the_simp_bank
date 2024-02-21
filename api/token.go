package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type refreshTokenResponse struct {
	Token          string    `json:"token"`
	TokenExpiresAt time.Time `json:"tokenExpiresAt"`
}

func (srv *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := srv.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//get user via token payload's username
	ssn, err := srv.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check session blocked
	if ssn.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Session blocked")))
		return
	}

	// Check session username to be the same as the token's username

	if ssn.Username != payload.Username {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("incorrect session user"))
		return
	}

	// Check that the request's refresh token is the same as the session's refresh token
	if req.RefreshToken != ssn.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("unmatching token"))
		return
	}

	if time.Now().After(ssn.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("expired token"))
		return
	}

	// Create token
	token, tokenPayload, err := srv.tokenMaker.CreateToken(payload.Username, srv.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// respond
	res := refreshTokenResponse{
		Token:          token,
		TokenExpiresAt: tokenPayload.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, res)
}
