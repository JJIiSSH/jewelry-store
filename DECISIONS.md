## ADR-001: Docker Compose for Local Development Environment
**Context:** A local environment with PostgreSQL, Redis, MinIO and other services is needed.
**Decision:** We use Docker Compose — all services start with a single `make up` command.
**Tradeoff:** Requires Docker on the machine, but eliminates the need to install PostgreSQL, Redis, and MinIO manually.

## ADR-002: Validation Strategy

**Status:** Accepted

**Context:**
Without a unified error-handling strategy, each layer handles errors differently, leading to unexpected errors at the output boundary.

**Decision:**
1. Handler rejects completely malformed requests not worth validating against business logic, using go-playground/validator via struct tags
2. Service processes requests that have already passed primary validation and adds its own — checking for business logic errors.

**Consequences:**
Request validation must be described in both the Handler and Service layers.
