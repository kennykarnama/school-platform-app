FROM golang:1.25.7-alpine AS goose-builder

ARG GOOSE_VERSION=v3.27.1
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}

FROM cockroachdb/cockroach:v26.2.2
COPY --from=goose-builder /go/bin/goose /usr/local/bin/goose

