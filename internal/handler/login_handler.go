package handler

import (
	"github.com/One-Regular-Guy/shop-page-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginHandler struct {
	UseCase *usecase.LoginUserUseCase
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	output, err := h.UseCase.Execute(usecase.LoginUserInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to login user"})
		return
	}
	if output.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": output.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": output.Token})
}
