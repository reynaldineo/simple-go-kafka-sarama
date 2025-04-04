# Simple Go Kafka Sarama

This project demonstrates a simple implementation of a Kafka producer and consumer using the [Sarama](https://github.com/Shopify/sarama) library in Go. It also includes a REST API built with [Fiber](https://gofiber.io/) to send messages to Kafka.

## Features

- **Producer**: Sends messages to a Kafka topic.
- **Consumer**: Listens to a Kafka topic and processes messages.
- **REST API**: Provides an endpoint to send messages to Kafka.

## Prerequisites

- Docker and Docker Compose
- Go 1.24.1 or higher
- Kafka and Zookeeper (configured via `docker-compose.yml`)

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/reynaldineo/simple-go-kafka-sarama.git
   cd simple-go-kafka-sarama
   ```

2. Start Kafka and Zookeeper using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Running the Project

### Producer (API Server)

The producer is implemented in the `producer/producer.go` file. It provides a REST API endpoint (`/api/comments`) to accept user input and send it as a message to a Kafka topic.

1. Navigate to the `producer` directory:
   ```bash
   cd producer
   ```

2. Run the producer:
   ```bash
   go run producer.go
   ```

3. The API server will start on `http://localhost:3000`.

4. Use the `/api/comments` endpoint to send messages to Kafka:
   ```bash
   curl -X POST http://localhost:3000/api/comments -H "Content-Type: application/json" -d '{"text":"Your comment here"}'
   ```

   - The producer connects to Kafka using the Sarama library.
   - It sends the message to the `comments` topic.
   - The message is logged with details about the partition and offset.

### Consumer (Worker)

The consumer is implemented in the `worker/worker.go` file. It listens to the `comments` topic and processes incoming messages.

1. Navigate to the `worker` directory:
   ```bash
   cd worker
   ```

2. Run the consumer:
   ```bash
   go run worker.go
   ```

3. The consumer will:
   - Connect to Kafka using the Sarama library.
   - Subscribe to the `comments` topic.
   - Process messages in real-time, logging the message content, topic, and message count.

4. To stop the consumer, press `Ctrl+C`. It will gracefully shut down and log the total number of processed messages.

## Project Structure

- **`producer/producer.go`**: Contains the REST API and Kafka producer logic.
- **`worker/worker.go`**: Contains the Kafka consumer logic.
- **`docker-compose.yml`**: Configures Kafka and Zookeeper services.
- **`go.mod`**: Manages Go dependencies.

## Dependencies

- [Sarama](https://github.com/Shopify/sarama): Kafka client library for Go.
- [Fiber](https://gofiber.io/): Web framework for building REST APIs.

## Notes

- Ensure Kafka is running and accessible at `localhost:9092`.
- The default Kafka topic used is `comments`.

