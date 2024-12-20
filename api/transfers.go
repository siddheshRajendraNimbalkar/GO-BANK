package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type CreateTransfersParams struct {
	FromAccID int64 `db:"from_acc_id" binding:"required"`
	ToAccID   int64 `db:"to_acc_id"  binding:"required"`
	Amount    int64 `db:"amount"  binding:"required"`
}

func (server *Server) Transfer(ctx *gin.Context) {
	var req CreateTransfersParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Amount cant be negative or zero",
		})
	}

	arg := db.CreateTransfersParams{
		FromAccID: req.FromAccID,
		ToAccID:   req.ToAccID,
		Amount:    req.Amount,
	}

	transfer, err := server.store.CreateTransfers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, transfer)

}
