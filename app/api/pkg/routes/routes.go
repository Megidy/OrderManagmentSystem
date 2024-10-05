package routes

import (
	"github.com/Megidy/OrderManagmentSystem/pkg/orders/producer"
	"github.com/gin-gonic/gin"
)

var InitRoutes = func(router *gin.Engine) {
	router.POST("/order", producer.HandleCreateOrder)

}
