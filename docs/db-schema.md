# Database Schema

> The authoritative source of the DB schema. All changes are made here BEFORE writing migrations.
> Read by: Backend Core agent, gRPC agent.

---

## Enums

```sql
CREATE TYPE product_status AS ENUM ('draft', 'published', 'sold', 'archived');
CREATE TYPE product_type   AS ENUM ('ring', 'necklace', 'earrings', 'bracelet', 'brooch', 'pendant');
CREATE TYPE user_role      AS ENUM ('customer', 'admin');
CREATE TYPE order_status   AS ENUM ('pending', 'confirmed', 'in_progress', 'shipped', 'delivered', 'cancelled');
```

---

## Tables

### categories

```sql
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,  -- 'rings', 'necklaces'
    description TEXT,
    image_url   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### products

```sql
CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id),
    title       VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,  -- SEO-friendly URL
    description TEXT,
    story       TEXT,                          -- "item story" — key attribute for handmade products
    stone       VARCHAR(100),                  -- primary stone: 'amethyst', 'lazurite', ...
    materials   TEXT[]       NOT NULL DEFAULT '{}', -- ['925 silver', 'gold plating']
    price       DECIMAL(10,2) NOT NULL,
    weight_g    DECIMAL(8,2),                  -- weight in grams
    size        VARCHAR(50),                   -- '17' for a ring, '18 cm' for a bracelet
    is_unique   BOOLEAN NOT NULL DEFAULT TRUE, -- one-of-a-kind item
    stock       INT     NOT NULL DEFAULT 1,    -- always 1 when is_unique=true
    status      product_status NOT NULL DEFAULT 'draft',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT stock_unique_check CHECK (
        is_unique = FALSE OR stock <= 1        -- unique item cannot have quantity > 1
    )
);

CREATE INDEX idx_products_category  ON products(category_id);
CREATE INDEX idx_products_status    ON products(status);
CREATE INDEX idx_products_stone     ON products(stone);
CREATE INDEX idx_products_price     ON products(price);
```

---

### product_images

```sql
CREATE TABLE product_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,                 -- MinIO object URL
    alt         VARCHAR(255),
    order_index INT  NOT NULL DEFAULT 0,       -- position in gallery
    is_primary  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(product_id, order_index)
);

CREATE INDEX idx_product_images_product ON product_images(product_id);
```

---

### users

```sql
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(255),
    phone         VARCHAR(50),
    role          user_role NOT NULL DEFAULT 'customer',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

---

### refresh_tokens

```sql
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE, -- we store the hash, not the raw token
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user    ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

---

### carts

```sql
CREATE TABLE carts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE, -- 1 cart per user
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### cart_items

```sql
CREATE TABLE cart_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id    UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity   INT  NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(cart_id, product_id)
);

CREATE INDEX idx_cart_items_cart ON cart_items(cart_id);
```

---

### orders

```sql
CREATE TABLE orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    status           order_status NOT NULL DEFAULT 'pending',
    total            DECIMAL(10,2) NOT NULL,
    comment          TEXT,                      -- customer notes
    shipping_address TEXT,
    contact_phone    VARCHAR(50),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user   ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
```

---

### order_items

```sql
CREATE TABLE order_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id   UUID         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID         NOT NULL REFERENCES products(id),
    title      VARCHAR(255) NOT NULL,     -- snapshot of product name at order time
    price      DECIMAL(10,2) NOT NULL,    -- snapshot of price at order time
    quantity   INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_order_items_order ON order_items(order_id);
```

---

## ER Diagram (text)

```
categories ──< products ──< product_images
                  │
users ──< carts ──< cart_items ──> products
  │
  └──< orders ──< order_items ──> products
  │
  └──< refresh_tokens
```

---

## Key Design Decisions

| Decision | Why |
|----------|-----|
| `is_unique + stock CHECK` | Uniqueness business rule is enforced at the DB level, not only in code |
| Price snapshot in `order_items` | Price is fixed at order time — order history is not broken when a product is edited |
| `slug` in products and categories | SEO-friendly URLs (`/catalog/ring-amethyst-silver`) without JOIN by id |
| `token_hash` in refresh_tokens | Raw tokens are never stored — only the hash (SHA-256) |
| `TIMESTAMPTZ` everywhere | Timestamps always include timezone — avoids problems when changing servers |
| `UUID` PK everywhere | Safer than serial (id is not guessable), easier horizontal scaling |
