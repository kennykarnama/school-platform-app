# syntax=docker/dockerfile:1
FROM golang:1.25.7-alpine AS goose-builder

ARG GOOSE_VERSION=v3.27.1
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=goose-builder /go/bin/goose /usr/local/bin/goose
COPY services/school-administration-api/database/migrations/cockroachdb /migrations/administration
COPY services/school-user-api/database/migrations/cockroachdb /migrations/user
COPY docker/migrations/migrate.sh /usr/local/bin/migrate
RUN chmod 0755 /usr/local/bin/migrate
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/migrate"]
