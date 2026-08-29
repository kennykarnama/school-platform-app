# School Platform App

Monorepo for the school administration UI, administration API, and user API.

## Layout

- `apps/school-administration-ui`: Riot/Webpack administration interface
- `services/school-administration-api`: attendance and school administration API
- `services/school-user-api`: user registration and authentication API

The Go services remain independent modules and are coordinated by the root `go.work` file.

## Prerequisites

- Go 1.22 or newer
- Node.js 22 (see `.nvmrc`)
- Docker with Compose v2

## First-time setup

```sh
make setup
make deps-up
```

`make deps-up` starts CockroachDB and Redis, creates both local databases, applies pending migrations, and loads idempotent seed data.

Local development logins:

| Service | Username | Password |
| --- | --- | --- |
| Administration UI/API | `teacher.demo` | `school123` |
| User API | `demo.user` | `school123` |

The administration seed also includes two academic years, four attendance types, six students in `KELAS I A`, and attendance records for August 24–26, 2026. These credentials and records are for local development only.

Run the complete development stack:

```sh
make dev
```

Or run each host process in a separate terminal:

```sh
make run-administration
make run-user
make run-ui
```

Local endpoints:

- Administration UI: http://localhost:3000
- Administration API: http://localhost:8081
- User API: http://localhost:8082
- CockroachDB SQL: localhost:26257
- CockroachDB Console: http://localhost:8088
- Redis: localhost:6379

## Verification and cleanup

```sh
make test
make build
make deps-down
```

To delete only this project's local CockroachDB and Redis volumes and rebuild the databases:

```sh
make db-reset
```

The checked-in files under `config/local` contain development-only values. Production deployments must provide their own database TLS settings and secrets.
