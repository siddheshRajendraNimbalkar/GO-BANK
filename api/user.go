package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserParams struct {
	Username       string `db:"username" binding:"required"`
	HashedPassword string `db:"hashed_password" binding:"required,min=8"`
	FullName       string `db:"full_name" binding:"required"`
	Email          string `db:"email" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserParams

	// Validate request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid arguments",
			"error":   err.Error(),
		})
		return
	}

	// Hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while hashing password",
			"error":   err.Error(),
		})
		return
	}

	// Prepare user parameters
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: string(hashPassword),
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// Create user in the database
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	// Generate JWT token
	maker, err := token.NewPasetoMaker(server.config.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating token maker",
			"error":   err.Error(),
		})
		return
	}

	tokenStr, err := maker.CreateToken(user.Username, server.config.JwtDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating token",
			"error":   err.Error(),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   tokenStr,
	})
}

func (server *Server) getUser(ctx *gin.Context) {
	username := ctx.Param("userName")

	// Fetch user from database
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}

	// Return user details
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User found",
		"user": gin.H{
			"Username":  user.Username,
			"FullName":  user.FullName,
			"Email":     user.Email,
			"CreatedAt": user.CreatedAt,
		},
	})
}

type SignInUserParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) compareUser(ctx *gin.Context) {
	var req SignInUserParams

	// Validate input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid arguments",
			"error":   err.Error(),
		})
		return
	}

	// Get user from database
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong password",
			"error":   err.Error(),
		})
		return
	}

	maker, err := token.NewPasetoMaker(server.config.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating token maker",
			"error":   err.Error(),
		})
		return
	}

	tokenStr, err := maker.CreateToken(user.Username, server.config.JwtDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating token",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User found successfully",
		"token":   tokenStr,
	})
}
