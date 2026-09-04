# School Platform VM Deployment

## Requirements

- Linux amd64 VM with Docker Engine and Docker Compose v2
- Existing `school_administration` and `school_user` CockroachDB databases
- Network access from the VM to CockroachDB Cloud

The images are public, so Docker Hub login is not required.

## First deployment

```sh
cp .env.example .env
chmod 600 .env
```

Edit `.env`, replace both database URLs and both JWT secrets, then run:

```sh
docker compose pull
docker compose up -d
docker compose ps
```

The one-shot `migrations` service applies both Goose migration sets before either API starts. It never loads demo seed data. The UI is available on port `3000`, the administration API on `8081`, and the user API on `8082` by default.

## Upgrade to the latest stable release

```sh
docker compose pull
docker compose up -d --remove-orphans
docker compose ps
```

## Deploy the experimental main build

Set `IMAGE_TAG=experimental` in `.env`, then pull and recreate the stack:

```sh
docker compose pull
docker compose up -d --remove-orphans
docker compose ps
```

The `experimental` tag follows the newest verified commit on `main` and must not be treated as a stable production release. For a reproducible test deployment or rollback, replace it with the corresponding immutable `sha-...` tag shown by the workflow.

## Deploy a manually built branch

Run **Build and release Docker images** from the GitHub Actions page and select the desired branch in the **Run workflow** dialog. The images are published with a sanitized `branch-<branch-name>` tag; for example, `feature/tenant-admin` becomes `branch-feature-tenant-admin`.

Set that value in `.env`, or export it for the current shell, then recreate the stack:
## Deploy an experimental or branch build

Pushes to `main` publish the `experimental` image tag. Other branches can be built manually from **Actions → Build and release Docker images → Run workflow** and receive a sanitized `branch-<branch-name>` tag. For example, `feature/tenant-admin` becomes `branch-feature-tenant-admin`.

Set the selected tag in `.env`, or export it for the current shell, then recreate the stack:

```sh
export IMAGE_TAG=branch-feature-tenant-admin
docker compose pull
docker compose up -d --remove-orphans
```

Use the accompanying immutable `sha-...` tag when the branch deployment must be reproducible.
Use the accompanying immutable `sha-...` tag when the deployment must be reproducible.

## Logs and migration status

```sh
docker compose logs migrations
docker compose logs -f school-administration-api school-user-api school-administration-ui
```

## Bootstrap the platform administrator

Migration 20 moves existing data into the `legacy` school and promotes its oldest active teacher to school administrator. The first global platform administrator is created directly in CockroachDB; there is intentionally no public bootstrap endpoint.

Generate a bcrypt hash outside the database and keep the plaintext password out of shell history and source control. Then insert the account into the administration database with a globally unique username:

```sql
INSERT INTO teacher (
  alternative_id, name, password, school_id, role, active,
  must_change_password, created_at, updated_at
) VALUES (
  'platform.admin', 'Platform Administrator', '<BCRYPT_HASH>', NULL,
  'platform_admin', true, true, now(), now()
);
```

The administrator must replace the temporary password after first login. From the **Sekolah** page, they can create each tenant and its first school administrator.

## Pin or roll back

Set `IMAGE_TAG` in `.env` to the desired release without the leading `v`, then recreate the stack:

```sh
sed -i 's/^IMAGE_TAG=.*/IMAGE_TAG=1.2.3/' .env
docker compose pull
docker compose up -d --remove-orphans
```

Database migrations are forward-applied before startup. Rolling application images back does not automatically roll database migrations back.

## HTTPS

A ready-to-use Nginx + Let's Encrypt reverse proxy is in `proxy/`. See `proxy/README.md` for setup instructions.

When TLS terminates at a reverse proxy in front of the UI, set `SESSION_COOKIE_SECURE=true`. Keep it `false` only when accessing the UI directly over HTTP.
