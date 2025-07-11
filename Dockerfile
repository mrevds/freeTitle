# Используем официальный образ Go
FROM golang:1.24-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]