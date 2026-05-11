# Vision — Sokker.org Auto Bidder

## Project Name

Sokker.org Auto Bidder

---

## Product Vision

Sokker.org Auto Bidder is a lightweight automation tool for players of the online game Sokker.org.

The application automates participation in player auctions by placing bids on behalf of the user using the lowest possible valid bid amount, up to a configured maximum price.

The goal of the application is to save user time and eliminate the need to manually monitor auctions close to their end.

---

## Problem Statement

In Sokker.org, users must currently stay online and manually observe auctions near their ending time.

The auction ending time may change dynamically due to additional bids, which forces players to:
- repeatedly check auction status,
- remain available for long periods of time,
- manually react to competing bids,
- spend significant time on repetitive actions.

This creates unnecessary friction and wastes user time.

The application solves this problem by automating the bidding process.

---

## Target Users

### Primary User

Players of Sokker.org who:
- actively participate in transfer auctions,
- do not want to manually monitor auctions,
- want to save time,
- prefer simple automation over manual interaction.

### Technical Profile

Target users are expected to:
- have basic technical knowledge,
- be capable of running a CLI/server application,
- configure environment variables and cron jobs if needed.

---

## Core Use Cases

### Automated Auction Participation

User provides:
- Sokker.org credentials,
- player identifiers,
- maximum acceptable bid price.

The application:
- authenticates to Sokker.org,
- monitors selected auctions,
- places bids automatically,
- always bids the lowest possible valid amount,
- never exceeds configured maximum price.

---

## Product Superpowers

The application prioritizes:

- Simplicity
- Automation
- Security
- Low operational cost

---

## MVP Scope

The MVP version must:

- support a single user only,
- run as a CLI/server process,
- work on PC or server environments,
- support automated bidding,
- use cron-based scheduling,
- persist required data locally.

---

## Explicitly Out of Scope (MVP)

The MVP will NOT include:

- multi-user support,
- graphical user interface,
- automatic price suggestions,
- AI-assisted bidding,
- cloud infrastructure,
- distributed architecture,
- mobile support.

---

## Product Model

- Open-source project
- Single-user architecture
- Self-hosted
- Local execution on PC/server
- Monolithic application

---

## AI Usage

At the current stage, AI is NOT part of the runtime product flow.

AI may be used only during software development and code generation processes.

---

## Technical Priorities

Priority order:

1. Fast development
2. Simplicity
3. Low cost
4. Monolithic architecture
5. Easy maintenance

Architecture and implementation decisions should always favor:
- lower complexity,
- easier debugging,
- minimal infrastructure requirements,
- reduced operational overhead.

---

## Technology Decisions

### Language

- Golang

### Database

- SQLite

### Logging

- sirupsen/logrus

### Testing

- stretchr/testify

### Mocking

- jarcoal/httpmock
- DATA-DOG/go-sqlmock

### Configuration

- credentials provided through environment variables,
- scheduling controlled externally via crontab.

---

## Security Assumptions

The application handles user credentials for Sokker.org accounts.

Security expectations:
- credentials must never be logged,
- secrets should be read only from environment variables,
- application should minimize credential exposure,
- local-first architecture preferred over external services.

---

## Architectural Principles

The project should prefer:
- simple code over abstraction,
- explicit logic over magic,
- low dependency count,
- deterministic behavior,
- standard library solutions where reasonable,
- maintainability over premature optimization.

The project should avoid:
- unnecessary interfaces,
- overengineering,
- microservices,
- event-driven complexity,
- introducing infrastructure without clear need.

---

## Success Criteria

The project is successful if:

- users no longer need to manually monitor auctions,
- bidding works reliably,
- setup is simple,
- operational cost is near zero,
- the application can run unattended on a local machine or server.
