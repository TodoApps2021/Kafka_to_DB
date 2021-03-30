# Kafka_to_DB

This service accepts all tasks (from bot telegrams and from the site) and sends them to the Postgres SQL database.

- Kafka
    - Consumer (topics: "auth", "todo_list", "todo_item")
- Postgres SQL

## Libs:

[confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)  
[pgx/pgxpool](https://github.com/jackc/pgx)  
[logrus](https://github.com/sirupsen/logrus)  
[viper](https://github.com/spf13/viper)

### Developed: Vlad1k_Zhuk