# Используем официальный Go образ под ARM64
FROM golang:1.24.4

# Рабочая директория
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник прямо здесь
RUN go build -o app ./services/restaurant

# Делаем исполняемым на всякий случай
RUN chmod +x app

# Точка входа
CMD ["./app"]
