package kafka

import "github.com/TodoApps2021/Kafka_to_DB/pkg/kafka/consumer"

type ConsumerConfig struct {
	LogLevel  string
	InfoPort  int
	TopicName string
	Consumer  consumer.Config
}

// type ProducerConfig struct {
// 	LogLevel  string
// 	InfoPort  int
// 	TopicName string
// 	Producer  producer.Config
// }
