# API Contracts

> The authoritative source of the API. Backend implements, Frontend consumes.
> Read by: Backend Core agent, Frontend agent.
>
> Base URL: `http://localhost:8080/api/v1`
> All protected endpoints require: `Authorization: Bearer <access_token>`

---

## Conventions

- All responses are wrapped in an envelope:
  ```json
  { "data": { ... } }           // success
  { "error": "message" }        // error
  ```
- Pagination via query params: `?page=1&limit=20`
- Paginated response:
  ```json
  { "data": [...], "total": 100, "page": 1, "limit": 20 }
  ```
- Dates — ISO 8601: `"2026-04-25T12:00:00Z"`

---

## Auth

### POST /auth/register
Register a new customer.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "name": "Anna",
  "phone": "+79001234567"
}
```

**Response 201:**
```json
{
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "user": { "id": "uuid", "email": "user@example.com", "name": "Anna", "role": "customer" }
  }
}
```

---

### POST /auth/login

**Request:**
```json
{ "email": "user@example.com", "password": "secret123" }
```

**Response 200:** same as `/register`

---

### POST /auth/refresh
Refresh access_token using refresh_token.

**Request:**
```json
{ "refresh_token": "eyJ..." }
```

**Response 200:**
```json
{ "data": { "access_token": "eyJ...", "refresh_token": "eyJ..." } }
```

---

### POST /auth/logout 🔒
Invalidate refresh_token.

**Request:**
```json
{ "refresh_token": "eyJ..." }
```

**Response 200:**
```json
{ "data": { "message": "logged out" } }
```

---

## Products (Public)

### GET /products
Catalog with filtering.

**Query params:**
| Parameter | Type | Example |
|-----------|------|---------|
| `category` | string | `rings` |
| `stone` | string | `amethyst` |
| `min_price` | number | `1000` |
| `max_price` | number | `9999` |
| `sort` | string | `price_asc`, `price_desc`, `newest` |
| `page` | int | `1` |
| `limit` | int | `20` |

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Amethyst Ring",
      "slug": "amethyst-ring-silver",
      "price": 3500.00,
      "stone": "amethyst",
      "status": "published",
      "is_unique": true,
      "primary_image": "https://minio.../image.jpg",
      "category": { "id": "uuid", "name": "Rings", "slug": "rings" }
    }
  ],
  "total": 42,
  "page": 1,
  "limit": 20
}
```

---

### GET /products/:id
Full product card.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "title": "Amethyst Ring",
    "slug": "amethyst-ring-silver",
    "description": "A delicate ring...",
    "story": "This stone was found...",
    "stone": "amethyst",
    "materials": ["925 silver", "gold plating"],
    "price": 3500.00,
    "weight_g": 4.2,
    "size": "17",
    "is_unique": true,
    "stock": 1,
    "status": "published",
    "category": { "id": "uuid", "name": "Rings", "slug": "rings" },
    "images": [
      { "id": "uuid", "url": "https://...", "alt": "front view", "is_primary": true, "order_index": 0 }
    ],
    "created_at": "2026-04-25T10:00:00Z"
  }
}
```

---

## Categories (Public)

### GET /categories
List all categories.

**Response 200:**
```json
{
  "data": [
    { "id": "uuid", "name": "Rings", "slug": "rings", "image_url": "https://..." }
  ]
}
```

---

## Cart 🔒

### GET /cart
Current user's cart.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "items": [
      {
        "id": "uuid",
        "product": { "id": "uuid", "title": "...", "price": 3500, "primary_image": "..." },
        "quantity": 1
      }
    ],
    "total": 3500.00
  }
}
```

---

### POST /cart/items
Add a product to the cart.

**Request:**
```json
{ "product_id": "uuid", "quantity": 1 }
```

**Response 200:** full cart (same as GET /cart)

---

### DELETE /cart/items/:product_id
Remove a product from the cart.

**Response 200:** full cart

---

### DELETE /cart
Clear the cart.

**Response 200:**
```json
{ "data": { "message": "cart cleared" } }
```

---

## Orders 🔒

### POST /orders
Place an order from the cart.

**Request:**
```json
{
  "shipping_address": "123 Main St, Belgrade",
  "contact_phone": "+38161234567",
  "comment": "Please use nice packaging"
}
```

**Response 201:**
```json
{
  "data": {
    "id": "uuid",
    "status": "pending",
    "total": 3500.00,
    "items": [
      { "product_id": "uuid", "title": "Amethyst Ring", "price": 3500.00, "quantity": 1 }
    ],
    "created_at": "2026-04-25T12:00:00Z"
  }
}
```

---

### GET /orders
Order history for the current user.

**Response 200:** paginated list of orders

---

### GET /orders/:id
Details of a specific order.

**Response 200:** full order object

---

## Admin: Products 🔒👑

### POST /admin/products

**Request:**
```json
{
  "category_id": "uuid",
  "title": "Amethyst Ring",
  "slug": "amethyst-ring-silver",
  "description": "...",
  "story": "...",
  "stone": "amethyst",
  "materials": ["925 silver"],
  "price": 3500.00,
  "weight_g": 4.2,
  "size": "17",
  "is_unique": true
}
```

**Response 201:** full product object

---

### PUT /admin/products/:id
Update a product. Body same as POST.

**Response 200:** updated product

---

### PATCH /admin/products/:id/status

**Request:**
```json
{ "status": "published" }
```

**Response 200:** updated product

---

### POST /admin/products/:id/images
Upload a photo. `multipart/form-data`.

**Form fields:**
- `image` — file (JPEG/PNG/WebP, max 10MB)
- `alt` — alt text (optional)
- `is_primary` — `true`/`false`

**Response 201:**
```json
{ "data": { "id": "uuid", "url": "https://minio/...", "order_index": 2 } }
```

---

### DELETE /admin/products/:id/images/:image_id

**Response 200:**
```json
{ "data": { "message": "image deleted" } }
```

---

### DELETE /admin/products/:id

**Response 200:**
```json
{ "data": { "message": "product deleted" } }
```

---

## Admin: Orders 🔒👑

### GET /admin/orders

**Query params:** `status`, `page`, `limit`

**Response 200:** paginated list of all orders with user data

---

### GET /admin/orders/:id

**Response 200:** full order + user data

---

### PATCH /admin/orders/:id/status

**Request:**
```json
{ "status": "confirmed" }
```

**Response 200:** updated order

---

## Admin: Categories 🔒👑

### POST /admin/categories

**Request:**
```json
{ "name": "Rings", "slug": "rings", "description": "..." }
```

**Response 201:** created category

---

### PUT /admin/categories/:id

**Response 200:** updated category

---

### DELETE /admin/categories/:id

**Response 200:**
```json
{ "data": { "message": "category deleted" } }
```

---

## Legend

- 🔒 — requires JWT (any authenticated user)
- 👑 — requires `admin` role

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request (validation) |
| 401 | Unauthorized |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (e.g. product already in cart) |
| 500 | Internal Server Error |
