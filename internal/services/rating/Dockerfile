FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o rating ./cmd/rating

EXPOSE 8083

CMD ["./rating"]
