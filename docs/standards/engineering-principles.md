# Purpose

This document defines engineering rules that guide all implementation decisions made by AI agents.

The goal is to ensure that the system is:

- simple
- testable
- predictable
- maintainable
- safe to run without supervision

The AI is expected to follow these rules even when no explicit instruction is given in a task.

## Core Engineering Mindset

The system must always prioritize:

- Correctness over cleverness
- Simplicity over abstraction
- Explicit behavior over implicit behavior
- Maintainability over extensibility
- Predictability over optimization
- Testability over convenience

If a solution feels "more advanced", it is usually incorrect unless explicitly justified.

## AI Autonomy Rules

AI is allowed to:

- design internal implementation freely
- refactor code if needed
- introduce helper functions
- optimize for clarity and testability

AI is NOT allowed to:

- introducing unnecessary architectural layers
- adding dependencies without strong justification
- assuming future requirements
- over-engineering abstractions
- creating frameworks inside the project

If unsure, AI must choose the simplest solution.

## Code Generation Principles

All generated code must be:

- readable without context switching
- explicit in behavior
- minimal in abstraction layers
- deterministic in execution
- structured around clear responsibilities

Hidden behavior is forbidden.

Magic behavior without explicit flow is forbidden.

## Function Design Rules

Functions should:

- do one thing clearly
- have predictable inputs and outputs
- avoid side effects unless explicitly required
- prefer early returns over nested logic
- remain small enough to be understood in a single read

If a function becomes hard to reason about, it must be split.

## Error Handling Standard

Errors must:

- clearly describe what failed
- preserve contextual information
- propagate upward unless handled explicitly
- never be silently ignored

Errors are part of system design, not edge cases.

User-facing errors must be understandable and actionable.

Internal errors must be detailed for debugging.

## Testing Requirements

All critical logic must be testable without external dependencies.

Tests must:

- be deterministic
- not rely on real HTTP calls
- not rely on real database state
- use mocks for external systems
- validate behavior, not implementation details

If code is hard to test, it is considered poorly designed.

## HTTP Communication Rules

All external HTTP communication must:

- be encapsulated in the client layer
- have clear timeout handling
- be testable via mocks
- avoid hidden retries unless explicitly defined
- never leak credentials in logs or errors

HTTP logic must not exist inside business logic layers.

## Database Interaction Rules

Database access must:

- be isolated in the repository layer
- use explicit queries or simple abstractions
- avoid ORM complexity
- allow easy replacement of SQLite with another backend
- be testable via mocks or in-memory substitutes

Business logic must never directly access SQL.

## Logging Principles

Logging must:

- support debugging without code inspection
- be structured and consistent
- never include sensitive data
- describe system behavior, not implementation noise

Log levels:

- trace: internal execution flow
- info: user-visible actions
- debug: technical state
- warn: recoverable issues
- error: failures requiring attention

Logging must not affect system behavior.

## Concurrency Rules

Concurrency is allowed only if:

- it simplifies execution flow
- it is fully deterministic
- it does not introduce hidden state

Hidden goroutines or background processes without explicit control are forbidden.

If concurrency increases complexity, it must be avoided.

## Dependency Rules

Dependencies must be:

- minimal
- justified
- replaceable

Standard library is preferred in all cases.

No dependency should be added for convenience alone.

Each dependency must solve a clear, unavoidable problem.

## Refactoring Rules

AI is allowed to refactor code when:

- it improves clarity
- it reduces duplication without adding abstraction layers
- it simplifies logic
- it improves testability

Refactoring must never introduce speculative architecture.

## Security Rules

The system must always:

- avoid logging credentials
- read secrets only from environment variables
- avoid persistent storage of sensitive data when not required
- sanitize logs before output

Security defaults must be safe by design.

## Performance Philosophy

Performance optimization is secondary.

Rules:

- Ensure correctness first
- Ensure readability second
- Optimize only when there is a measured need

Premature optimization is considered a design flaw.

## Anti-Patterns (Strictly Forbidden)

The following patterns must not be introduced:

- unnecessary abstraction layers
- global mutable state
- hidden side effects
- dependency injection frameworks
- reflection without strong reason
- panic-driven control flow
- overly generic utility packages
- speculative future-proofing


## Decision Rule Hierarchy

When multiple solutions exist, choose in this order:

1. simplest working solution
2. most explicit solution
3. most testable solution
4. most maintainable solution
5. most performant solution (only if needed)

## Final Rule

If a decision is unclear or ambiguous:

Always choose the simplest implementation that correctly solves the current requirement.
