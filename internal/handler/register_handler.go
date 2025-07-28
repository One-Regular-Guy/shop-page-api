package handler

import (
	"github.com/One-Regular-Guy/shop-page-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterHandler struct {
	UseCase *usecase.RegisterUserUseCase
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	output, err := h.UseCase.Execute(usecase.RegisterUserInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
		return
	}
	if output.Message == "Username already exists" || output.Message == "Email already exists" {
		c.JSON(http.StatusBadRequest, gin.H{"message": output.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": output.Message})
}
