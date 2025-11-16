FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

RUN go install github.com/go-task/task/v3/cmd/task@latest \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/backend ./cmd/main.go

RUN task install:migrator


FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/backend/taskfile.yml ./taskfile.yml

COPY --from=builder /app/backend/bin/backend ./bin/backend

COPY --from=builder /go/bin/task /usr/local/bin/task
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY backend/db/migrations ./db/migrations

ENTRYPOINT ["./bin/backend"]
