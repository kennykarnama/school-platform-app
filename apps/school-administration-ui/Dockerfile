# syntax=docker/dockerfile:1
FROM node:22-alpine AS builder

WORKDIR /src
COPY apps/school-administration-ui/package.json apps/school-administration-ui/package-lock.json ./
RUN npm ci
COPY apps/school-administration-ui/ ./
RUN SCHOOL_ADMINISTRATION_API_BASE_URL= npm run build

FROM nginx:1.28-alpine
COPY apps/school-administration-ui/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /src/dist /usr/share/nginx/html
EXPOSE 80
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -q -O- http://127.0.0.1/healthz >/dev/null || exit 1
