PROTO_DIR=api/proto
OUT_DIR=generated

gen:
	mkdir -p $(OUT_DIR)
	protoc \
		-I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
