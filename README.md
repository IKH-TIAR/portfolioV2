# Portfolio — personal site & brand home

Minimal + warm portfolio for an AI developer & automation engineer, with a journal (learning logs + photography), projects showcase, and a full content admin.

## Architecture

| Layer     | Tech            | Hosting        |
| --------- | --------------- | -------------- |
| Frontend  | Next.js + Tailwind | Vercel (free) |
| Backend   | Go (stdlib `net/http` only) | Render (free) |
| Database  | PostgreSQL       | Neon (free)    |
| Storage   | Supabase Storage | Supabase (free) |
| Email     | Resend           | Resend (free)  |

```
Vercel                 Render                 Neon                Supabase
┌──────────────┐       ┌──────────────┐       ┌──────────┐       ┌──────────────┐
│ Next.js      │◄─────►│ Go API       │◄─────►│ Postgres │       │ Storage      │
│ App Router   │ REST  │ auth/CRUD/   │       │ content  │◄──────│ (photos)     │
│ + Tailwind   │ +JWT  │ uploads/     │       │          │       │              │
│ /admin panel │       │ comments/mail│       │          │       │              │
└──────────────┘       └──────────────┘       └──────────┘       └──────────────┘
```

## Repos layout

```
backend/    Go API — stdlib net/http, pgx, JWT auth, embed-based migrations
frontend/   Next.js app (App Router, TypeScript, Tailwind)
```

## Development

### Backend

Requires Go 1.26+ and a running PostgreSQL.

```sh
cp backend/.env.example backend/.env   # fill in values
cd backend
go run ./cmd/migrate                   # apply migrations
go run ./cmd/api                       # start API on :8080
```

Key endpoints:

- `GET  /api/health`       — liveness + DB ping
- `POST /api/auth/login`   — admin login → JWT
- `POST /api/auth/logout`  — invalidates session
- `GET  /api/auth/me`      — current admin (auth-guarded)

### Frontend

Requires Node 22+.

```sh
cp frontend/.env.example frontend/.env.local   # fill in API URL
cd frontend
npm install
npm run dev
```

## Configuration (backend env vars)

| Variable          | Required | Description                              |
| ----------------- | -------- | ---------------------------------------- |
| `PORT`            | no       | HTTP port (default `8080`)               |
| `DATABASE_URL`    | yes      | PostgreSQL connection string             |
| `JWT_SECRET`      | yes      | Secret used to sign auth tokens          |
| `JWT_TTL`         | no       | Token lifetime (default `24h`)           |
| `ADMIN_EMAIL`     | yes      | Admin account email (seeded on boot)     |
| `ADMIN_PASSWORD`  | yes      | Admin account password (seeded on boot)  |
| `ALLOWED_ORIGINS` | no       | Comma-separated CORS origins             |
