package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/handler"
	"github.com/TodoApps2021/Kafka_to_DB/pkg/kafka/consumer"
	"github.com/TodoApps2021/Kafka_to_DB/pkg/options"
	"github.com/TodoApps2021/Kafka_to_DB/pkg/repository"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if err := initConfig(); err != nil {
		log.Fatal("Error init")
	}

	config := options.ReadFromYAML()
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	wg := new(sync.WaitGroup)

	pool, err := repository.NewPostgresDB(repository.Config{
		DB_URL: viper.GetString("db.url"),
	})
	if err != nil {
		log.Fatal("Error init pool DB")
	}

	handler := &handler.Handler{Repo: repository.NewRepository(pool)}

	client, err := consumer.New(config, viper.GetStringSlice("kafka.consumer.topics"), handler)
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
