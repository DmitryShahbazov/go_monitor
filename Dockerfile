# Использование официального образа Golang как базового
FROM golang:1.21.5 AS builder

# Установка рабочей директории
WORKDIR /app

# Копирование модулей и их установка
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копирование исходного кода и сборка приложения
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Финальный этап, использование легковесного базового образа
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование исполняемого файла в финальный образ
COPY --from=builder /app/ .

# Команда для запуска приложения
CMD ["./myapp"]