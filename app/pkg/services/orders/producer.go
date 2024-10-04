package orders

import (
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/gin-gonic/gin"
)

func HandleCreateOrder(c *gin.Context) {
	var NewOrder types.Order

	err := c.ShouldBindBodyWithJSON(&NewOrder)
}
