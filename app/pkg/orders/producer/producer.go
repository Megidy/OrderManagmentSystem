package producer

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/db"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/Megidy/OrderManagmentSystem/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	Topic string = "orders"
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

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal body ", http.StatusInternalServerError, nil)
		return
	}

	err = PushOrderToQueue(Topic, orderInBytes, "Create_order")
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

func HandleCheckOrders(c *gin.Context) {

	checkOrders, err := db.CheckOrdersStatus()
	if err != nil {
		utils.HandleError(c, err, "failed to get CheckOrderStatus", http.StatusInternalServerError, nil)
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

func PushOrderToQueue(topic string, message []byte, key string) error {
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
