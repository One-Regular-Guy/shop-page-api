package login

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return
	}
	if !ServiceInstance.CheckUsername(req.Username) {
		c.JSON(401, ErrorResponse{Message: "Credentials aren't valid"})
		return
	}
	if !ServiceInstance.CheckPassword(req.Username, req.Password) {
		c.JSON(401, ErrorResponse{Message: "Credentials aren't valid"})
		return
	}
	signedToken, err := ServiceInstance.ProvideToken(req.Username)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to generate token"})
		return
	}
	c.JSON(200, Response{Token: signedToken})
}
