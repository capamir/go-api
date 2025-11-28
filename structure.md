# Project Structure

```
.
├── .gitattributes
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── cmd/
│   ├── main.go
│   └── api/
│       └── api.go
├── configs/
│   └── env.go
├── db/
│   └── db.go
├── services/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   ├── cart/
│   │   ├── routes.go
│   │   └── service.go
│   ├── order/
│   │   └── store.go
│   ├── product/
│   │   ├── routes.go
│   │   └── store.go
│   └── user/
│       ├── routes_test.go
│       ├── routes.go
│       └── store.go
├── types/
│   ├── cart.go
│   ├── product.go
│   └── user.go
├── utils/
│   └── utils.go
└── cmd/migrate/
    ├── main.go
    └── migrations/
        ├── 20251125220243_add-user-table.down.sql
        └── 2025112520243_add-user-table.up.sql
```