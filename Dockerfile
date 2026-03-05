# Сборочный этап
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

# Финальный этап
FROM alpine:latest

# Устанавливаем зависимости для работы yt-dlp и ffmpeg
RUN apk --no-cache add ca-certificates python3 ffmpeg curl \
    && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app

COPY --from=builder /app/bot .

# Создаем папку для временных файлов и конфигов
RUN mkdir -p /tmp /app/configs

CMD ["./bot"]