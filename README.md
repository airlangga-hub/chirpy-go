# Chirpy HTTP Server

This is a backend HTTP Server for a social media app inspired by Twitter called **Chirpy**, where users can post short messages called **chirps**. This documentation outlines the HTTP API endpoints, environment setup, and core functionality of the `chirpy` server implemented in Go.

---

## 🚀 Overview

- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **ORM/Query Tool**: [sqlc](https://sqlc.dev/)
- **Migrations**: [Goose](https://github.com/pressly/goose)
- **Authentication**: JWT (JSON Web Tokens)
- **File Serving**: Static assets served under `/app`

---

## 🛠️ Environment Setup

The server requires the following environment variables to run:

| Variable       | Description                                   | Required |
|----------------|-----------------------------------------------|----------|
| `JWT_SECRET`   | Secret key used to sign/verify JWT tokens     | ✅ Yes   |
| `POLKA_KEY`    | API key for webhook integration (made up webhook just for practice)     | ✅ Yes   |
| `DB_URL`       | PostgreSQL connection string (e.g., `postgres://user:pass@localhost:5432/chirpy?sslmode=disable`) | ✅ Yes |
| `PLATFORM`     | Deployment platform identifier (e.g., `dev`, `prod`) | ✅ Yes |

> 💡 Load these via a `.env` file at the project root. The app uses [`joho/godotenv`](https://github.com/joho/godotenv) to read it automatically.

---

## 🗂️ Project Structure
```bash
.
├── assets
│   └── logo.png
├── go.mod
├── go.sum
├── handler_chirp_delete.go
├── handler_chirps_create.go
├── handler_chirps_get.go
├── handler_login.go
├── handler_polka_webhook.go
├── handler_refresh.go
├── handler_users_create.go
├── handler_users_update.go
├── index.html
├── internal
│   ├── auth
│   └── database
├── json.go
├── main.go
├── metrics.go
├── readiness.go
├── README.md
├── reset.go
├── sql
│   ├── queries
│   └── schema
└── sqlc.yaml
```

---

## 🌐 API Endpoints

### 🔌 Health & Metrics

| Method | Path               | Description                          |
|--------|--------------------|--------------------------------------|
| `GET`  | `/api/healthz`     | Returns `OK` if server is running    |
| `GET`  | `/admin/metrics`   | Returns number of file server hits   |
| `POST` | `/admin/reset`     | Resets in-memory metrics (dev only)  |

---

### 👤 User Management

| Method | Path             | Description                              |
|--------|------------------|------------------------------------------|
| `POST` | `/api/users`     | Create a new user                        |
| `PUT`  | `/api/users`     | Update user (email or password)          |
| `POST` | `/api/login`     | Authenticate and return access/refresh tokens |
| `POST` | `/api/refresh`   | Exchange refresh token for new access token |
| `POST` | `/api/revoke`    | Invalidate a refresh token               |

> 🔐 All user-related endpoints (except `/api/users` POST) require valid authentication.

---

### 🐦 Chirp Management

| Method | Path                     | Description                          |
|--------|--------------------------|--------------------------------------|
| `POST` | `/api/chirps`            | Create a new chirp                   |
| `GET`  | `/api/chirps`            | Retrieve all chirps (supports `author_id` and `sort` query params) |
| `GET`  | `/api/chirps/{chirpID}`  | Get a single chirp by ID             |
| `DELETE`| `/api/chirps/{chirpID}`  | Delete a chirp |

> 🔒 `POST` and `DELETE` require authentication.
> 📏 Chirps are limited to **140 characters**; longer content will be rejected.

---

### 💸 Webhook (to Upgrade Users to Premium)

| Method | Path                    | Description                         |
|--------|-------------------------|-------------------------------------|
| `POST` | `/api/polka/webhooks`   | Handle upgrade to paid membership |

> 🔐 This endpoint validates the `Authorization: ApiKey {POLKA_KEY}` header.

---

## 📁 Static File Serving

All files in the project root (e.g., `index.html`, `assets/`) are served under `GET /app/{path}`
- Example: `http://localhost:8080/app/index.html` serves `./index.html`
- Requests to `/app/` are handled by Go’s `http.FileServer`

---

## 🔐 Authentication

Uses **Bearer JWT** tokens in the `Authorization` header
- Access tokens expire after **1 hour**.
- Refresh tokens are long-lived and stored securely in the database.
- Passwords are **hashed** using `bcrypt` before storage.

---

## 🗃️ Database Schema (via Goose Migrations)

Migrations in `sql/schema/`:

1. `001_users.sql` – Create `users` table (`id`, `email`, `created_at`, `updated_at`)
2. `002_chirps.sql` – Create `chirps` table (`id`, `body`, `user_id`, `created_at`)
3. `003_password.sql` – Add `password_hash` to `users`
4. `004_refresh_tokens.sql` – Store refresh tokens (`token_hash`, `user_id`, `expires_at`)
5. `005_is_chirpy_red.sql` – Add `is_chirpy_red` boolean to `users`

All queries are type-safe and generated by **sqlc** from files in `sql/queries/`.

---

## ▶️ Running the Server

1. Ensure PostgreSQL is running and `DB_URL` points to a valid database.
2. Run migrations:
 ```bash
 goose -dir sql/schema postgres "$DB_URL" up
 ```
 3. Start the server:
 ```bash
 go run .
 ```
 4. Server runs on `http://localhost:8080`

 ---

 ## 🧪 Testing

- Unit tests exist for auth logic (`internal/auth/auth_test.go`)
- Manual testing recommended via curl, Postman, or frontend integration