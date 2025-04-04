package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/comments", createComment)
	app.Listen(":3000")
}

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		log.Println("Error creating producer:", err)
		return nil, err
	}
	return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
	brokersUrl := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		log.Println("Error connecting to producer:", err)
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	fmt.Printf("Message sent to topic %s partition %d at offset %d\n", topic, partition, offset)
	return nil
}

func createComment(c *fiber.Ctx) error {
	cmt := new(Comment)
	if err := c.BodyParser(cmt); err != nil {
		log.Println("Error parsing request body:", err)
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return err
	}

	cmtInBytes, err := json.Marshal(cmt)
	if err != nil {
		log.Println("Error marshalling comment:", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
		return err
	}

	PushCommentToQueue("comments", cmtInBytes)

	err = c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment created successfully",
		"comment": cmt,
	})

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
		return err
	}

	return err
}
