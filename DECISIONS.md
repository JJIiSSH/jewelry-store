## ADR-001: Docker Compose для локальной среды разработки
**Context:** Нужна локальная среда с PostgreSQL, Redis, MinIO и другими сервисами.
**Decision:** Используем Docker Compose — все сервисы запускаются одной командой `make up`.
**Tradeoff:** Нужен Docker на машине, зато не нужно устанавливать PostgreSQL, Redis и MinIO вручную.

