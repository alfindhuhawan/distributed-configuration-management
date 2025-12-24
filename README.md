# Distributed Configuration Management

This repository is a Go monorepo that implements a small distributed configuration system with three runnable services:

- Controller service: source of truth for global configuration and agent registration.
- Agent service: auto-registers, polls for config changes, and pushes updates to workers.
- Worker service: stores the latest config URL in memory and can call it on demand.

It follows a Clean Architecture / Gogen-style layout with explicit usecases, inports/outports, and infrastructure wiring.

---

## Tech Stack

- Language: Go
- Web framework: Gin (HTTP router + middleware)
- Persistence: MongoDB (controller service)
- Architecture: Clean Architecture (Gogen-style layering)
- Config: JSON-based config files (config.json, agent_cache.json)
- Auth: SHA256 signature in Authorization: Bearer <hash> header (You can get the signature by, encrypt SHA256 credentials.admin.secret_key or credentials.agent.secret_key from config.json)

---

## Repository Structure

High-level layout:

```bash
.
|-- application/                # Composition roots for each service
|-- domain_controllerappservice # Controller domain: REST + usecases
|-- domain_workerappservice     # Worker domain: REST + usecases
|-- shared/                     # Infra, models, gateways, SDKs
|-- config.sample.json          # Sample config
|-- main.go                     # Application entrypoint (selects service)
|-- go.mod
|-- go.sum
```

---

## How To Run

### Prerequisites

- Go 1.25+ (or your local Go version)
- MongoDB (required for controller persistence)

### 1) Configure

Copy and edit the config:

```bash
cp config.sample.json config.json
```

Update ports, MongoDB settings, credentials, and agent init values as needed.

### 2) Start Services

Run each service in its own terminal:

```bash
go run main.go appcontroller
go run main.go appworker
go run main.go appagent
```

Default ports (from config.json):

- Controller: 8080
- Agent: 8081
- Worker: 8082

### 3) Key API Endpoints

Controller:

- POST /api/v1/register (agent registration)
- GET /api/v1/config (get global config with ETag header)
- POST /api/v1/config (update global config)

Worker:

- POST /api/v1/config (update worker config URL)
- GET /api/v1/hit (call the configured URL)

OpenAPI spec is available at docs/openapi/openapi.yaml.
Postman JSON is available at docs/postman/distributed-configuration-management.json.
