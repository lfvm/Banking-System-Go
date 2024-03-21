package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshSessionRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshSessionRes struct {
	AccesToken string `json:"acces_token"`
}

func (server *Server) refreshToken(ctx *gin.Context) {

	var req refreshSessionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fmt.Printf("Token id:  %s", payload.ID)

	session, err := server.store.GetSession(ctx, payload.ID)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// create acces token
	accesToken, _, err := server.tokenMaker.CreateToken(session.Username, time.Minute*15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := refreshSessionRes{
		AccesToken: accesToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}
