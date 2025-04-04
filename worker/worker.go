package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	topic := "comments"
	worker, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		log.Fatal("Error connecting to consumer:", err)
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal("Error consuming partition:", err)
		panic(err)
	}

	fmt.Println("Consumer started, waiting for messages...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0

	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println("Error:", err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Message received: %s | Topic: %s | Count: %d\n", string(msg.Value), msg.Topic, msgCount)
			case <-sigchan:
				fmt.Println("Interrupt is detected, exiting...")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
	if err := consumer.Close(); err != nil {
		log.Println("Error closing consumer:", err)
		panic(err)
	} else {
		fmt.Println("Consumer closed")
	}
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		log.Println("Error creating consumer:", err)
		return nil, err
	}

	return conn, nil
}
