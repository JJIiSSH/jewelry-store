# Jewelry Store

A backend for a small handmade-jewelry shop. Most items are one of a kind, so
the data model is built around that idea: every product carries a "story", a
stone, a list of materials, and an `is_unique` flag that keeps stock at one.

I started this as a way to get comfortable with idiomatic Go service code —
clean separation between transport, business logic and storage — rather than to
ship the fanciest catalog on the internet. The product side works end to end and
user registration is live. JWT and refresh-token generation are implemented as
service helpers, but login, token rotation and protected routes are not wired yet.

## Status

What works today:

- Product catalog: create, read, update, delete, change status
- Public listing with filtering (category, stone, price range), sorting and pagination
- User registration (`POST /auth/register`) with bcrypt-hashed passwords
- Postgres storage with migrations, config from the environment, graceful shutdown
- Access-JWT generation with HS256 and opaque 256-bit refresh-token generation

Auth is still in progress: the token helpers are not called by a login flow, the
refresh-token repository is unfinished, and `/admin/*` has no JWT guard. Categories,
product images, cart and orders currently exist only as schema/domain groundwork;
there are no HTTP endpoints for them. The `docker-compose` file also starts Redis,
MinIO, Jaeger and Prometheus, but the running API does not use them yet.

The public product queries currently return every product status, including drafts.
A published-only public catalog is a pending contract fix.

## Tech stack

- Go 1.26, [Gin](https://github.com/gin-gonic/gin) for HTTP
- PostgreSQL 16 via [sqlx](https://github.com/jmoiron/sqlx)
- [Viper](https://github.com/spf13/viper) for config (12-factor, env vars)
- [golang-migrate](https://github.com/golang-migrate/migrate) for schema migrations
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) for access tokens
- `testify` for tests

## Layout

```
cmd/api/            entrypoint: wiring + graceful shutdown
internal/
  domain/           core types and repository interfaces, no external deps
  service/          business logic (slug generation, validation, stock rules)
  repository/postgres/  sqlx-backed storage
  handler/httphandler/  Gin handlers, DTOs, request mapping
  config/           env loading
migrations/         SQL up/down
deployments/        docker-compose for local dependencies
```

The dependency direction points inward: handlers depend on a `ProductService`
interface, the service depends on a repository interface from `domain`, and the
Postgres implementation sits at the edge. That makes the service testable
without a database.

## Running it locally

You'll need Go, Docker and the `migrate` CLI.

Start the dependencies (Postgres and friends):

```sh
make up
```

Apply the schema:

```sh
make migrate-up
```

Load the required secrets into the current shell. Local secrets are kept outside
the repository in `~/.Codex/secrets/api-keys.env`:

```sh
source ~/.Codex/secrets/api-keys.env
```

The file must export both `DB_PASSWORD` and `JWT_SECRET`.

Run the API:

```sh
make run
```

It listens on `:8080` by default.

| Variable          | Default      |
|-------------------|--------------|
| `SERVER_PORT`     | `8080`       |
| `DB_HOST`         | `127.0.0.1`  |
| `DB_PORT`         | `5432`       |
| `DB_USER`         | `IvanDev`    |
| `DB_NAME`         | `mydb`       |
| `DB_SSLMODE`      | `disable`    |
| `DB_PASSWORD`     | *(required)* |
| `JWT_SECRET`      | *(required)* |
| `JWT_ACCESS_TTL`  | `15m`        |
| `JWT_REFRESH_TTL` | `720h`       |

Tests and linting:

```sh
make test
make lint
```

At the current auth work-in-progress checkpoint, tests and `go vet` pass. The
linter reports the token helpers as unused until login is wired to them.

## API

Everything lives under `/api/v1`. The table below lists routes that are actually
registered today. `docs/api-contracts.md` also describes target endpoints that
have not been implemented yet.

> **Warning:** `/admin/*` is not authenticated yet. Run the API only in a trusted
> local development environment until the JWT admin guard is implemented.

| Method | Path                          | What it does                  |
|--------|-------------------------------|-------------------------------|
| POST   | `/auth/register`              | Register a customer account   |
| GET    | `/products`                   | List products (filter/sort/paginate) |
| GET    | `/products/:id`               | Get one product by ID         |
| POST   | `/admin/products`             | Create a product              |
| PUT    | `/admin/products/:id`         | Update a product              |
| PATCH  | `/admin/products/:id/status`  | Change status only            |
| DELETE | `/admin/products/:id`         | Delete a product              |

Registration currently returns only the created user ID:

```json
{ "data": { "id": "uuid" } }
```

Login, refresh and logout routes do not exist yet.

List query parameters: `category`, `stone`, `min_price`, `max_price`,
`sort` (`price_asc` \| `price_desc` \| `newest`), `page`, `limit`.

A quick look at the catalog:

```sh
curl 'http://localhost:8080/api/v1/products?stone=amethyst&sort=price_asc&page=1&limit=20'
```

Creating a product:

```sh
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Amethyst silver ring",
    "price": 4500,
    "description": "Hand-forged 925 silver band with a raw amethyst.",
    "story": "Made from a single stone found on a trip to the Urals.",
    "stone": "amethyst",
    "materials": ["silver 925"],
    "weight_g": 6.2,
    "size": "17",
    "is_unique": true,
    "stock": 1,
    "category_id": "00000000-0000-0000-0000-000000000000"
  }'
```

Products start as `draft`. The intended public contract is to expose only
`published` products, but that filter is not enforced yet.

## What's next

- Finish auth: login, refresh-token persistence/rotation, logout and auth tests
- Add JWT authentication and an admin-role guard to `/admin/*`
- Stabilize the frontend catalog contract: published-only reads, categories,
  product images and CORS or a Next.js API proxy
- Cart and orders (price/title snapshots at checkout are already in the schema)
- Image uploads to MinIO
- Observability: wire up the Prometheus and Jaeger that compose already starts
