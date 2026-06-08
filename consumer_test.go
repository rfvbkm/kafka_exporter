package main

import (
	"regexp"
	"sync"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNormalizeConsumerHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "/host.example", want: "host.example"},
		{in: "host.example", want: "host.example"},
		{in: "", want: ""},
		{in: "/127.0.0.1", want: "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			if got := normalizeConsumerHost(tt.in); got != tt.want {
				t.Fatalf("normalizeConsumerHost(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTryEmitTopicPartitionConsumer(t *testing.T) {
	initTestMetricDescs()

	topicFilter := regexp.MustCompile("^allowed-.*")
	topicExclude := regexp.MustCompile("^$")
	exporter := &Exporter{
		topicFilter:  topicFilter,
		topicExclude: topicExclude,
	}

	ch := make(chan prometheus.Metric, 4)
	seen := make(map[consumerMetricKey]struct{})
	var mu sync.Mutex

	emit := func(topic string) {
		exporter.tryEmitTopicPartitionConsumer(
			ch, seen, &mu,
			topic, "0", "group-a", "member-1", "host", "client-1", 1,
		)
	}

	emit("allowed-topic")
	emit("allowed-topic")
	emit("denied-topic")

	close(ch)

	var metrics []*dto.Metric
	for m := range ch {
		metrics = append(metrics, metricToDTO(t, m))
	}

	if len(metrics) != 1 {
		t.Fatalf("got %d metrics, want 1 (filter + dedup)", len(metrics))
	}

	m := metrics[0]
	if m.GetGauge().GetValue() != 1 {
		t.Fatalf("metric value = %v, want 1", m.GetGauge().GetValue())
	}
	labels := labelMap(m)
	if labels["topic"] != "allowed-topic" {
		t.Fatalf("topic label = %q, want %q", labels["topic"], "allowed-topic")
	}
	if labels["consumergroup"] != "group-a" {
		t.Fatalf("consumergroup label = %q, want %q", labels["consumergroup"], "group-a")
	}
}

func TestTryEmitTopicPartitionConsumerExcludedTopic(t *testing.T) {
	initTestMetricDescs()

	exporter := &Exporter{
		topicFilter:  regexp.MustCompile(".*"),
		topicExclude: regexp.MustCompile("^internal-.*"),
	}

	ch := make(chan prometheus.Metric, 1)
	seen := make(map[consumerMetricKey]struct{})
	var mu sync.Mutex

	exporter.tryEmitTopicPartitionConsumer(
		ch, seen, &mu,
		"internal-topic", "0", "group-a", "member-1", "host", "client-1", 1,
	)
	close(ch)

	if len(ch) != 0 {
		t.Fatal("expected no metric for excluded topic")
	}
}

func metricToDTO(t *testing.T, m prometheus.Metric) *dto.Metric {
	t.Helper()
	var out dto.Metric
	if err := m.Write(&out); err != nil {
		t.Fatalf("write metric: %v", err)
	}
	return &out
}

func labelMap(m *dto.Metric) map[string]string {
	labels := make(map[string]string, len(m.Label))
	for _, l := range m.Label {
		labels[l.GetName()] = l.GetValue()
	}
	return labels
}
