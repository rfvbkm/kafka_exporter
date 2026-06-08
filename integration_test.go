package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestIntegrationSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	brokers := kafkaBrokersFromEnv()
	if !kafkaAvailable(brokers) {
		if os.Getenv("REQUIRE_KAFKA") == "1" {
			t.Fatal("Kafka is not available at " + strings.Join(brokers, ","))
		}
		t.Skip("Kafka is not running, skipping integration test")
	}

	ensureTestTopic(t, brokers)

	const listenAddress = "127.0.0.1:19304"
	baseURL := "http://" + listenAddress

	initMetricDescs(nil)
	opts := kafkaOpts{
		uri:                     brokers,
		kafkaVersion:            "4.3.0",
		metadataRefreshInterval: "30s",
		groupMetricsTimeout:     "5m",
	}
	exporter, err := NewExporter(opts, ".*", "^$", ".*", "^$")
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	t.Cleanup(func() {
		exporter.client.Close()
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	})

	waitForEndpoint(t, baseURL+"/healthz", func(body string, status int) bool {
		return status == http.StatusOK && body == "ok"
	})

	var metricsBody string
	waitForEndpoint(t, baseURL+"/metrics", func(body string, status int) bool {
		if status != http.StatusOK {
			return false
		}
		metricsBody = body
		return strings.Contains(body, "kafka_brokers") && strings.Contains(body, "kafka_topic_partitions")
	})

	if !strings.Contains(metricsBody, "kafka_brokers") {
		t.Fatal("metrics body missing kafka_brokers")
	}
	if !strings.Contains(metricsBody, "kafka_topic_partitions") {
		t.Fatal("metrics body missing kafka_topic_partitions")
	}
}

func ensureTestTopic(t *testing.T, brokers []string) {
	t.Helper()

	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("4.3.0")
	if err != nil {
		t.Fatalf("parse kafka version: %v", err)
	}
	cfg.Version = version

	admin, err := sarama.NewClusterAdmin(brokers, cfg)
	if err != nil {
		t.Fatalf("create cluster admin: %v", err)
	}
	defer admin.Close()

	const topic = "kafka-exporter-integration-test"
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil && !errors.Is(err, sarama.ErrTopicAlreadyExists) {
		t.Fatalf("create topic: %v", err)
	}
}

func kafkaBrokersFromEnv() []string {
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		return strings.Split(brokers, ",")
	}
	return []string{"localhost:9092"}
}

func kafkaAvailable(brokers []string) bool {
	client, err := sarama.NewClient(brokers, nil)
	if err != nil {
		return false
	}
	defer client.Close()
	_, err = client.Topics()
	return err == nil
}

func waitForEndpoint(t *testing.T, url string, ok func(body string, status int) bool) {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if ok(string(body), resp.StatusCode) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("endpoint %s did not become ready in time", url)
}
