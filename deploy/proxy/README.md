# HTTPS Reverse Proxy

Nginx terminates TLS and reverse-proxies to the application containers. Certbot issues and renews Let's Encrypt certificates automatically.

## Requirements

- Application stack (`../compose.yaml`) running on the VM
- Domain with an A record pointing to the VM external IP
- GCP firewall rules allowing inbound TCP on ports 80 and 443

## First setup

```sh
cp .env.example .env
chmod 600 .env
```

Edit `.env` and set `DOMAIN` and `CERT_EMAIL`. Verify the network name matches the application stack:

```sh
docker network ls | grep school
```

Adjust `APP_NETWORK` in `.env` if it differs from `school-platform_default`.

### Generate a temporary certificate

The HTTPS server block references certificate files that do not exist yet. Generate a self-signed certificate so Nginx can start, then replace it with a real one from Let's Encrypt:

```sh
set -a; . ./.env; set +a
mkdir -p certbot/conf/live/"$DOMAIN"
openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout certbot/conf/live/"$DOMAIN"/privkey.pem \
  -out certbot/conf/live/"$DOMAIN"/fullchain.pem \
  -subj "/CN=$DOMAIN" -days 1
```

### Start Nginx and obtain the real certificate

```sh
docker compose up -d nginx
```

Remove the temporary certificate so Certbot can create its own directory structure:

```sh
rm -rf certbot/conf/live/"$DOMAIN" certbot/conf/archive/"$DOMAIN" certbot/conf/renewal/"$DOMAIN".conf
```

```sh
docker compose run --rm certbot certonly \
  --webroot --webroot-path /var/www/certbot \
  -d "$DOMAIN" \
  --email "$CERT_EMAIL" \
  --agree-tos --no-eff-email
docker compose exec nginx nginx -s reload
docker compose up -d certbot-renew
```

The `certbot` service uses the image's default entrypoint for one-shot commands. The `certbot-renew` service runs the renewal loop. The `certbot` service is excluded from `docker compose up` via a profile — it is only used through `docker compose run`.

The ACME challenge is served over HTTP (port 80) and is not affected by the self-signed certificate.

## Verification

```sh
curl -I http://"$DOMAIN"       # expect 301 → https
curl -I https://"$DOMAIN"      # expect 200
docker compose run --rm certbot certificates
```

## Renewal

The `certbot-renew` service runs `certbot renew` every 12 hours. Certificates are valid for 90 days and renewable 30 days before expiry.

Nginx must be reloaded after a successful renewal. Add a cron job on the VM:

```sh
sudo crontab -e
```

```
0 3 * * * docker exec nginx nginx -s reload
```

## Post-setup

### Secure session cookies

Set `SESSION_COOKIE_SECURE=true` in `../.env` and recreate the application stack:

```sh
cd ..
sed -i 's/^SESSION_COOKIE_SECURE=.*/SESSION_COOKIE_SECURE=true/' .env
docker compose up -d --remove-orphans
```

### Restrict application ports

Application containers no longer need to expose ports publicly. Bind them to localhost so they are only reachable through the proxy:

```sh
# In ../.env
UI_PORT=127.0.0.1:3000
ADMIN_API_PORT=127.0.0.1:8081
USER_API_PORT=127.0.0.1:8082
```

Recreate the application stack after changing ports.

## Logs

```sh
docker compose logs -f nginx
docker compose logs certbot-renew
```
