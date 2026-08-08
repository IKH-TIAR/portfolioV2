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

Public:

- `GET  /api/health`            — liveness + DB ping
- `POST /api/auth/login`        — admin login → JWT
- `POST /api/auth/logout`       — invalidates session (auth)
- `GET  /api/auth/me`           — current admin (auth)
- `GET  /api/site`              — site settings (hero, about, ...)
- `GET  /api/projects`          — published projects
- `GET  /api/projects/{slug}`   — single published project
- `GET  /api/posts`             — published posts (`?limit=&offset=`)
- `GET  /api/posts/{slug}`      — single published post
- `GET  /api/posts/{slug}/comments` — comments on a post
- `POST /api/posts/{slug}/comments` — post a comment
- `GET  /api/photos`            — published photos
- `POST /api/contact`           — contact form (stores message, emails you if Resend configured)

Admin (all auth-guarded):

- `GET/POST    /api/admin/projects`         — list/create
- `GET/PUT/DELETE /api/admin/projects/{id}` — read/update/delete
- `GET/POST    /api/admin/posts`            — list/create
- `GET/PUT/DELETE /api/admin/posts/{id}`    — read/update/delete
- `GET/POST    /api/admin/photos`           — list/create
- `GET/PUT/DELETE /api/admin/photos/{id}`   — read/update/delete
- `GET         /api/admin/site-settings`    — list all settings
- `PUT         /api/admin/site-settings/{key}` — upsert one setting
- `DELETE      /api/admin/comments/{id}`    — remove a comment
- `GET/DELETE  /api/admin/messages`         — read/delete contact messages
- `POST        /api/admin/uploads`          — multipart image upload → Supabase

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
| `SUPABASE_URL`    | no       | Enables image uploads (Supabase Storage) |
| `SUPABASE_SERVICE_ROLE_KEY` | no | Service role key for uploads        |
| `SUPABASE_BUCKET` | no       | Storage bucket name                      |
| `RESEND_API_KEY`  | no       | Enables contact-form emails (Resend)     |
| `CONTACT_FROM_EMAIL` | no    | Sender for contact emails                |
