package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
)

func main() {
	// Подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	// Создаем тестовый заказ
	req := &pb.CreateOrderRequest{
		UserId: "user123",
		RestId: "rest456",
		Items: []*pb.OrderItem{
			{Name: "Pizza", Quantity: 2},
			{Name: "Cola", Quantity: 1},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка создания заказа: %v", err)
	}

	fmt.Printf("Заказ создан успешно! ID: %s\n", resp.OrderId)
}
