# sqlc — Type-Safe SQL Code Generation

## What it does

sqlc reads your SQL query files and generates idiomatic, fully type-safe Go code. You write plain SQL; sqlc produces the structs, interfaces, and functions needed to call those queries from Go — with zero reflection and no runtime query building.

The result is a DB layer that is:
- **type-checked at compile time** — wrong column names or parameter types are caught before the app runs
- **zero boilerplate** — no manual `rows.Scan(...)` or struct mapping to write
- **readable** — the generated code is straightforward Go, not a DSL

In this project, the relevant paths are:

| Path | Purpose |
|---|---|
| `internal/adapters/postgresql/sqlc/queries.sql` | SQL queries you write |
| `internal/adapters/postgresql/migrations/` | Schema files sqlc reads to validate queries |
| `internal/adapters/postgresql/sqlc/` | Generated Go code (do not edit manually) |
| `sqlc.yaml` | Configuration file |

---

## Configuration (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./internal/adapters/postgresql/sqlc/queries.sql"
    schema: "./internal/adapters/postgresql/migrations"
    gen:
      go:
        package: "repo"
        out: "./internal/adapters/postgresql/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
```

| Field | Meaning |
|---|---|
| `engine` | Target database dialect |
| `queries` | File containing annotated SQL queries |
| `schema` | Directory sqlc reads to understand the table structure |
| `package` | Go package name for generated files |
| `out` | Output directory for generated code |
| `sql_package` | Driver used in generated code (`pgx/v5` here) |
| `emit_json_tags` | Adds `json:"..."` tags to generated structs |
| `emit_interface` | Generates a `Querier` interface for mocking in tests |

---

## Writing queries

Queries are annotated SQL in `queries.sql`. The annotation tells sqlc the function name and return type:

```sql
-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: CreateProduct :one
INSERT INTO products (name, description, price_in_cents, quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
```

| Annotation | Generated return type |
|---|---|
| `:one` | Single struct, returns `(T, error)` |
| `:many` | Slice of structs, returns `([]T, error)` |
| `:exec` | No rows returned, returns `error` |
| `:execresult` | Returns `(sql.Result, error)` |

---

## Installation

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

---

## Commands

### Generate Go code from SQL

```bash
sqlc generate
```

Reads `sqlc.yaml`, validates queries against the schema, and writes the generated files to the configured output directory. Run this every time you add or modify a query.

### Validate queries without generating

```bash
sqlc vet
```

Checks that all queries are syntactically valid and consistent with the schema. Useful in CI to catch errors without writing files.

---

## What gets generated

Running `sqlc generate` produces four files in `internal/adapters/postgresql/sqlc/`:

| File | Contents |
|---|---|
| `models.go` | Go structs for each database table |
| `db.go` | `DBTX` interface and `Queries` struct that wraps the connection |
| `querier.go` | `Querier` interface listing all query methods (enabled by `emit_interface: true`) |
| `queries.sql.go` | The actual query implementations |

Example of generated code for `GetProduct`:

```go
func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
    row := q.db.QueryRow(ctx, getProduct, id)
    var i Product
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Description,
        &i.PriceInCents,
        &i.Quantity,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}
```

**Never edit the generated files.** They are overwritten on every `sqlc generate` run.

---

## Workflow

```
1. Update the schema via a goose migration
2. Add or modify a query in queries.sql
3. Run `sqlc generate`
4. Use the generated function in your service layer
```

The `Querier` interface (generated because `emit_interface: true`) lets you inject a mock in unit tests without hitting a real database.
