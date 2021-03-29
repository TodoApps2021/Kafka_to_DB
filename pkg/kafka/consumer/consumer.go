package consumer

import (
	"context"
	"errors"
	"sync"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/options"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

// var (
// 	commands = map[string]func() error{
// 		"auth":  func() error { return nil },
// 		"lists": func() error { return nil },
// 		"items": func() error { return nil },
// 	}
// )

type Handler interface {
	// If 'Handle' returns an error, a message will not be committed.
	Handle(ctx context.Context, key, value []byte, topic string, partition int32) error
}

type Consumer struct {
	consumer      *kafka.Consumer
	handler       Handler
	runErr        error
	topics        []string
	name          string
	pollTimeoutMs int
}

func New(config *options.ConfigConsumer, topics []string, handler Handler) (*Consumer, error) {
	configMap := kafka.ConfigMap{
		"bootstrap.servers":  config.BootstrapServers,
		"group.id":           config.GroupId,
		"session.timeout.ms": config.SessionTimeoutMs,
		"auto.offset.reset":  config.AutoOffsetReset,
		"enable.auto.commit": config.EnableAutoCommit,
	}

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:      consumer,
		topics:        topics,
		handler:       handler,
		name:          "todo_app",
		pollTimeoutMs: 300,
	}, nil
}

func (c *Consumer) Run(ctx context.Context, wg *sync.WaitGroup) {
	log.Infof("kafka consumer '%s': begin run", c.name)
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		log.Errorf("kafka consumer '%s': end run: failed to subscribe to topics: %v", c.name, err)
		c.runErr = err
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Poll(ctx)
		c.close()
		log.Infof("kafka consumer '%s': end run", c.name)
	}()
}

func (c *Consumer) Poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := c.consumer.Poll(c.pollTimeoutMs)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				topic := e.TopicPartition.Topic
				partition := e.TopicPartition.Partition

				err := c.handler.Handle(ctx, e.Key, e.Value, *topic, partition)
				if err != nil {
					log.Errorf("kafka consumer '%s': failed to handle message: %v", c.name, err)
					continue
				}
				_, err = c.consumer.CommitMessage(e)
				if err != nil {
					log.Errorf("kafka consumer '%s': failed to commit message %v: %v", c.name, e, err)
				}
			case kafka.Error:
				if e.IsFatal() {
					log.Errorf("kafka consumer '%s': fatal error: %v", c.name, e)
					c.runErr = e
					return
				}
				log.Tracef("kafka consumer '%s': error: %v", c.name, e)
			}
		}
	}
}

func (c *Consumer) HealthCheck() error {
	if c.runErr != nil {
		return errors.New("kafka consumer " + c.name + ": run issue: " + c.runErr.Error())
	}
	return nil
}

func (c *Consumer) close() {
	if err := c.consumer.Close(); err != nil {
		log.Errorf("kafka consumer '%s': failed to close: %v", c.name, err)
	}
}
