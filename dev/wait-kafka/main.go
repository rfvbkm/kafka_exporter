package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		brokers = strings.Split(v, ",")
	}

	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		client, err := sarama.NewClient(brokers, nil)
		if err == nil {
			_, err = client.Topics()
			_ = client.Close()
			if err == nil {
				return
			}
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Fprintf(os.Stderr, "Kafka not ready at %s\n", strings.Join(brokers, ","))
	os.Exit(1)
}
