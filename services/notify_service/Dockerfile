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
RUN CGO_ENABLED=0 GOOS=linux go build -o notify_service cmd/notify_service/main.go

# =========================
# Этап 2: Минимальный образ
# =========================
FROM alpine:3.20


# Рабочая директория
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/notify_service .
COPY --from=builder /app/translations ./translations

CMD [ "/app/notify_service" ]

