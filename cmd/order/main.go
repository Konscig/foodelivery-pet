package order

import "github.com/Konscig/foodelivery-pet/generated/orderpb"

func main() {
	cfg := config.MustLoad()

	pg := bootstrap.InitPostgres(cfg)
	redis := bootstrap.InitRedis(cfg)

	producer := bootstrap.InitKafkaProducer(cfg, cfg.Kafka.OrderCreatedTopic)

	orderPublisher := kafkaadapter.NewOrderProducer(producer)
	orderService := order.NewService(pg, redis, orderPublisher)

	grpcServer := bootstrap.InitGRPCServer()
	orderpb.RegisterOrderServiceServer(
		grpcServer,
		grpcadapter.NewOrderServer(orderService),
	)

	bootstrap.RunGRPC(grpcServer, cfg.GRPC.Port)
}
