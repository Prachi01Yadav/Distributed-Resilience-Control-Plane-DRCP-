# DRCP — Distributed Resilience Control Plane

A production-grade, distributed control plane that enforces Service Level Agreements (SLAs) across microservice architectures through real-time telemetry ingestion, automated circuit breaking, and immutable blockchain audit trails.

---

## Architecture

```
                    ┌─────────────┐
                    │  Dashboard  │  (HTML/CSS/JS)
                    └──────┬──────┘
                           │ HTTP
                    ┌──────▼──────┐
                    │   Gin API   │  :8080
                    │  (Go/REST)  │
                    └──┬───┬───┬──┘
                       │   │   │
              ┌────────┘   │   └────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │PostgreSQL│ │  Kafka   │ │  Redis   │
        │ Registry │ │ Telemetry│ │ Budget   │
        └──────────┘ └────┬─────┘ └──────────┘
                          │
                    ┌─────▼─────┐
                    │  Worker   │
                    │ (OPA Eval)│
                    └─────┬─────┘
                          │
                ┌─────────┴─────────┐
                ▼                   ▼
         ┌────────────┐     ┌────────────┐
         │ xDS Server │     │  Ethereum  │
         │ (Envoy CP) │     │  Anchor    │
         └────────────┘     └────────────┘
```

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| API Framework | Gin |
| Database | PostgreSQL + GORM |
| Messaging | Apache Kafka (Sarama) |
| Cache | Redis |
| Policy Engine | Open Policy Agent (OPA) |
| Service Mesh | Envoy Proxy (xDS gRPC) |
| Blockchain | Ethereum (go-ethereum) |
| Frontend | Vanilla HTML/CSS/JS + Lucide Icons |
| Infrastructure | Docker, Terraform, Helm, Kubernetes |

## Features

- **Service Registry** — Register, discover, and manage microservices with full CRUD
- **SLA Contract Management** — Define latency and error-rate budgets per service using OPA/Rego policies
- **Real-Time Telemetry Pipeline** — Ingest metrics via Kafka, compute sliding-window error budgets in Redis
- **Automated Circuit Breaking** — Trip Envoy proxy circuit breakers via xDS when SLA thresholds are breached
- **Blockchain Audit Trail** — Anchor breach incidents as immutable transactions on Ethereum
- **Incident Tracking** — Full lifecycle tracking: OPEN → RESOLVED → ANCHORED
- **Modern Dashboard** — Beautiful, responsive UI with glassmorphism design, real-time data, and toast notifications

## Project Structure

```
sentinelmesh/
├── cmd/
│   ├── api/          # REST API server (main entry point)
│   ├── worker/       # Kafka consumer + OPA evaluation worker
│   ├── xds/          # Envoy xDS gRPC control plane
│   └── anchor/       # Blockchain anchoring microservice
├── internal/
│   ├── registry/     # Service & SLA CRUD (models, handlers, repository)
│   ├── telemetry/    # Telemetry ingestion handler
│   ├── budget/       # Redis sliding-window error budget calculator
│   ├── policy/       # OPA policy evaluation engine
│   ├── xds/          # Envoy snapshot builder
│   └── anchor/       # Ethereum transaction signer
├── pkg/
│   ├── db/           # PostgreSQL connection (GORM)
│   ├── kafka/        # Kafka producer & consumer
│   ├── cache/        # Redis client wrapper
│   └── logger/       # Zap structured logging
├── web/              # Frontend dashboard (HTML, CSS, JS)
├── contracts/        # Solidity smart contracts
├── deployments/      # Terraform & Helm charts
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Quick Start

### Prerequisites
- Go 1.23+
- Docker Desktop

### Run Locally

```bash
# 1. Start infrastructure
docker compose up -d

# 2. Run the API server
CGO_ENABLED=0 go run ./cmd/api

# 3. Open the dashboard
# http://localhost:8080
```

### API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/api/v1/services` | List all services |
| `POST` | `/api/v1/services` | Register a service |
| `GET` | `/api/v1/services/:id` | Get service by ID |
| `POST` | `/api/v1/services/:id/contracts` | Create SLA contract |
| `GET` | `/api/v1/contracts` | List all SLA contracts |
| `GET` | `/api/v1/incidents` | List all incidents |
| `POST` | `/api/v1/telemetry` | Ingest telemetry data |

## Deployment

### Docker

```bash
docker build -t drcp .
docker run -p 8080:8080 -e DATABASE_URL=<your-postgres-url> drcp
```

### Railway / Render

This project includes a `Dockerfile` for one-click deployment on Railway or Render. Connect your GitHub repo and the platform will auto-detect the Dockerfile.

## License

by: Prachi01Yadav
