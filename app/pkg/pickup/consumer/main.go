package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/db"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/goccy/go-json"
)

const (
	Topic string = "Kitchen"
)

func main() {
	worker, err := ConnectConsumer([]string{"kafka:9092"})
	if err != nil {
		log.Fatal(err)
	}
	consumer, err := worker.ConsumePartition(Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				key := string(msg.Key)

				if key == "Send_order" {
					var order types.Order
					err := json.Unmarshal(msg.Value, &order)
					if err != nil {
						log.Println(err)
					}

					log.Println("received ready order from service kitchen  ", order)
				}
				if key == "Take_order" {
					var customer types.Customer

					err = json.Unmarshal(msg.Value, &customer)
					if err != nil {
						log.Println(err)
					}
					err = db.DeleteOrder(customer.OrderId)
					if err != nil {
						log.Println(err)
					}

				}
			case <-sigChan:
				log.Println("interrupt")
				doneCh <- struct{}{}
			}

		}
	}()
	<-doneCh
	log.Println("exited goroutine")
	err = worker.Close()
	if err != nil {
		log.Fatal(err)
	}

}
func ConnectConsumer(brokers []string) (sarama.Consumer, error) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer(brokers, config)

}
