# Changelog

## Unreleased

* [ENHANCEMENT] Add unit and integration test suite with `make test`, `make test-integration`, and CI coverage.
* [ENHANCEMENT] Local and CI Kafka test environment uses `apache/kafka:4.3.0` (KRaft).

## 2.0.0 / 2026-06-03

* [FEATURE] Add `kafka_topic_partition_consumer` metric for topics without active consumers.
* [ENHANCEMENT] Batch topic offset fetches by leader broker.
* [ENHANCEMENT] Reuse topic offset map when collecting consumer group lag.
* [ENHANCEMENT] Use group coordinator for offset fetch.
