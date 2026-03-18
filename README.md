# Nexus

> A high-performance distributed event streaming platform built in Go.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/alex-corbin-dev/nexus/actions)
[![Coverage](https://img.shields.io/badge/coverage-87%25-green)](https://github.com/alex-corbin-dev/nexus)

Nexus is a cloud-native event streaming platform designed for high-throughput, low-latency workloads. It supports multi-region replication, exactly-once delivery semantics, and sub-millisecond end-to-end latency at scale.

## Features

- **Exactly-once delivery** — Idempotent producers and transactional consumers guarantee no duplicate or lost messages
- **Multi-region replication** — Active-active replication across geographic regions with configurable consistency levels
- **Sub-millisecond latency** — Optimised I/O path using io_uring and zero-copy networking
- **Schema registry** — Built-in Avro/Protobuf schema management with full backward/forward compatibility checks
- **Consumer groups** — Kafka-compatible consumer group protocol with automatic partition rebalancing
- **Observability** — Native OpenTelemetry integration with distributed tracing and Prometheus metrics

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Nexus Cluster                        │
│                                                             │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐               │
│  │  Broker  │   │  Broker  │   │  Broker  │               │
│  │  Node 1  │◄──│  Node 2  │──►│  Node 3  │               │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘               │
│       │              │              │                       │
│  ┌────▼──────────────▼──────────────▼────┐                 │
│  │           Raft Consensus Layer         │                 │
│  └───────────────────────────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/alex-corbin-dev/nexus.git
cd nexus

# Start a local 3-node cluster
make cluster-up

# Produce a message
./bin/nexus-cli produce --topic events --message '{"user_id": 42, "action": "login"}'

# Consume messages
./bin/nexus-cli consume --topic events --group my-consumer-group
```

## Configuration

```yaml
# nexus.yaml
cluster:
  node_id: "broker-1"
  listen_addr: "0.0.0.0:9092"
  advertised_addr: "broker-1.example.com:9092"

replication:
  factor: 3
  min_in_sync_replicas: 2

storage:
  data_dir: "/var/lib/nexus"
  segment_size_mb: 256
  retention_hours: 168

observability:
  metrics_port: 9090
  tracing_endpoint: "http://jaeger:14268/api/traces"
```

## Performance Benchmarks

| Scenario | Throughput | P99 Latency |
|---|---|---|
| Single producer, single partition | 1.2M msg/s | 0.8ms |
| 10 producers, 100 partitions | 8.4M msg/s | 1.2ms |
| Cross-region replication (us-east → eu-west) | 450K msg/s | 12ms |

## Tech Stack

- **Language:** Go 1.22
- **Consensus:** Raft (via `hashicorp/raft`)
- **Storage:** Custom LSM-tree with WAL
- **Networking:** `io_uring` via `iceber/iouring-go`
- **Serialisation:** Protocol Buffers
- **Observability:** OpenTelemetry, Prometheus

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please ensure all tests pass and coverage remains above 85%.

## License

[Apache 2.0](LICENSE)
