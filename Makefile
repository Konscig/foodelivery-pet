COMPOSE=docker-compose

PROTO_DIR=api/proto
OUT_DIR=generated

ORDER_DIR=services/order
RESTAURANT_DIR=services/restaurant
DELIVERY_DIR=services/delivery

.PHONY: build up down logs clean

## Генерация protobuf
gen:
	mkdir -p $(OUT_DIR)/orderpb
	mkdir -p $(OUT_DIR)/eventspb

	protoc -I=$(PROTO_DIR) \
		--go_out=paths=source_relative:$(OUT_DIR)/orderpb \
		--go-grpc_out=paths=source_relative:$(OUT_DIR)/orderpb \
		$(PROTO_DIR)/order_models.proto \
		$(PROTO_DIR)/order_service.proto

	protoc -I=$(PROTO_DIR) \
		--go_out=paths=source_relative:$(OUT_DIR)/eventspb \
		--go-grpc_out=paths=source_relative:$(OUT_DIR)/eventspb \
		$(PROTO_DIR)/order_events.proto

build:
	docker build -t fd-order -f services/order/Dockerfile .
	docker build -t fd-restaurant -f services/restaurant/Dockerfile .
	docker build -t fd-delivery -f services/delivery/Dockerfile .

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

clean:
	$(COMPOSE) down -v
	docker system prune -af

## Создание Kafka топиков
create-topics:
	$(COMPOSE) exec kafka kafka-topics --bootstrap-server localhost:9092 --create --topic order.created --partitions 1 --replication-factor 1 || true
	$(COMPOSE) exec kafka kafka-topics --bootstrap-server localhost:9092 --create --topic order.ready --partitions 1 --replication-factor 1 || true
	$(COMPOSE) exec kafka kafka-topics --bootstrap-server localhost:9092 --create --topic order.coming --partitions 1 --replication-factor 1 || true
	$(COMPOSE) exec kafka kafka-topics --bootstrap-server localhost:9092 --create --topic order.done --partitions 1 --replication-factor 1 || true