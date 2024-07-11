gen-sensors: 
	@protoc \
		--proto_path=protobuf "protobuf/Scan.proto" \
		--go_out=services/common/genproto/sensors --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/sensors \
		--go-grpc_opt=paths=source_relative
