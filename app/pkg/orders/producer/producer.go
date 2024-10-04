package producer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/Megidy/OrderManagmentSystem/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	Topic string = "orders"
)

func HandleCreateOrder(c *gin.Context) {

	var order types.Order

	err := c.ShouldBindJSON(&order)
	if err != nil {
		utils.HandleError(c, err, "failed to read body", http.StatusBadRequest, nil)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal body ", http.StatusInternalServerError, nil)
		return
	}
	err = PushOrderToQueue(Topic, orderInBytes)
	if err != nil {
		utils.HandleError(c, err, "failed to push order ", http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "order is handling ",
	})
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)

}

func PushOrderToQueue(topic string, message []byte) error {
	brokers := []string{"kafka:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := sarama.ProducerMessage{
		Topic: topic,
		//add key later
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	log.Printf("status is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topic, partition, offset)
	return nil
}
