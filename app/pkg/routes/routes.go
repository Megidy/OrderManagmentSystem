package routes

import (
	"github.com/Megidy/OrderManagmentSystem/pkg/middleware"
	"github.com/Megidy/OrderManagmentSystem/pkg/orders/producer"
	"github.com/gin-gonic/gin"
)

var InitRoutes = func(router *gin.Engine) {
	router.POST("/order/create", producer.HandleCreateOrder)
	router.GET("/orders", producer.HandleCheckOrders)
	router.GET("/myorders", middleware.RequireOrder, producer.HandlerCheckCusatomersOrder)
	router.DELETE("/orders/take", middleware.RequireOrder, producer.HandleTakeOrder)
}
