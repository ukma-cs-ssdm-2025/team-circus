FROM golang:1.25.1 AS builder
WORKDIR /app

RUN go install github.com/go-task/task/v3/cmd/task@latest
COPY . .
WORKDIR /app/backend
RUN go mod download 

RUN go build -o ./bin/backend cmd/main.go
RUN task install:migrator

FROM debian:bookworm
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/

COPY --from=builder /app/backend/taskfile.yml ./taskfile.yml

COPY --from=builder /app/backend/bin/backend ./bin/backend
COPY --from=builder /go/bin/task /usr/local/bin/task
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY backend/db/migrations ./db/migrations