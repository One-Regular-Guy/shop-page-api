package login

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := ServiceInstance.CheckUsername(req.Username); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(200, gin.H{"username": req.Username})

	/*token, err := Login(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, Response{Token: token})*/
}
