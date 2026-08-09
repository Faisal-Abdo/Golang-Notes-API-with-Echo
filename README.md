# Go Notes API

A small CRUD Notes API built to practice a layered Go architecture with real
authentication/authorization, not just a bare `net/http` toy.

Built with:

- Go
- [Echo](https://echo.labstack.com/) (HTTP framework)
- PostgreSQL
- [sqlc](https://sqlc.dev/) (generated, type-safe DB queries)
- Docker / Docker Compose
- Keycloak (identity provider)
- JWT / OIDC

The project started as a `net/http` CRUD API and was later migrated to Echo.

## Architecture

```
Client
  |  HTTP + JWT
  v
Echo
  |
  v
Authentication Middleware   (verifies the JWT against Keycloak)
  |
  v
Authorization Middleware    (checks realm role, only where required)
  |
  v
Handler                     (HTTP concerns: binding, validation, status codes)
  |
  v
Service                     (business logic)
  |
  v
Repository (sqlc-generated) (SQL)
  |
  v
PostgreSQL
```

Business logic and data access are kept independent of the HTTP framework -
`internal/notes/service.go` and the generated `internal/database` package
don't know Echo exists.

## Authentication

```
Keycloak -> OIDC discovery -> JWT Verifier -> Echo middleware -> Claims
```

`internal/auth/middleware.go` implements `AuthMiddleware.Authenticate`:

1. Reads the `Authorization` header.
2. Requires a `Bearer <token>` value.
3. Verifies the JWT's signature and issuer against Keycloak's OIDC verifier.
4. Decodes the verified claims into the project's own `auth.Claims` type
   (keeping the rest of the app independent of the OIDC library's types).
5. Stores the claims in Echo's context: `c.Set("user", claims)`.

Missing or malformed tokens return `401 Unauthorized` before any network call
to Keycloak is made.

## Authorization

Two realm roles exist in Keycloak: `USER` and `ADMIN` (role names are
case-sensitive). `internal/auth/role_middleware.go` implements
`RequireRole(role)`, which reads the claims set by `Authenticate` and checks
`claims.RealmAccess.Roles`.

Policy for this project:

| Route              | Requires        |
|---------------------|-----------------|
| `GET /notes`         | any authenticated user |
| `GET /notes/:id`     | any authenticated user |
| `POST /notes`        | any authenticated user |
| `PUT /notes/:id`     | any authenticated user |
| `DELETE /notes/:id`  | `ADMIN` role |

`RequireRole` returns `401` if there are no claims in context (auth didn't
run / failed) and `403` if the user is authenticated but lacks the role.

## Running the project

### 1. Start Postgres and Keycloak

```
docker compose up -d postgres keycloak
```

This brings up:
- Postgres on `localhost:5432` (persisted in the `postgres_data` volume)
- Keycloak admin console on `http://localhost:8081` (persisted in the
  `keycloak_data` volume, admin/root)

In Keycloak you need a realm `notes-realm`, a client `notes-api` with Direct
Access Grants enabled, realm roles `USER` and `ADMIN`, and at least one user
with **email, first name, and last name filled in** (Keycloak's User Profile
requires those for login even though the admin console lets you create a user
without them).

### 2. Run the API

Set `.env` (copy from the example below) or export the variables directly:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/notesdb?sslmode=disable
KEYCLOAK_URL=http://localhost:8081
KEYCLOAK_REALM=notes-realm
KEYCLOAK_CLIENT_ID=notes-api
```

Then:

```
run.cmd
```

or directly:

```
go run ./cmd/api
```

The API listens on `:8080`.

> The `api` service defined in `docker-compose.yml` can also containerize the
> API itself, but Keycloak issues tokens with an `iss` claim tied to whatever
> host:port the client used to log in. If you obtain tokens via
> `http://localhost:8081` (the normal case on your machine), a containerized
> API talking to Keycloak via a different hostname will fail issuer
> validation. Running the API locally against Dockerized Postgres/Keycloak
> (as above) is the supported dev workflow; `docker compose up` for the full
> stack works if your clients also authenticate through the same hostname the
> API container uses to reach Keycloak (`host.docker.internal` by default).

### 3. Get a token and call the API

```
curl -X POST http://localhost:8081/realms/notes-realm/protocol/openid-connect/token \
  -d grant_type=password \
  -d client_id=notes-api \
  -d username=<your-user> \
  -d password=<your-password>
```

Use the returned `access_token` as a Bearer token against the API. A Postman
collection covering all CRUD requests is included under `postman/`.

## API endpoints

| Method | Path          | Auth              | Description        |
|--------|---------------|-------------------|---------------------|
| GET    | `/notes`      | Bearer token       | List all notes      |
| GET    | `/notes/:id`  | Bearer token       | Get a note by ID    |
| POST   | `/notes`      | Bearer token       | Create a note       |
| PUT    | `/notes/:id`  | Bearer token       | Update a note       |
| DELETE | `/notes/:id`  | Bearer token, `ADMIN` role | Delete a note |

Request/response body for create/update:

```json
{
  "title": "Groceries",
  "content": "Milk, eggs"
}
```

Both `title` and `content` are required (non-empty after trimming), or the
API returns `400`.

### Status codes

| Status | Meaning |
|--------|---------|
| 200 / 201 / 204 | Success |
| 400 | Invalid request body, parameters, or note ID |
| 401 | Missing/invalid/expired token |
| 403 | Authenticated but missing the required role |
| 404 | Note does not exist |
| 500 | Unexpected server/database failure (details are logged server-side, not returned to the client) |

## Tests

```
go test ./...
```

Covers the authentication header parsing (missing/malformed tokens),
`RequireRole` (401/403/200 cases), and note validation. Full CRUD behavior
against a live database is exercised via the Postman collection rather than
automated tests, in keeping with this project's toy-project scope.
