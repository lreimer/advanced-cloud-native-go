package main

import (
	"fmt"
	"os"

	"time"

	"github.com/Shopify/sarama"
)

func main() {
	fmt.Println("Starting synchronous Kafka producer...")
	time.Sleep(5 * time.Second)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokers := []string{brokerAddr()}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	topic := topic()
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})

	go func() {
		for {
			msgCount++

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(fmt.Sprintf("Hello Kafka %v", msgCount)),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
			time.Sleep(5 * time.Second)
		}
	}()

	<-doneCh
}

func brokerAddr() string {
	brokerAddr := os.Getenv("BROKER_ADDR")
	if len(brokerAddr) == 0 {
		brokerAddr = "localhost:9092"
	}
	return brokerAddr
}

func topic() string {
	topic := os.Getenv("TOPIC")
	if len(topic) == 0 {
		topic = "default-topic"
	}
	return topic
}
