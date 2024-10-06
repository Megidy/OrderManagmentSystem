package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/Megidy/OrderManagmentSystem/pkg/db"
	"github.com/Megidy/OrderManagmentSystem/pkg/kitchen/producer"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
	"github.com/goccy/go-json"
)

const (
	Topic string = "orders"
)

var wg sync.WaitGroup

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
				if key == "Create_order" {
					var order types.Order

					err := json.Unmarshal(msg.Value, &order)
					if err != nil {

						log.Println(err)
					}
					// order := string(msg.Value)
					log.Println("received order from service orders ", order)
					wg.Add(1)
					fmt.Println(order)
					go ProcessOrder(order, &wg)

				}

			case <-sigChan:
				log.Println("interrupt")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Println("exited goroutine ")
	err = worker.Close()
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()

}

func ProcessOrder(order types.Order, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := db.ChangeOrderStatus(order, "in_progress")
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		time.Sleep(1000 * time.Millisecond)
		log.Println("preparing order", order.OrderId, " :", i*10, "%")

	}
	log.Println("order :", order.OrderId, " is ready!!")
	err = db.ChangeOrderStatus(order, "completed")
	if err != nil {
		return err
	}
	producer.SendKitchenMessage(order.OrderId, order.CustomerId)

	return nil
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer(brokers, config)
}
