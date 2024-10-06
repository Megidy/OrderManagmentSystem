package producer

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/db"
)

const (
	Topic string = "Kitchen"
)

func SendKitchenMessage(orderId, customerId int) {
	order, err := db.GetOrder(orderId, customerId)
	if err != nil {
		log.Println(err)
	}
	dishes, err := db.GetDishes(orderId)
	if err != nil {
		log.Println(err)
	}
	order.Dishes = dishes

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
	}
	PushKitchenMessageToQueue(Topic, orderInBytes, "Send_order")
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return sarama.NewSyncProducer(brokers, config)

}
func PushKitchenMessageToQueue(Topic string, message []byte, key string) error {
	brokers := []string{"kafka:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	log.Printf("Kitchen Message is stored in topic(%s)/partition(%d)/offset(%d),key(%s)\n",
		Topic, partition, offset, key)
	return nil
}
