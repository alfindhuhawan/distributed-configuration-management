# Architecture Explanation

This document expands on docs/memory_architecture.md and explains the repository architecture in detail: layout, responsibilities, runtime flow, and the way services communicate.

## 1. Repository Layout and Intent

This repository is a Go monorepo built around a Clean Architecture / Gogen-style organization. The code is split into distinct layers and bounded contexts so that use cases (business logic) stay separated from IO concerns (HTTP, databases, and config). The top-level directories are used as follows:

- application/ (composition root)
  - Wires configuration, logging, HTTP server, controllers, usecases, and gateways.
  - Each file represents a runnable application (controller, agent, worker).

- domain_controllerappservice/ (controller domain)
  - Handles HTTP endpoints for configuration and agent registration.
  - Implements usecases that read/write global config and agent registration data.

- domain_workerappservice/ (worker domain)
  - Handles HTTP endpoints that accept config updates and execute an outbound "hit".
  - Uses an in-memory SDK to store the current config and call an external URL.

- shared/
  - Infrastructure: config loader, database clients, server, logger, utilities.
  - Models: entities, repository contracts, request/response payloads, enums/errors.
  - Gateways: concrete persistence implementations (Mongo) and config access.
  - SDKs: agent package for polling/registration, worker SDK for in-memory config.

## 2. Entry Point and Application Selection

The application entrypoint is main.go. It reads configuration once and then decides which application to run based on the first CLI argument:

- go run main.go appcontroller
- go run main.go appagent
- go run main.go appworker

In main.go:

- config.ReadConfig() validates that config.json is readable and well-formed.
- A map of app factories connects an app name string to a function that builds that app.
- driver.Run() is called on the selected registry contract.

This design supports multiple services in one repo without a separate binary for each, and keeps wiring logic in application/ while domain logic stays independent.

## 3. Composition Root (application/)

Each file in application/ is a composition root. It creates concrete implementations and connects them to controllers and usecases.

Common pattern:

- Read config from config.json
- Generate application metadata (name, instance ID, start time)
- Create a logger (structured JSON logger)
- Initialize the Gin HTTP handler with CORS and a /ping route
- Wire the controller with the correct inports (usecases) and outports (gateways)

### appcontroller

- Starts a Gin HTTP server on the controller port from config.json.
- Instantiates prod gateway with Mongo support and config access.
- Wires REST handlers to usecases for agent registration and config management.

### appagent

- Starts a Gin HTTP server on the agent port from config.json.
- Runs background processes for agent registration and polling.
- Initializes agent_cache.json, restores last-known config URL if present.
- Calls the worker service when new config is detected.

### appworker

- Starts a Gin HTTP server on the worker port from config.json.
- Builds the in-memory worker SDK with a configurable timeout.
- Wires REST handlers for accepting config updates and performing outbound hits.

## 4. Controller Domain Service

Location: domain_controllerappservice/

The controller service is the source of truth for global configuration and agent registration. It exposes endpoints and persists data in MongoDB.

### REST API

- POST /api/v1/register
  - Validates agent payload (name and IP).
  - Persists agent metadata.
  - Initializes agent_cache.json values (agent id, poll interval, config URL).

- GET /api/v1/config
  - Reads global config.
  - Returns current URL, poll interval, and version.
  - Supports version checking (client sends ETag-like version).

- POST /api/v1/config
  - Updates the global configuration.
  - Increments version if URL or poll interval changes.

### Usecases

- registeragent
  - Validates input.
  - Saves a new agent entry.
  - Reads the global config (if any) to return current config settings.
  - Writes agent_cache.json so the agent can continue even if controller is down.

- getconfig
  - Looks up the global config by type (default is global).
  - Returns URL, poll interval, and version.
  - Indicates if the data changed compared to a version in the request.

- updateconfig
  - Validates input and updates the global config.
  - Creates a new config document if none exists.
  - Increments version only when values change.

### Gateway

The controller gateway is Mongo-backed. It combines several repository implementations into one shared data source:

- AgentImpl: reads/writes agent records.
- GlobalConfigImpl: reads/writes global configuration.
- ConfigImpl: supplies application config access where needed.

Mongo access is abstracted through shared/infrastructure/database and is injected into the gateway during app startup.

## 5. Worker Domain Service

Location: domain_workerappservice/

The worker service stores config in memory and can make outbound HTTP calls using that config.

### REST API

- POST /api/v1/config
  - Accepts a config update (URL) and stores it via the worker SDK.
  - Returns success or validation errors (URL required).

- GET /api/v1/hit
  - Uses the worker SDK to call the stored URL.
  - Returns the status code, response headers, and raw body bytes from the external service.

### Usecases

- runconfig
  - Validates and stores the new config URL in the worker SDK.

