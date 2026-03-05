FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates python3 ffmpeg \
    && apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing yt-dlp

COPY --from=builder /app/telegram-bot .

CMD ["./telegram-bot"]
