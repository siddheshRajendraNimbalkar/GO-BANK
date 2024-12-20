package api

import (
	"database/sql"
	"net/http"
	"strconv"

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

func (server *Server) GetAcount(ctx *gin.Context) {
	id := ctx.Param("id")

	accountID, err := strconv.Atoi(id)
	if err != nil || accountID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := server.store.GetAccount(ctx, int64(accountID))

	if err != nil {
		// Handle the case where the account is not found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		// Handle other potential errors
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		return
	}

	// Respond with the account details
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(ctx *gin.Context) {
	// Optional: Parse pagination query parameters
	var params struct {
		Page     int `form:"page" binding:"min=1"`
		PageSize int `form:"page_size" binding:"min=5,max=100"`
	}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination values if not provided
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}

	// Calculate offset for pagination
	offset := (params.Page - 1) * params.PageSize

	// Fetch accounts from the database
	accounts, err := server.store.ListAccount(ctx, db.ListAccountParams{
		Limit:  int32(params.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	// Respond with the list of accounts
	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) DeleteAccounts(ctx *gin.Context) {
	id := ctx.Param("id")
	accountID, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = server.store.DeleteAccount(ctx, int64(accountID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})

}

type UpdateAccountParams struct {
	Balance int64 `db:"balance"  binding:"required"`
}

func (server *Server) UpdateAccount(ctx *gin.Context) {

	id := ctx.Param("id")
	accountID, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req UpdateAccountParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	arg := db.UpdateAccountParams{
		ID:      int64(accountID),
		Balance: req.Balance,
	}

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	ctx.JSON(http.StatusFound, updatedAccount)
	return
}
