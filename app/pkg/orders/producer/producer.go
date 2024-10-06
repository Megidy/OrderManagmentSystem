package producer

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/db"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/Megidy/OrderManagmentSystem/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	Topic1 string = "orders"
	Topic2 string = "Kitchen"
)

func HandleCreateOrder(c *gin.Context) {

	var NewOrder types.CreateOrder

	err := c.ShouldBindJSON(&NewOrder)
	if err != nil {
		utils.HandleError(c, err, "failed to read body", http.StatusBadRequest, nil)
		return
	}

	//TODO
	//create more uniqu ids
	orderId := rand.Intn(100)
	customerId := rand.Intn(100)
	order := types.Order{
		OrderId:    orderId,
		CustomerId: customerId,
		Dishes:     NewOrder.Dishes,
		Status:     "pending",
	}

	customer := types.Customer{
		OrderId: orderId,
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Order", strconv.Itoa(customer.OrderId), 3600*24*10, "", "", false, true)

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal body ", http.StatusInternalServerError, nil)
		return
	}

	err = PushMessageToQueue(Topic1, orderInBytes, "Create_order")
	if err != nil {
		utils.HandleError(c, err, "failed to push order ", http.StatusInternalServerError, nil)
		return
	}

	err = db.CreateOrder(order)
	if err != nil {
		utils.HandleError(c, err, "failed to create order", http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "order is handling ",
	})

}

func HandleTakeOrder(c *gin.Context) {
	customer, ok := c.Get("customer")
	if !ok {
		utils.HandleError(c, nil, "you dont have orders", http.StatusBadRequest, nil)
		return
	}

	customerInBytes, err := json.Marshal(customer)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal data", http.StatusInternalServerError, nil)
		return
	}
	order, err := db.GetOrder(customer.(*types.Customer).OrderId, 0)
	if err != nil {
		utils.HandleError(c, err, "failed GetOrder", http.StatusInternalServerError, nil)
		return
	}
	order.Dishes, err = db.GetDishes(customer.(*types.Customer).OrderId)
	if err != nil {
		utils.HandleError(c, err, "failed GetDishes", http.StatusInternalServerError, nil)
		return
	}
	if order.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"you order is not ready yet": ":(",
		})
		return
	}
	err = PushMessageToQueue(Topic2, customerInBytes, "Take_order")
	if err != nil {
		utils.HandleError(c, err, "failed to marshal data", http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": " you took your order",
		"order":   order,
	})

}

func HandlerCheckCusatomersOrder(c *gin.Context) {
	customer, ok := c.Get("customer")
	if !ok {
		utils.HandleError(c, nil, "you dont have orders", http.StatusBadRequest, nil)
		return
	}

	orders, err := db.CheckcustomersOrdersStatus(customer.(*types.Customer).OrderId)
	if err != nil {
		utils.HandleError(c, err, "failed to  CheckcustomersOrdersStatus", http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"your orders": orders,
	})

}

func HandleCheckOrders(c *gin.Context) {

	checkOrders, err := db.CheckOrdersStatus()
	if err != nil {
		utils.HandleError(c, err, "failed to CheckOrderStatus", http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders ": checkOrders,
	})
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)

}

func PushMessageToQueue(topic string, message []byte, key string) error {
	brokers := []string{"kafka:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	log.Printf("status is stored in topic(%s)/partition(%d)/offset(%d),key(%s)\n",
		topic, partition, offset, key)
	return nil
}
