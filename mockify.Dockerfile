FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mock ./utils/mockify

RUN chmod +x mock

CMD ["./mock"]