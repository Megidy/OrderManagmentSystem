package orders

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

	err := c.ShouldBindBodyWithJSON(&order)
	if err != nil {
		utils.HandleError(c, err, "failed to read body ", http.StatusBadRequest, nil)
		return
	}
	orderInBytes, err := json.Marshal(order)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal order ", http.StatusBadRequest, nil)
		return
	}
	err = PushOrderToQueue(Topic, orderInBytes)
	if err != nil {
		utils.HandleError(c, err, "failed to push order ", http.StatusInternalServerError, nil)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "your order is already processing ",
	})

}

func PushOrderToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}

	defer producer.Close()
	msg := sarama.ProducerMessage{
		Topic: topic,
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

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return sarama.NewSyncProducer(brokers, config)
}
