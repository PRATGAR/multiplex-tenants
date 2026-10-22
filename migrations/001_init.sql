CREATE TABLE tenants (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    slug           TEXT NOT NULL UNIQUE,
    plan           TEXT NOT NULL DEFAULT 'starter',
    invoice_footer TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     BIGINT NOT NULL REFERENCES tenants(id),
    email         TEXT NOT NULL,
    display_name  TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'member',
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, email)
);

CREATE TABLE projects (
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  BIGINT NOT NULL REFERENCES tenants(id),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE invitations (
    id          BIGSERIAL PRIMARY KEY,
    tenant_id   BIGINT NOT NULL REFERENCES tenants(id),
    email       TEXT NOT NULL,
    role        TEXT NOT NULL,
    token       TEXT NOT NULL UNIQUE,
    status      TEXT NOT NULL DEFAULT 'pending',
    accepted_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE domains (
    id           BIGSERIAL PRIMARY KEY,
    tenant_id    BIGINT NOT NULL REFERENCES tenants(id),
    name         TEXT NOT NULL,
    verify_token TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'pending'
);

CREATE TABLE exports (
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  BIGINT NOT NULL REFERENCES tenants(id),
    row_count  INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE billing_accounts (
    tenant_id     BIGINT PRIMARY KEY REFERENCES tenants(id),
    seats         INT NOT NULL DEFAULT 1,
    balance_cents BIGINT NOT NULL DEFAULT 0,
    currency      TEXT NOT NULL DEFAULT 'EUR'
);

CREATE TABLE plan_prices (
    code       TEXT PRIMARY KEY,
    unit_cents INT NOT NULL,
    seat_cents INT NOT NULL
);

CREATE TABLE audit_log (
    id        BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    actor     TEXT NOT NULL,
    action    TEXT NOT NULL,
    at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO plan_prices (code, unit_cents, seat_cents) VALUES
    ('starter', 0, 0),
    ('team', 4900, 500),
    ('business', 19900, 900),
    ('enterprise', 99900, 1500);
