# Goose — Database Migrations

## What it does

Goose is a database migration tool. It tracks which SQL migrations have been applied to a database and lets you move the schema forward (`up`) or backward (`down`) in a controlled, versioned way.

Each migration is a plain `.sql` file prefixed with a sequential number (e.g. `00001_create_products.sql`). Goose uses that number to determine order and to record which migrations have already run — so running `goose up` twice is safe.

In this project, migration files live in:

```
internal/adapters/postgresql/migrations/
```

---

## Migration file format

Every file contains two annotated sections:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price_in_cents INTEGER NOT NULL CHECK (price_in_cents >= 0),
    ...
);

-- +goose Down
DROP TABLE IF EXISTS products;
```

| Section | When it runs |
|---|---|
| `+goose Up` | When applying the migration |
| `+goose Down` | When rolling back the migration |

---

## Installation

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Commands

### Apply all pending migrations

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "$DATABASE_URL" up
```

Runs every migration that has not been applied yet, in ascending order.

### Roll back the last migration

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "$DATABASE_URL" down
```

Executes the `-- +goose Down` block of the most recently applied migration.

### Check current status

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "$DATABASE_URL" status
```

Lists every migration file with its applied/pending state and timestamp.

### Apply up to a specific version

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "$DATABASE_URL" up-to 2
```

Stops after applying migration number `2`, regardless of how many files exist.

### Roll back to a specific version

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "$DATABASE_URL" down-to 1
```

Rolls back until version `1` is the latest applied migration.

### Create a new migration file

```bash
goose -dir ./internal/adapters/postgresql/migrations create add_users_table sql
```

Generates a new numbered file with the correct `+goose Up` / `+goose Down` stubs.

---

## How goose tracks state

Goose creates a `goose_db_version` table in your database on first run. It stores the version number of each applied migration. This is what makes `up` idempotent — it only runs files whose version number is not already in that table.

---

## Workflow

```
1. Write a new migration file (or use `goose create`)
2. Run `goose up` to apply it
3. Verify with `goose status`
4. If something is wrong: `goose down` to roll back, fix, then `goose up` again
```

When adding a new table or column, **always write the corresponding `Down` block** — it is the only safe way to undo the change in development or staging.
