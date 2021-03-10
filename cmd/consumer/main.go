package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/kafka/consumer"
	options "github.com/TodoApps2021/Kafka_to_DB/pkg/options/kafka"
)

func main() {
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
	log.Info("end")
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

func (h *Handler) Handle(ctx context.Context, key, value []byte, timestamp time.Time) error {
	log.Printf("key:%s, value:%s", string(key), string(value))
	return nil
}
