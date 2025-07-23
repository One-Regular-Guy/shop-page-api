package register

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Handler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	if ServiceInstance.CheckUsername(request.Username) {
		c.JSON(400, Response{Message: "Username already exist"})
		return
	}
	if ServiceInstance.CheckEmail(request.Email) {
		c.JSON(400, Response{Message: "Email already exist"})
		return
	}

	err := ServiceInstance.RegisterUser(request.Name, request.Username, request.Email, request.Password)
	if err != nil {
		log.Println("Error registering user:", err)
		c.JSON(500, Response{Message: "Failed to register user"})
		return
	}

	c.JSON(200, Response{Message: "User registered successfully"})
}