- getconfig
  - Calls worker SDK Hit() to perform the outbound request.

### Worker SDK

The SDK is located in shared/pkg/workersdk and is intentionally simple:

- Stores config in-memory with a read/write lock.
- Ensures a URL exists before performing any request.
- Uses a configurable timeout for outbound requests.
- Returns status code, headers, and raw body bytes for the last hit.

This separation allows the worker domain to remain independent of the HTTP transport and external network details.

## 6. Agent Service and Polling Logic

Location: shared/pkg/agent

The agent service is a background participant that registers itself, polls the controller for config changes, and informs the worker service when updates happen.

Key responsibilities:

- AutoRegister
  - Sends a signed request to the controller register endpoint.
  - Writes agent_cache.json with agent id, poll URL, version, interval.

- StartDynamicPolling
  - Performs periodic polling with a configurable interval.
  - Adjusts polling interval if controller returns a new one.

- pollData
  - Loads last known version from agent_cache.json.
  - Calls controller config endpoint with version in ETag header.
  - If a new version is available, updates cache and calls worker service.

- fetchConfigInfiniteBackoff
  - If polling fails, enters an exponential backoff loop.
  - Keeps retrying until controller is reachable.

This pattern lets the system tolerate controller restarts and transient outages while keeping worker configuration up to date.

## 7. Shared Infrastructure and Models

### Infrastructure

- shared/infrastructure/config
  - Reads config.json and agent_cache.json into strongly typed structs.

- shared/infrastructure/server
  - Builds the Gin HTTP server and wires CORS and /ping.
  - Includes graceful shutdown support.

- shared/infrastructure/database
  - MongoDB connection and transaction helpers.

- shared/infrastructure/logger
  - Structured JSON logger with app metadata and trace IDs.

- shared/infrastructure/util
  - Helpers for IDs, hashing, JSON parsing, cache read/write, and general utilities.

### Models and Contracts

- shared/model/entity
  - Core domain entities (Agent, GlobalConfig, etc.).

- shared/model/repository
  - Repository interfaces (SaveAgentRepo, FindOneGlobalConfigRepo, etc.).
  - Request/response structures used across layers.

- shared/model/request and shared/model/response
  - HTTP payload models for worker and controller endpoints.

- shared/model/errorenum
  - Error enums used to standardize errors across usecases.

These boundaries keep domain logic insulated from infrastructure changes.

## 8. Configuration and State Files

### config.json

Defines application behavior and wiring:

- Port assignments for controller, agent, worker servers.
- Credentials (admin/agent) for request signing.
- MongoDB connection information.
- Polling and backoff values (poll interval, base delay, max delay).
- Agent initial settings (agent name, IP, worker URL, cache path).

### agent_cache.json

Stores agent state used for recovery and polling:

- Agent ID and name
- Last version of config seen
- Poll interval and URL
- Last poll timestamp

This enables a warm restart of the agent even when the controller is unavailable.

## 9. End-to-End Runtime Flow

A typical system run looks like this:

1) Controller service is started.
   - Accepts register/config requests and stores global config.

2) Worker service is started.
   - Waits for config to be set.

3) Agent service is started.
   - If agent_cache.json exists, it pushes the last-known URL to the worker.
   - Registers with the controller to fetch poll config.
   - Starts polling loop for global config changes.

4) Config update occurs (admin POST /api/v1/config).
   - Controller increments version and stores new URL.

5) Agent detects new version.
   - Updates agent_cache.json.
   - Calls worker /api/v1/config with the new URL.

6) Worker now uses the updated URL.
   - /api/v1/hit triggers a GET to the external config URL and returns the result.

This flow creates a simple distributed configuration pipeline with separation of concerns between config source (controller), distribution/monitoring (agent), and execution (worker).

## 10. Design Notes and Tradeoffs

- Clean Architecture boundaries are enforced through inport/outport interfaces.
- The controller is the single source of truth for global config, simplifying versioning.
- The agent uses local cache to ensure resilience when controller is offline.
- The worker holds config in-memory, which is fast but not persistent across restarts.
- Versioning is tracked with a numeric value and ETag-like checks to reduce updates.

## 11. Files to Explore

- main.go
- application/app_controllerappservice.go
- application/app_agentappservice.go
- application/app_workerappservice.go
- domain_controllerappservice/controller/restapi/router.go
- domain_controllerappservice/usecase/registeragent/interactor.go
- domain_controllerappservice/usecase/getconfig/interactor.go
- domain_controllerappservice/usecase/updateconfig/interactor.go
- domain_workerappservice/controller/restapi/router.go
- domain_workerappservice/usecase/runconfig/interactor.go
- shared/pkg/agent/agent.go
- shared/pkg/workersdk/worker.go
- shared/infrastructure/config/config.go
