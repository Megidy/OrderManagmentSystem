package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

const (
	Topic string = "orders"
)

func main() {
	worker, err := ConnectConsumer([]string{"kafka:9092"})
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := worker.ConsumePartition(Topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneCh := make(chan struct{})
	go func() {
		select {
		case err := <-consumer.Errors():
			log.Println(err)
		case msg := <-consumer.Messages():

			order := string(msg.Value)
			log.Println("received order from service orders ", order)
		case <-sigChan:
			log.Println("interrupt")
			doneCh <- struct{}{}
		}

	}()

	<-doneCh
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
