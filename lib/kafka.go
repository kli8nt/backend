package lib

import (
	"fmt"
	"log"
	"os"

	KF "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/joho/godotenv/autoload"
)

var (
	partition = 0
)

type KConfig struct {
	host     string
	port     string
	deadline int
	topic    string
}

type Kf struct {
	config   *KF.ConfigMap
	topic    string
	cleanups []func()
}

type KReader struct {
	topic string
	c     *KF.Consumer
}

var Kafka Kf

func onError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func init() {
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaUser := os.Getenv("KAFKA_USERNAME")
	kafkaPass := os.Getenv("KAFKA_PASSWORD")

	fmt.Println("Kafka Port: ", kafkaPort)
	fmt.Println("Kafka Host: ", kafkaHost)
	fmt.Println("Kafka User: ", kafkaUser)
	fmt.Println("Kafka Pass: ", kafkaPass)

	config := &KF.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s:%s", kafkaHost, kafkaPort),
		"group.id":           "myGroup",
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"sasl.username":      kafkaUser,
		"sasl.password":      kafkaPass,
		"session.timeout.ms": "45000",
		"auto.offset.reset":  "earliest",
	}

	log.Println("Connecting to Kafka...")

	Kafka.Init(config)
	defer Kafka.Close()
}

func (kf *Kf) Init(config *KF.ConfigMap) {
	kf.config = config
}

func (kf *Kf) Reader(topic string) KReader {
	c, err := KF.NewConsumer(kf.config)

	onError(err, "Failed to create consumer")

	kf.cleanups = append(kf.cleanups, func() {
		c.Close()
	})

	fmt.Printf("Created Consumer %v\n", c)

	c.Assign([]KF.TopicPartition{{Topic: &topic, Partition: int32(partition)}})
	kReader := KReader{
		topic: topic,
		c:     c,
	}

	return kReader
}

func (kReader *KReader) Close() (func()) {
	return func() {
		kReader.c.Close()
	}
}

func (kReader *KReader) Read(cb func(key string, msg []byte), OnError func(err error, msg string)) {
	log.Printf(" [*] Waiting for logs. for Application %s", kReader.topic)

	for {
		msg, err := kReader.c.ReadMessage(-1)
		if err == nil {
			cb(string(msg.Key), msg.Value)
		} else {
			onError(err, "Consumer error")
			OnError(err, "Failed to listen to this message")
		}
	}

	// 	for {
	// 		message, err := kReader.connection.ReadMessage(context.Background())
	// 		OnError(err, "Failed to listen to this message")

	//		key := string(message.Key)
	//		value :=message.Value
	//		cb(key, value)
	//	}
}

func (kf *Kf) Close() {
	for _, cleanup := range kf.cleanups {
		cleanup()
	}
}
