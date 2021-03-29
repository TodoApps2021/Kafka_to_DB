package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

// FOR TEST MESSAGE FOR TEST MESSAGE FOR TEST MESSAGE FOR TEST MESSAGE FOR TEST MESSAGE FOR TEST MESSAGE FOR TEST MESSAGE
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfig(); err != nil {
		log.Fatal("not initialized config")
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("kafka.url")})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	user := message.User{
		Name:     "Vlad",
		Username: "Vlad",
		Password: "Vlad",
	}

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)
	topic := "accounts"
	text := map[string]interface{}{
		"status": "create",
		"item":   user,
	}
	value, err := json.Marshal(text)
	if err != nil {
		log.Print(err.Error())
	}
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0},
		Value:          []byte(value),
		Key:            []byte(fmt.Sprint(time.Now().UTC())),
		// Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)
	log.Print(err)
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
