# School Platform App

Monorepo for the school administration UI, administration API, and user API.

## Layout

- `apps/school-administration-ui`: Riot/Webpack administration interface
- `services/school-administration-api`: attendance and school administration API
- `services/school-user-api`: user registration and authentication API

The Go services remain independent modules and are coordinated by the root `go.work` file.

The administration application is multi-tenant. School data is isolated by `school_id`; globally unique usernames sign in without selecting a school. Platform administrators create and activate schools, school administrators manage their school's setup, teachers, classes, and access assignments, and teachers can only work with the academic-year/class combinations assigned to them. New and reset accounts must replace their temporary password on first login.

## Prerequisites

- Go 1.22 or newer
- Node.js 22 (see `.nvmrc`)
- Docker with Compose v2

## First-time setup

```sh
make setup
make deps-up
```

`make deps-up` starts CockroachDB, creates both local databases, applies pending migrations, and loads idempotent seed data.

Schema migrations are managed by [Goose](https://github.com/pressly/goose). The pinned Goose runner is built automatically by Docker, so a host installation is not required. Goose records applied versions in each database's `goose_db_version` table; seed SQL remains separate from schema migrations.

Create and validate new sequential migrations with:

```sh
make db-new-administration NAME=add_school_profile
make db-new-user NAME=add_user_status
make db-validate
```

Every SQL migration must contain one `-- +goose Up` section. Add a later `-- +goose Down` section when the change can be rolled back safely.

Local development logins:

| Service | Username | Password |
| --- | --- | --- |
| Administration UI/API | `teacher.demo` | `school123` |
| User API | `demo.user` | `school123` |

The administration seed also includes two academic years, four attendance types, six students in `KELAS I A`, and attendance records for August 24–26, 2026. These credentials and records are for local development only.

After signing in, open **Data Awal** in the administration UI to add or update academic years, attendance types, students, and class assignments. Student data can be entered in the form or imported from Excel (`.xlsx`) or CSV. Use **Unduh Template Excel** for a ready-to-fill workbook, or provide either file format with these exact headers:

```csv
alternativeID,name,academicYearLabel,classLabel
,Nama Siswa,2026/2027 - Semester 1,KELAS I A
```

Leave `alternativeID` empty to generate a unique `STU-…` value during preview. Use **Validasi & Pratinjau** before applying. Imports accept up to 5,000 student rows, upsert matching records, and never delete records omitted from the form or imported file. A student is matched by alternative ID, and each student may have one active class assignment per academic year. Student names are not treated as unique. The Go API generates the workbook at `GET /api/v1/setup/students/template` and parses multipart Excel/CSV uploads at `POST /api/v1/setup/students/import`; the authenticated preview/apply endpoints then validate and persist the returned rows.

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

## Verification and cleanup

```sh
make test
make build
make deps-down
```

To delete only this project's local CockroachDB volume and rebuild the databases:

```sh
make db-reset
```

The checked-in files under `config/local` contain development-only values. Production deployments must provide their own database TLS settings and secrets.

The user API stores JWT session and revocation state in the `school_user.jwt_session` table. Redis is not required.

## Docker releases

Pushing to `main` or pushing a semantic-version tag such as `v1.2.3` runs the workflow in `.github/workflows/release.yml`. It verifies the monorepo, builds `linux/amd64` images, and pushes these public Docker Hub repositories:

- `kennykarnama/school-administration-api`
- `kennykarnama/school-user-api`
- `kennykarnama/school-administration-ui`
- `kennykarnama/school-platform-migrations`

Before the first release, create those four repositories as **Public** in Docker Hub. Create a Docker Hub personal access token with read/write permission, then add it to the GitHub repository as the Actions secret `DOCKERHUB_TOKEN`.

Create a stable release with:

```sh
git tag v1.0.0
git push origin v1.0.0
```

Stable tags publish `1.0.0`, `sha-...`, and `latest` image tags. Prerelease tags such as `v1.1.0-rc.1` publish the version and commit tags without moving `latest`.

Every verified push to `main` publishes `experimental` and `sha-...` tags for all four images. `experimental` is mutable and intended for integration testing; use its immutable `sha-...` companion when a deployment must be reproducible. Main-branch builds do not move `latest` or create GitHub Releases.

To build another branch, open **Actions → Build and release Docker images → Run workflow**, select the branch, and start the workflow. A branch named `feature/tenant-admin`, for example, publishes `branch-feature-tenant-admin` plus the immutable `sha-...` tag. Manual branch runs do not update `experimental` or `latest` and do not create a GitHub Release.

After the images are published, the workflow creates a GitHub Release containing a VM-ready `school-platform-app-v1.0.0-linux-amd64.tar.gz` archive and its SHA-256 checksum. The archive contains Compose configuration and deployment instructions; its default `IMAGE_TAG=latest` pulls the newest stable images.

## Hosted CockroachDB

Both APIs accept a complete CockroachDB PostgreSQL URL through `DATABASE_URL`. It takes precedence over the individual `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME`, `DB_CLUSTER`, and `DB_SSL_MODE` variables.

Create the two application databases on the cluster, then configure each deployed service with its own secret:

```sh
# school-administration-api
DATABASE_URL='postgresql://USER:PASSWORD@HOST:26257/school_administration?sslmode=verify-full'

# school-user-api
DATABASE_URL='postgresql://USER:PASSWORD@HOST:26257/school_user?sslmode=verify-full'
```

Use the connection URLs generated by the CockroachDB Cloud console when possible, because some clusters include additional routing parameters. Do not commit either URL. If a database client reports that multi-database mode is disabled, enable its **Show all databases** connection option; that setting controls the client UI, not this application.

Apply or inspect Goose migrations on a hosted cluster by supplying the URL for one database at a time:

```sh
DATABASE_URL='postgresql://USER:PASSWORD@HOST:26257/school_administration?sslmode=verify-full' \
  make db-migrate-administration

DATABASE_URL='postgresql://USER:PASSWORD@HOST:26257/school_user?sslmode=verify-full' \
  make db-migrate-user

DATABASE_URL='postgresql://USER:PASSWORD@HOST:26257/school_administration?sslmode=verify-full' \
  make db-status-administration
```

For a local insecure cluster, use `sslmode=disable`. Hosted CockroachDB should retain the TLS parameters supplied by its connection dialog. The remote migration targets apply schema only; they do not load demo seed data.
