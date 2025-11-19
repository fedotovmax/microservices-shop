# =========================
# Этап 1: Сборка бинарника
# =========================
FROM golang:1.25 AS builder

# Рабочая директория внутри контейнера
WORKDIR /app

# Кэшируем зависимости (go.mod/go.sum)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o user_service cmd/user_service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator cmd/migrator/main.go

# =========================
# Этап 2: Минимальный образ
# =========================
FROM alpine:3.20


# Рабочая директория
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/migrator .
COPY --from=builder /app/user_service .
COPY docker-entrypoint.sh .

RUN chmod +x docker-entrypoint.sh

# Опционально: копируем статические файлы (если есть)
# COPY --from=builder /app/static ./static


# Запуск
ENTRYPOINT ["./app/docker-entrypoint.sh"]