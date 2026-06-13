# Jewelry Store

A backend for a small handmade-jewelry shop. Most items are one of a kind, so
the data model is built around that idea: every product carries a "story", a
stone, a list of materials, and an `is_unique` flag that keeps stock at one.

I started this as a way to get comfortable with idiomatic Go service code —
clean separation between transport, business logic and storage — rather than to
ship the fanciest catalog on the internet. The product side is working
end to end; the rest of the schema (users, cart, orders) is designed but not
wired up yet.

## Status

What works today:

- Product catalog: create, read, update, delete, change status
- Public listing with filtering (category, stone, price range), sorting and pagination
- Postgres storage with migrations, config from the environment, graceful shutdown

Designed in the DB schema but not implemented yet: auth (users + refresh tokens),
cart, orders, product images. The `docker-compose` file already brings up Redis,
MinIO, Jaeger and Prometheus so those pieces have somewhere to plug into when I
get to them — none of them are touched by the running service right now.

## Tech stack

- Go 1.26, [Gin](https://github.com/gin-gonic/gin) for HTTP
- PostgreSQL 16 via [sqlx](https://github.com/jmoiron/sqlx)
- [Viper](https://github.com/spf13/viper) for config (12-factor, env vars)
- [golang-migrate](https://github.com/golang-migrate/migrate) for schema migrations
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

Set the database password (the only required secret — everything else has a
sensible default):

```sh
export DB_PASSWORD=1111
```

Run the API:

```sh
make run
```

It listens on `:8080` by default. Other knobs, all optional:

| Variable      | Default     |
|---------------|-------------|
| `SERVER_PORT` | `8080`      |
| `DB_HOST`     | `127.0.0.1` |
| `DB_PORT`     | `5432`      |
| `DB_USER`     | `IvanDev`   |
| `DB_NAME`     | `mydb`      |
| `DB_SSLMODE`  | `disable`   |
| `DB_PASSWORD` | *(required)*|

Tests and linting:

```sh
make test
make lint
```

## API

Everything lives under `/api/v1`. Public reads are open; writes sit under
`/admin` and will move behind a JWT admin check once auth lands.

| Method | Path                          | What it does                  |
|--------|-------------------------------|-------------------------------|
| GET    | `/products`                   | List products (filter/sort/paginate) |
| GET    | `/products/:id`               | Get one product by ID         |
| POST   | `/admin/products`             | Create a product              |
| PUT    | `/admin/products/:id`         | Update a product              |
| PATCH  | `/admin/products/:id/status`  | Change status only            |
| DELETE | `/admin/products/:id`         | Delete a product              |

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

Products start as `draft` and become visible in the public listing once their
status is `published`.

## What's next

- Auth: registration, login, JWT + refresh tokens, admin guard on `/admin/*`
- Cart and orders (price/title snapshots at checkout are already in the schema)
- Image uploads to MinIO
- Observability: wire up the Prometheus and Jaeger that compose already starts
