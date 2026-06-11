# Changelog

## 3.0.0 / 2026-06-11

This release marks the repository as a standalone project, fully detached from
[danielqsj/kafka_exporter](https://github.com/danielqsj/kafka_exporter).

**Breaking changes** (relative to upstream and to `v2.0.0`; see README "Breaking Changes from Upstream"):

* [CHANGE] Binary is named `kafka-exporter` (upstream and `v2.0.0`: `kafka_exporter`).
* [CHANGE] Docker image entrypoint and binary path use `kafka-exporter` (upstream and `v2.0.0`: `kafka_exporter`).
* [CHANGE] Go module path is `github.com/rfvbkm/kafka-exporter` (`v2.0.0`: `github.com/rfvbkm/kafka_exporter`).
* [CHANGE] Helm chart and Kubernetes manifests reference the new binary and image names.
* [CHANGE] Git history was rewritten to drop commits that included the `vendor/` directory and reduce repository size; clones and forks created before the rewrite are not compatible — re-clone instead of pulling.

Other changes since `v2.0.0`:

* [ENHANCEMENT] Rewrite the Grafana dashboard (`kafka_exporter_overview.json`) using the current Grafana schema (timeseries/stat/bargauge/table panels, datasource variable) and add panels for `kafka_topic_partition_consumer`: active consumers by topic/partition and partitions without an active consumer.
* [ENHANCEMENT] Add unit tests for certificate helpers, consumer metric filters, and OAuth token caching (`make test`).
* [ENHANCEMENT] Add integration tests that scrape a real Kafka broker (`make test-integration`, `make test-all`).
* [ENHANCEMENT] Add `dev/wait-kafka` helper and `make ensure-kafka` for local Kafka startup.
* [ENHANCEMENT] CI runs unit tests in the lint job and integration tests in a separate job with an `apache/kafka:4.3.0` service container.
* [ENHANCEMENT] Local and CI Kafka test environment uses `apache/kafka:4.3.0` (KRaft).
* [ENHANCEMENT] Document the unit and integration test workflow in README.
* [ENHANCEMENT] Refactor metric descriptor initialization into `initMetricDescs`.
* [BUGFIX] Abort `make test-integration` when Kafka never becomes ready instead of hanging.
* [BUILD] Promote `github.com/prometheus/client_model` to a direct dependency.
* [BUILD] Bump `github.com/IBM/sarama` from 1.47.0 to 1.50.2.
* [BUILD] Bump `github.com/panjf2000/ants/v2` from 2.11.3 to 2.12.1.
* [BUILD] Bump `golang.org/x/oauth2` from 0.32.0 to 0.36.0.
* [BUILD] Bump `k8s.io/klog/v2` from 2.130.1 to 2.140.0.
* [BUILD] Bump `github.com/xdg-go/scram` from 1.1.2 to 1.2.0.
* [BUILD] Bump `github.com/aws/aws-msk-iam-sasl-signer-go` from 1.0.0 to 1.0.4.
* [BUILD] Bump `github.com/prometheus/client_model` from 0.6.1 to 0.6.2.

## 2.0.0 / 2026-06-03

* [FEATURE] Add `kafka_topic_partition_consumer` metric for topics without active consumers.
* [ENHANCEMENT] Batch topic offset fetches by leader broker.
* [ENHANCEMENT] Reuse topic offset map when collecting consumer group lag.
* [ENHANCEMENT] Use group coordinator for offset fetch.
