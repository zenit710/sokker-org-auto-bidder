# Architectural Philosophy

The project prioritizes:

- internal testability,
- simplicity,
- predictability,
- maintainability,
- low operational complexity.

The architecture should allow validating application logic without participating in real Sokker.org auctions.

All critical business flows should be testable in isolation using mocks and controlled environments.

The project intentionally avoids enterprise-grade abstractions and infrastructure complexity.

## Project Structure

Code is organized by responsibility.

Preferred structure:

```
/cmd/main.go
/internal/client
/internal/model
/internal/repository
/internal/service
/internal/subcommands
/tools
```

### Responsibilities
#### /cmd

Application bootstrap layer:

- dependency wiring,
- environment loading,
- logger configuration,
- CLI registration,
- application startup.

#### /internal/client

HTTP communication layer:

- authentication,
- session handling,
- requests to Sokker.org,
- HTTP response parsing.

#### /internal/model

Domain and transport models.

Models should remain simple data structures without infrastructure logic.

#### /internal/repository

Persistence layer.

Responsible for:

- database access,
- SQL queries,
- persistence logic.

#### /internal/service

Business logic layer.

Responsible for:

- auction automation,
- bidding decisions,
- orchestration of repositories and HTTP clients.

#### /internal/subcommands

CLI command implementations.

Examples:

- add player to bidding,
- validate authorization,
- run bidding process.

#### /tools

Reusable utility functions shared across the project.

Tools may be imported into internal packages.

## Dependency Rules

The project uses dependency injection through explicit constructor wiring.

The project does NOT use dependency injection frameworks.

Dependencies should remain explicit and traceable.

Preferred dependency direction:

- services depend on repositories and clients,
- repositories depend on database layer,
- subcommands depend on services,
- models should not depend on infrastructure layers,
- tools may be imported into internal packages.

Avoid circular dependencies.

## Interface Philosophy

Interfaces are used only when they provide clear value.

Primary reasons for interfaces:

- isolated testing,
- mocking external systems,
- preventing unnecessary HTTP requests in tests,
- avoiding duplication of test logic.

Do NOT create interfaces:

- prematurely,
- for every struct,
- for future hypothetical implementations.

Prefer concrete implementations unless abstraction is clearly justified.

## Persistence Philosophy

Repositories are responsible for all database communication.

SQLite is treated as a replaceable implementation detail.

The architecture should allow replacing SQLite with another storage engine by replacing repository implementations.

Raw SQL is allowed and preferred when simpler and more explicit than additional abstractions.

Avoid unnecessary ORM complexity.

## Error Handling Philosophy

Errors should provide:

- clear information for the end user,
- full contextual information for debugging.

Avoid generic errors whenever possible.

Business errors should:

- describe exactly what failed,
- contain actionable context.

Internal errors should:

- preserve original error chain,
- include operation context,
- support debugging and tracing.

Do not use panic for normal control flow.

## Logging Principles

Logging is an important operational and debugging mechanism.

The application logs:

- application initialization,
- service initialization,
- repository initialization,
- client initialization,
- CLI command execution,
- outgoing HTTP requests,
- successful HTTP communication,
- database operations,
- business logic execution,
- operational warnings,
- runtime errors.

## Log Levels

### trace

Detailed operational flow:

- service initialization,
- repository initialization,
- HTTP requests.

### debug

Technical debugging details:

- HTTP response statuses,
- diagnostic information.

### info

User-visible operational actions:

- selected CLI commands,
- application flow milestones.

### warn

Recoverable problems that do not stop execution.

### error

Failures causing invalid behavior or application termination.

## Security Logging Rules

Never log:

- credentials,
- authentication tokens,
- session cookies,
- secrets,
- sensitive user data.

Sensitive data must be sanitized before logging.

## Configuration Philosophy

Configuration is environment-variable driven.

Rules:

- no hardcoded secrets,
- minimal configuration surface,
- secure defaults,
- environment-first configuration.

The MVP should avoid unnecessary configuration systems.

## Concurrency Philosophy

Concurrency is used only when:

- behavior remains predictable,
- debugging complexity stays manageable,
- implementation remains understandable.

Prefer synchronous execution unless concurrency provides clear practical value.

Avoid hidden goroutines and uncontrolled async behavior.

## Testing Philosophy

Testability is one of the highest architectural priorities.

The project primarily uses:

- unit tests,
- isolated logic testing,
- mocked HTTP communication,
- mocked database interactions.

The goal is to validate business logic without:

- real HTTP communication,
- participation in real auctions,
- dependence on external systems.

Tests should remain:

- deterministic,
- isolated,
- fast,
- easy to debug.


## Security Principles

The application handles Sokker.org credentials and must follow strict security rules.

Rules:

- never log credentials,
- credentials only from environment variables,
- minimize secret persistence,
- sanitize logs,
- use secure defaults.


## Performance Philosophy

The project prioritizes:

- correctness,
- reliability,
- maintainability.

Rules:

- correctness over optimization,
- simplicity over micro-optimization,
- optimize only after measurement.

Avoid premature optimization.

## Dependency Management

Dependencies must remain minimal and intentional.

Rules:

- prefer standard library solutions,
- avoid dependencies for trivial functionality,
- every dependency must have clear justification.

Avoid dependency bloat.

## Operational Philosophy

The application should remain operationally simple.

Rules:

- single binary preferred,
- no external infrastructure,
- local execution first,
- graceful failure handling,
- restart-safe behavior.

The application should be easy to deploy and maintain on a local machine or simple server.

## Forbidden Patterns

The following patterns are discouraged or forbidden:

- global mutable state,
- hidden goroutines,
- reflection unless clearly necessary,
- dependency injection frameworks,
- panic-based control flow,
- unnecessary abstractions,
- speculative architecture,
- overengineering.

## Decision Rules

When architectural tradeoffs appear, prefer:

- simplicity over extensibility
- explicitness over abstraction
- maintainability over cleverness
- standard library over external packages
- synchronous flow over concurrency
