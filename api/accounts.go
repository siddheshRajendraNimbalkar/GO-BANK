package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type CreateAccountParams struct {
	Owner    string `db:"owner" binding:"required"`
	Currency string `db:"currency" binding:"required,oneof=INR"`
}

func (server *Server) createAcount(ctx *gin.Context) {
	var req CreateAccountParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"account": account,
	})
	return
}
