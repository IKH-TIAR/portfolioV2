-- 0001: initial schema — users, content, comments, messages

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE projects (
    id          BIGSERIAL PRIMARY KEY,
    slug        TEXT NOT NULL UNIQUE,
    title       TEXT NOT NULL,
    tagline     TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    tech        TEXT[] NOT NULL DEFAULT '{}',
    repo_url    TEXT NOT NULL DEFAULT '',
    demo_url    TEXT NOT NULL DEFAULT '',
    image_url   TEXT NOT NULL DEFAULT '',
    featured    BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    published   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE posts (
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT NOT NULL UNIQUE,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    post_type  TEXT NOT NULL DEFAULT 'journal' CHECK (post_type IN ('journal', 'note')),
    published  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE photos (
    id         BIGSERIAL PRIMARY KEY,
    url        TEXT NOT NULL,
    caption    TEXT NOT NULL DEFAULT '',
    location   TEXT NOT NULL DEFAULT '',
    taken_at   DATE,
    published  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE site_settings (
    key   TEXT PRIMARY KEY,
    value JSONB NOT NULL
);

CREATE TABLE comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author     TEXT NOT NULL,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX comments_post_id_idx ON comments (post_id);

CREATE TABLE messages (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
