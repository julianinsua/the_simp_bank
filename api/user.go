package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=5"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func newUserResponse(user database.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := database.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hash,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	usr, err := s.store.CreateUser(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	rsp := newUserResponse(usr)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId             uuid.UUID    `json:"sessionId"`
	User                  userResponse `json:"user"`
	Token                 string       `json:"token"`
	TokenExpiresAt        time.Time    `json:"tokenExpiresAt"`
	RefreshToken          string       `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time    `json:"refreshTokenExpiresAt"`
}

func (srv *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//get user via username
	usr, err := srv.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// validate password
	err = util.CheckPassword(req.Password, usr.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create token
	token, tokenPayload, err := srv.tokenMaker.CreateToken(usr.Username, srv.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create refresh token
	refreshToken, refreshTokenPayload, err := srv.tokenMaker.CreateToken(usr.Username, srv.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create session
	ssn := database.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     usr.Username,
		RefreshToken: refreshToken,
		ClientAgent:  ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshTokenPayload.ExpiresAt,
	}
	session, err := srv.store.CreateSession(ctx, ssn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// respond
	res := loginUserResponse{
		SessionId:             session.ID,
		User:                  newUserResponse(usr),
		Token:                 token,
		TokenExpiresAt:        tokenPayload.ExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, res)
}
