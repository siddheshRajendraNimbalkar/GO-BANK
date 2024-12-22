package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type CreateTransfersParams struct {
	FromAccID int64 `db:"from_acc_id" binding:"required"`
	ToAccID   int64 `db:"to_acc_id"  binding:"required"`
	Amount    int64 `db:"amount"  binding:"required,gt=0"`
}

func (server *Server) Transfer(ctx *gin.Context) {
	var req CreateTransfersParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fromAccount, err := server.store.GetAccount(ctx, req.FromAccID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "From Account Not Found",
			"From_ID": req.FromAccID,
		})
		return
	}

	toAccount, err := server.store.GetAccount(ctx, req.ToAccID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "To Account Not Found",
			"To_ID":   req.ToAccID,
		})
		return
	}

	if fromAccount.Currency != toAccount.Currency {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Currency mismatch",
		})
		return
	}

	if req.FromAccID == req.ToAccID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "From and To account IDs must be different",
		})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccID,
		ToAccountID:   req.ToAccID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
