package middleware

import (
	"net/http"
	"strconv"

	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/Megidy/OrderManagmentSystem/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RequireOrder(c *gin.Context) {
	Id, err := c.Cookie("Order")
	if err != nil {
		utils.HandleError(c, err, "failed to get cookie", http.StatusBadRequest, nil)
		return
	}
	orderId, err := strconv.Atoi(Id)
	if err != nil {
		utils.HandleError(c, err, "failed to converte order ", http.StatusInternalServerError, nil)
		return
	}
	customer := types.Customer{
		OrderId: orderId,
	}

	if orderId != 0 {
		c.Set("customer", &customer)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusNoContent)
	}

}
