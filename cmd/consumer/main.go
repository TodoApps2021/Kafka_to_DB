package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/kafka/consumer"
	options "github.com/TodoApps2021/Kafka_to_DB/pkg/options/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := options.ReadKafkaConsumerEnv()

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	wg := new(sync.WaitGroup)

	handler := &Handler{}

	client, err := consumer.New(config.Consumer, []string{config.TopicName}, handler)
	if err != nil {
		log.Print("kafka consumer init error:", err.Error())
		os.Exit(1)
	}

	client.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Print("end")
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("got Interrupt signal")
		stop()
	}()
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, key, value []byte, timestamp time.Time, p kafka.TopicPartition) error {
	log.Printf("key: <%s>, value: <%s>, partition: %s", string(key), string(value), p.String())
	return nil
}
