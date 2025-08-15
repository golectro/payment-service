package config

import (
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ensureKafkaTopic(viper *viper.Viper, log *logrus.Logger) {
	brokers := viper.GetStringSlice("KAFKA_BROKERS")
	topic := viper.GetString("KAFKA_TOPIC")

	if len(brokers) == 0 {
		log.Fatal("KAFKA_BROKERS is not set in configuration")
	}
	if topic == "" {
		log.Fatal("KAFKA_TOPIC is not set in configuration")
	}

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		log.Fatalf("Failed to connect to Kafka broker: %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		log.Infof("Kafka topic '%s' already exists", topic)
		return
	}

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("Failed to get Kafka controller: %v", err)
	}

	ctrlConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port)))
	if err != nil {
		log.Fatalf("Failed to connect to Kafka controller: %v", err)
	}
	defer ctrlConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	if err := ctrlConn.CreateTopics(topicConfigs...); err != nil {
		log.Fatalf("Failed to create Kafka topic '%s': %v", topic, err)
	}

	log.Infof("Kafka topic '%s' created successfully", topic)
}

func NewKafkaWriter(viper *viper.Viper, log *logrus.Logger) *kafka.Writer {
	ensureKafkaTopic(viper, log)

	brokers := viper.GetStringSlice("KAFKA_BROKERS")
	topic := viper.GetString("KAFKA_TOPIC")

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	log.Infof("Kafka writer connected to brokers: %v, topic: %s", brokers, topic)
	return writer
}

func NewKafkaReader(viper *viper.Viper, log *logrus.Logger) *kafka.Reader {
	ensureKafkaTopic(viper, log)

	brokers := viper.GetStringSlice("KAFKA_BROKERS")
	topic := viper.GetString("KAFKA_TOPIC")
	groupID := viper.GetString("KAFKA_GROUP_ID")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})

	log.Infof("Kafka reader connected to brokers: %v, topic: %s, groupID: %s", brokers, topic, groupID)
	return reader
}
