package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

type createUserResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
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

	rsp := createUserResponse{
		Username:          user.Username,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		FullName:          user.FullName,
	}
	ctx.JSON(http.StatusCreated, rsp)
}
