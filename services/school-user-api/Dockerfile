# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY services/school-user-api/go.mod services/school-user-api/go.sum ./
RUN go mod download
COPY services/school-user-api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/server /usr/local/bin/server
USER 65532:65532
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD wget -q -O- http://127.0.0.1:8080/healthz >/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/server"]
