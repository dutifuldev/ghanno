FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/prtags ./cmd/prtags

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN useradd --system --create-home --uid 10001 prtags

WORKDIR /app

COPY --from=builder /out/prtags /usr/local/bin/prtags
COPY migrations /app/migrations

ENV PRTAGS_MIGRATIONS_DIR=/app/migrations

USER prtags

ENTRYPOINT ["/usr/local/bin/prtags"]
