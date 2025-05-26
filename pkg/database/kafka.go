package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CloseFunction func() error

func Publish(topic string, msg []byte) error {
	url := os.Getenv("KAFKA_URL")
	if url == "" {
		log.Println("KAFKA_URL is not defined")
	}

	config := kafka.ConfigMap{
		"bootstrap.servers": url,
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		err = fmt.Errorf("error while trying to create kafka-producer: %v", err)
		log.Printf("v", err)
		return err
	}
	defer producer.Close()

	ev := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, ev)
	defer close(ev)

	if err != nil {
		log.Printf("v", err)
		return err
	}

	for e := range ev {
		switch eMsg := e.(type) {
		case *kafka.Message:
			if eMsg.TopicPartition.Error != nil {
				return err
			} else {
				log.Printf("message was delevered on topic %s, partition %s, offset %s",
					*eMsg.TopicPartition.Topic, eMsg.TopicPartition.Partition, eMsg.TopicPartition.Offset)
				return nil
			}
		}
	}
	return nil
}

func Consume(topic string) (<-chan kafka.Message, CloseFunction, error) {
	url := os.Getenv("KAFKA_URL")
	if url == "" {
		log.Println("KAFKA_URL is not defined")
	}

	log.Printf("KAFKA_URL = %s", url)
	config := kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          "template-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		err = fmt.Errorf("error while trying to create kafka-consumer: %v", err)
		return nil, nil, err
	}

	for {
		if err = consumer.SubscribeTopics([]string{topic}, nil); err != nil {
			err = fmt.Errorf("error while trying to subscibe to topic: %v", err)
			log.Printf("error while trying to subscibe to topic: %v", err)

			time.Sleep(time.Second * 2)
		} else {
			break
		}
	}

	msgCh := make(chan kafka.Message)
	go func() {
		for {
			log.Println("start reading message")

			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else {
				switch e := err.(type) {
				case *kafka.Error:
					if e.Code() != kafka.ErrTimedOut {
						log.Printf("Kafka error: %v (%v)\n", e.Code(), e)
					}
				default:
					log.Printf("Unexpected error: %v (%T)\n", err, err)
				}
				continue
			}
			msgCh <- *msg
		}
	}()

	return msgCh, consumer.Close, nil
}
