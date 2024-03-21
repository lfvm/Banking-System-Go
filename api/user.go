package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/lfvm/simplebank/db/sqlc"
	"github.com/lfvm/simplebank/utils"
	"github.com/lib/pq"
)

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		FullName:          user.FullName,
	}
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hash,
		Username:       req.UserName,
	}

	user, err := server.store.CreateUser(context.Background(), arg)

	if err != nil {

		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusCreated, rsp)
}

type logInRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type logInRes struct {
	AccesToken   string       `json:"acces_token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
	SessonId     uuid.UUID    `json:"session_id"`
}

func (server *Server) logIn(ctx *gin.Context) {

	var req logInRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserName)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// validate passwords
	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// create acces token
	accesToken, payload, err := server.tokenMaker.CreateToken(user.Username, time.Minute*15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(payload.Username, time.Hour*24)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		Username:     user.Username,
		ID:           refreshTokenPayload.ID,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := logInRes{
		AccesToken:   accesToken,
		User:         newUserResponse(user),
		RefreshToken: refreshToken,
		SessonId:     session.ID,
	}
	ctx.JSON(http.StatusOK, rsp)
}
