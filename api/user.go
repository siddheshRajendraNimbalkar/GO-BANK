package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type CreateUserParams struct {
	Username       string `db:"username" binding:"required`
	HashedPassword string `db:"hashed_password" binding:"required`
	FullName       string `db:"full_name" binding:"required`
	Email          string `db:"email" binding:"required`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong Arrguments",
			"error":   err.Error(),
		})
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.HashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "Fail to create user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})

	return
}

func (server *Server) getUser(ctx *gin.Context) {
	id := ctx.Param("userName")

	user, err := server.store.GetUser(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong Arrguments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusFound, gin.H{
		"message": "user found",
		"user":    user,
	})

	return
}
