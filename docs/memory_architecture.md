# Repository Architecture Memory

## Overview
- Go monorepo structured around a Clean Architecture/Gogen-style layout.
- Three runnable applications selected via `go run main.go <app_name>`:
  - `appcontroller`: controller service (HTTP API + persistence).
  - `appagent`: agent service (auto-register + polling loop).
  - `appworker`: worker service (in-memory config + outbound hit).
- Configuration is JSON-driven (`config.json`) and read at startup.

## Entry Point + App Wiring
- `main.go` loads config, maps CLI arg -> application factory, then calls `driver.Run`.
- `application/` constructs each app:
  - reads config
  - creates app metadata (name/instance/time)
  - builds logger + Gin HTTP handler
  - wires controller/usecases/gateways

## Domain: Controller App Service
- Location: `domain_controllerappservice/`
- REST API routes: `/api/v1/register`, `/api/v1/config` (GET/POST).
- Usecases:
  - `registeragent` saves an agent and initializes agent cache metadata.
  - `getconfig` returns global config + versioning info (ETag-style flow).
  - `updateconfig` writes global config, increments version if changed.
- Gateway (prod): Mongo-backed implementations for agents/config/global config.

## Domain: Worker App Service
- Location: `domain_workerappservice/`
- REST API routes: `/api/v1/config` (POST), `/api/v1/hit` (GET).
- Usecases:
  - `runconfig` updates the worker SDK with a new config URL.
  - `getconfig` calls the worker SDK `Hit()` which makes the outbound request.
- Gateway (prod): Mongo wrapper (minimal here; worker uses in-memory SDK).

## Shared Layer
- `shared/infrastructure/`: config loader, Gin server, MongoDB connector, logger, util.
- `shared/model/`: entities, repositories, request/response payloads, errors/enums.
- `shared/gateway/prod/`: concrete repository implementations (Mongo + config access).
- `shared/pkg/agent`: auto-register + polling + backoff; manages agent_cache.json.
- `shared/pkg/workersdk`: in-memory config store + outbound HTTP hit.

## Runtime Data Flow (High-Level)
1. Controller service exposes register/config APIs and stores global config in Mongo.
2. Agent service auto-registers with controller, writes `agent_cache.json`, and polls with version/ETag.
3. On config changes, agent pushes config URL to worker via `/api/v1/config`.
4. Worker stores URL in-memory; `/api/v1/hit` triggers outbound GET to that URL.

## Key Config + State
- `config.json` defines ports, credentials, MongoDB, polling/backoff, and agent init settings.
- `agent_cache.json` stores agent id/name, last version, poll interval, and last poll time.
