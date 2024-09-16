package utils

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, err error, message string, statusCode int) {
	if err != nil {
		c.JSON(statusCode, gin.H{
			"error":   err,
			"details": message,
		})
		c.Abort()
	}

}
