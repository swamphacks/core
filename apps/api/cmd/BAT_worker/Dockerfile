FROM golang:1.25-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Dev 
FROM base AS dev

RUN go install github.com/air-verse/air@latest

CMD ["air"]

# Production
FROM base AS prod

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o BAT_worker ./cmd/BAT_worker

RUN apk --no-cache add ca-certificates

CMD ["./BAT_worker"]


