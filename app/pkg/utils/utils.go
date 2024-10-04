package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, msg string, statusCode int, v interface{}) {
	c.JSON(statusCode, gin.H{
		"error":   err,
		"details": msg,
	})
	log.Println("error : ", err, " details : ", msg)

}
