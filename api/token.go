package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenParams

	// Validate input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid arguments",
			"error":   err.Error(),
		})
		return
	}

	refreshToken, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid refresh token",
			"error":   err.Error(),
		})
		return
	}

	// Get session from database
	session, err := server.store.GetSession(ctx, refreshToken.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Session not found",
			"error":   err.Error(),
		})
		return
	}

	// Check if session is blocked
	if session.IsBlocked {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Session is blocked",
		})
		return
	}

	// Check if session is expired
	if session.ExpireDate.Before(time.Now()) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Session is expired",
		})
		return
	}

	// Generate new access token
	accessToken, userPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.SessionDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
		return
	}

	// Return new access token
	ctx.JSON(http.StatusOK, renewTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   userPayload.ExpiresAt,
	})
}
