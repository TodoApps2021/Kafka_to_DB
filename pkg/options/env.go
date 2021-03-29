package options

import "github.com/spf13/viper"

type ConfigConsumer struct {
	BootstrapServers string
	GroupId          string
	SessionTimeoutMs int
	AutoOffsetReset  string
	EnableAutoCommit bool
}

func ReadFromYAML() *ConfigConsumer {
	return &ConfigConsumer{
		BootstrapServers: viper.GetString("kafka.url"),
		GroupId:          viper.GetString("kafka.consumer.group_id"),
		SessionTimeoutMs: viper.GetInt("kafka.consumer.session_timeout_ms"),
		AutoOffsetReset:  viper.GetString("kafka.consumer.auto_offset_reset"),
		EnableAutoCommit: viper.GetBool("kafka.consumer.enable_auto_commit"),
	}
}
