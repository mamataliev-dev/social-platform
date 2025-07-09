proto:
	protoc \
	  --proto_path=shared/proto \
	  --go_out=shared/gen/userpb --go_opt=paths=source_relative \
	  --go-grpc_out=shared/gen/userpb --go-grpc_opt=paths=source_relative \
	  shared/proto/user.proto
