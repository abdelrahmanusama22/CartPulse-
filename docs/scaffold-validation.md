# Initial scaffold validation

Validated on 10 November 2026 with Go 1.27.0 on Windows. The module declares Go 1.25.0 as its minimum version; the minimum toolchain itself was not separately run here.

Passed:

- gofmt on all Go source files.
- `go test ./...` (configuration and HTTP contract/recovery tests).
- `go vet ./...`.
- `go build ./...` for the native Windows environment.
- `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...` for the container target.
- `go mod verify`.
- YAML parsing of Compose, sqlc, OpenAPI, CI and the GitHub issue template.
- JSON parsing of the shared VS Code configuration files.
- A live temporary API process: liveness 200, readiness 503 with no database configured, and an unimplemented business route 404 with the expected JSON error/request ID.

Not run here:

- Docker/Compose execution or an actual PostgreSQL readiness success (Docker is unavailable in this environment).
- Race detector locally (a C compiler is unavailable); Linux GitHub CI is configured to run it.
- Business/database integration tests, sqlc generation and migrations: those implementations do not exist yet.

See README.md for local startup and CONTRIBUTING.md for checks required as features are implemented.
