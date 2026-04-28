CREATE TYPE product_status AS ENUM ('draft', 'published', 'sold', 'archived');
CREATE TYPE product_type   AS ENUM ('ring', 'necklace', 'earrings', 'bracelet', 'brooch', 'pendant');
CREATE TYPE user_role      AS ENUM ('customer', 'admin');
CREATE TYPE order_status   AS ENUM ('pending', 'confirmed', 'in_progress', 'shipped', 'delivered', 'cancelled');



CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,  -- 'rings', 'necklaces'
    description TEXT,
    image_url   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id),
    title       VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,  -- SEO-friendly URL
    description TEXT,
    story       TEXT,                          -- "история изделия" — ключевой атрибут для handmade
    stone       VARCHAR(100),                  -- основной камень: 'аметист', 'лазурит', ...
    materials   TEXT[]       NOT NULL DEFAULT '{}', -- ['серебро 925', 'позолота']
    price       DECIMAL(10,2) NOT NULL,
    weight_g    DECIMAL(8,2),                  -- вес в граммах
    size        VARCHAR(50),                   -- '17' для кольца, '18 см' для браслета
    is_unique   BOOLEAN NOT NULL DEFAULT TRUE, -- единственный экземпляр
    stock       INT     NOT NULL DEFAULT 1,    -- для is_unique=true всегда 1
    status      product_status NOT NULL DEFAULT 'draft',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT stock_unique_check CHECK (
        is_unique = FALSE OR stock <= 1        -- уникальное изделие не может быть в количестве > 1
    )
);

CREATE INDEX idx_products_category  ON products(category_id);
CREATE INDEX idx_products_status    ON products(status);
CREATE INDEX idx_products_stone     ON products(stone);
CREATE INDEX idx_products_price     ON products(price);

CREATE TABLE product_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,                 -- MinIO object URL
    alt         VARCHAR(255),
    order_index INT  NOT NULL DEFAULT 0,       -- порядок в галерее
    is_primary  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(product_id, order_index)
);

CREATE INDEX idx_product_images_product ON product_images(product_id);

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

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE, -- храним хэш, не сам токен
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user    ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);

CREATE TABLE carts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE, -- 1 корзина на пользователя
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cart_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id    UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity   INT  NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(cart_id, product_id)
);

CREATE INDEX idx_cart_items_cart ON cart_items(cart_id);

CREATE TABLE orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    status           order_status NOT NULL DEFAULT 'pending',
    total            DECIMAL(10,2) NOT NULL,
    comment          TEXT,                      -- пожелания клиента
    shipping_address TEXT,
    contact_phone    VARCHAR(50),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user   ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);

CREATE TABLE order_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id   UUID         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID         NOT NULL REFERENCES products(id),
    title      VARCHAR(255) NOT NULL,     -- снэпшот названия на момент заказа
    price      DECIMAL(10,2) NOT NULL,    -- снэпшот цены на момент заказа
    quantity   INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_order_items_order ON order_items(order_id);