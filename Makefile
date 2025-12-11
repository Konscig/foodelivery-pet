PROTO_DIR=api/proto
OUT_DIR=generated

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
