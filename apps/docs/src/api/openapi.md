# OpenAPI

The API uses [swag](https://github.com/swaggo/swag) to generate an OpenAPI 3.1 spec from annotations written directly in Go source files. The generated output lives in `apps/api/docs/` and is served at `/docs` by the running API server via [Scalar](https://scalar.com).

A hosted version of the spec is also available at [core.apidocumentation.com](https://core.apidocumentation.com/guide/swamphacks-core-api).

---

## Installation

The project uses a custom fork of swag on its `v2` branch. Install it from source:

```bash
git clone -b v2 https://github.com/hieunguyent12/swag.git
cd swag
go install ./cmd/swag
```

> Go must be installed. See [Installation & Setup](installation.md) for instructions.

Verify:

```bash
swag --version
```

---

## Generating the spec

From `apps/api/`:

```bash
make openapi-generate
```

This runs:

```
swag init --dir cmd/api,internal/api/handlers --parseDependency --requiredByDefault -v3.1
```

swag scans `cmd/api/main.go` (for top-level API metadata) and all handler files in `internal/api/handlers/` for route annotations. It writes three output files:

| File | Description |
|---|---|
| `docs/swagger.yaml` | OpenAPI 3.1 spec (YAML) |
| `docs/swagger.json` | OpenAPI 3.1 spec (JSON) |
| `docs/docs.go` | Embedded Go file for serving the spec at runtime |

Run this any time you add or change handler annotations.

## Formatting annotations

swag can also normalise the annotation comments in-place:

```bash
make openapi-format
```

This runs `swag fmt` over the handler files. Run it before committing annotation changes to keep formatting consistent.

---

## How annotations work

### API-level metadata

Top-level metadata is declared in `cmd/api/main.go` above the `main` function:

```go
//	@title		SwampHacks Test API
//	@version	1.0
//	@description	This is SwampHacks' OpenAPI documentation.
```

### Route annotations

Each exported handler method that corresponds to a route gets a block of annotations directly above its signature. swag reads these to build the spec.

**Example** (`internal/api/handlers/auth.go`):

```go
// GetMe
//
//	@Summary		Get Current User
//	@Description	Get the currently authenticated user's information.
//	@Tags			Authentication
//	@Produce		json
//	@Param			sh_session	cookie		string	true	"The authenticated session token/id"
//	@Success		200			{object}	middleware.UserContext
//	@Failure		401			{object}	response.ErrorResponse	"Unauthenticated"
//	@Failure		500			{object}	response.ErrorResponse
//	@Router			/auth/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
```

**Example** (`internal/api/handlers/events.go`):

```go
// Create a new event
//
//	@Summary		Create a new event
//	@Description	Create a new event with the provided details
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateEventFields		true	"Event creation data"
//	@Success		201		{object}	sqlc.Event				"Event created"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Router			/events [post]
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
```

### Common annotation fields

| Annotation | Description |
|---|---|
| `@Summary` | Short one-line description shown in the spec |
| `@Description` | Longer description |
| `@Tags` | Groups the route under a tag in the UI |
| `@Accept` | Request content type (e.g. `json`) |
| `@Produce` | Response content type (e.g. `json`) |
| `@Param` | Parameter definition — see format below |
| `@Success` | Success response with status code and body type |
| `@Failure` | Error response with status code and body type |
| `@Router` | Path and HTTP method — **required** |

### `@Param` format

```
@Param  <name>  <in>  <type>  <required>  "<description>"
```

`<in>` is one of: `query`, `path`, `body`, `header`, `cookie`.

```go
@Param  eventId  path    string           true   "Event UUID"
@Param  request  body    CreateEventFields true   "Request body"
@Param  limit    query   int              false  "Page size"
```

### Response body types

Pass a Go struct to `{object}` and swag will reflect its fields into the spec. Types from other packages work as long as `--parseDependency` is set (it is, via `make openapi-generate`):

```go
@Success  200  {object}  sqlc.Event
@Success  201  {object}  middleware.UserContext
@Failure  400  {object}  response.ErrorResponse
```
