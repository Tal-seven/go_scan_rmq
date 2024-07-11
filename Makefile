run-orders: 
	@go run services/orders/*.go

run-kitchen:
	@go run services/kitchen/*.go


gen-orders: 
	@protoc \
		--proto_path=protobuf "protobuf/orders.proto" \
		--go_out=services/common/genproto/orders --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/orders \
		--go-grpc_opt=paths=source_relative

gen-sensors: 
	@protoc \
		--proto_path=protobuf "protobuf/Scan.proto" \
		--go_out=services/common/genproto/sensors --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/sensors \
		--go-grpc_opt=paths=source_relative
