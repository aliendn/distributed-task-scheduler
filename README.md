# Distributed Task Scheduler (Go + Gin + GORM)

This is a distributed, priority-based task scheduling system written in Go. It features leader election,
horizontal scalability, persistence with PostgreSQL, and Prometheus-based observability.

## ✅ Features

- Task priority levels: High, Medium, Low
- REST API to submit and query tasks
- Worker pool with backpressure handling
- Leader election (pluggable)
- PostgreSQL persistence using GORM
- Prometheus metrics endpoint (`/metrics`)
- Docker + Docker Compose for easy deployment

## 🔧 Tech Stack

- Go 1.21
- Gin (HTTP server)
- GORM (ORM)
- PostgreSQL
- Prometheus
- Docker

## 📈 Observability

This project includes full Prometheus integration:

- `/metrics` exposed by the Gin server
- Monitored by Prometheus via Docker Compose
- Metrics include:
    - `task_queue_length`
    - `task_submitted_total`
    - `task_processed_total`
    - `task_processing_seconds`

### Access Prometheus

Open: [http://localhost:9090](http://localhost:9090)

Example queries:
```promql
task_queue_length
task_processed_total
rate(task_processing_seconds_sum[1m])
```

### Run with Docker:

```bash
docker-compose up --build
```
