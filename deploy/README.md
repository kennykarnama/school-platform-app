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

## Logs and migration status

```sh
docker compose logs migrations
docker compose logs -f school-administration-api school-user-api school-administration-ui
```

## Pin or roll back

Set `IMAGE_TAG` in `.env` to the desired release without the leading `v`, then recreate the stack:

```sh
sed -i 's/^IMAGE_TAG=.*/IMAGE_TAG=1.2.3/' .env
docker compose pull
docker compose up -d --remove-orphans
```

Database migrations are forward-applied before startup. Rolling application images back does not automatically roll database migrations back.

## HTTPS

When TLS terminates at a reverse proxy in front of the UI, set `SESSION_COOKIE_SECURE=true`. Keep it `false` only when accessing the UI directly over HTTP.
