Now that the cart and order system is implemented, you’re ready for full Phase 4 documentation! Below is a comprehensive `docs/phase4.md` for your project.

***

# Phase 4: Shopping Cart & Orders

## Overview

This phase delivers the core e-commerce user flow: persistent shopping carts per user, robust cart item management, and a full-featured order system including checkout, order history, and administrative oversight. All relationships, business rules, and API endpoints are defined for both customers and admins.

***

## Models

### Cart & CartItem

```
User (1) ─ (1) Cart
Cart (1) ─ (N) CartItem (N) ─ (1) Product
```
Each user has one cart. Each cart contains multiple cart items, each referencing a product and a quantity. CartItem stores a snapshot of product price at the time it was added.

### Order & OrderItem

```
User (1) ─ (N) Order
Order (1) ─ (N) OrderItem (N) ─ (1) Product
```
Each order contains a snapshot of the purchased products (price, SKU, name), so order history and analytics are robust even if catalog data changes.

***

## Database Tables (by GORM migration)

- `carts`
- `cart_items`
- `orders`
- `order_items`

***

## API Endpoints

### Cart (Authenticated)

| Method | Endpoint                    | Description                |
|--------|-----------------------------|----------------------------|
| GET    | /api/v1/cart                | Get user's cart            |
| POST   | /api/v1/cart/items          | Add product to cart        |
| PUT    | /api/v1/cart/items/:item_id | Update cart item quantity  |
| DELETE | /api/v1/cart/items/:item_id | Remove cart item           |
| DELETE | /api/v1/cart                | Clear cart                 |

### Orders (Authenticated)

| Method | Endpoint                      | Description                |
|--------|-------------------------------|----------------------------|
| POST   | /api/v1/orders                | Checkout (new order)       |
| GET    | /api/v1/orders                | User's order history       |
| GET    | /api/v1/orders/:id            | Order details              |
| PUT    | /api/v1/orders/:id/cancel     | Cancel order (if eligible) |

### Orders (Admin)

| Method | Endpoint                             | Description                    |
|--------|--------------------------------------|--------------------------------|
| GET    | /api/v1/admin/orders                 | All orders with filters        |
| GET    | /api/v1/admin/orders/stats           | Order statistics               |
| PUT    | /api/v1/admin/orders/:id/status      | Update order status            |

***

## Cart Logic

- Each authenticated user auto-gets a cart.
- Adding a product increases quantity, or adds new CartItem.
- Adding/updating checks for sufficient product stock.
- Removing item or clearing cart supported.
- Price in CartItem is snapshot at time of add.

***

## Order & Checkout Logic

- Checkout only possible if all cart items are in stock.
- Order creation auto-generates an order number (`ORD-YYYYMMDD-#####`).
- On checkout:
  - Copies cart items to Order, snapshotting price, SKU, product name.
  - Calculates subtotal, tax (default 10%), shipping (free for subtotal ≥ $100).
  - Decreases product inventory.
  - Clears user cart.
- User can view/track order history.
- Order can be cancelled by user if status is "pending" or "processing" (restores stock).
- Admin can update order status through `/admin/orders/:id/status`.

***

## Validation & Security

- All cart and order endpoints require JWT auth.
- Permissions checked for all write operations.
- Stock check on every add-update in cart or at checkout.
- Admin-only routes validated via user role.

***

## Extensible Business Logic

- Easy to attach custom shipping/tax/payments functions.
- Price/availability always checked at mutating actions.
- Order statistics support basic admin dashboard analytics.

***

## Example API Calls

**Add to Cart:**
```bash
curl -X POST /api/v1/cart/items -H "Authorization: Bearer <JWT>" \
  -d '{"product_id": 3, "quantity": 2}'
```

**Checkout:**
```bash
curl -X POST /api/v1/orders -H "Authorization: Bearer <JWT>" \
  -d '{"shipping_address": { ... }, "notes": "Leave at front desk"}'
```

**Get User Orders:**
```bash
curl -X GET /api/v1/orders -H "Authorization: Bearer <JWT>"
```

***

## What's Next

- Integrate payments (Stripe or gateway of choice)
- Support guest/anonymous carts (optional)
- Add order history UI, order tracking, fulfillment workflow

