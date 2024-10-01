package main

import (
	"log"
	"net/http"

	"github.com/Megidy/OrderManagmentSystem/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/v1", func(ctx *gin.Context) {
		config.ConnectDB()
		db := config.GetDB()

		// Перевіряємо чи db не nil
		if db == nil {
			log.Println("Database connection is nil")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect to the database",
			})
			return
		}

		// Виконуємо запит до бази даних
		_, err := db.Exec("select * from test")
		if err != nil {
			log.Println("Database query failed:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to execute query",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	router.Run(":8080")
}
