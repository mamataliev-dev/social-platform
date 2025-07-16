# Generate swagger api
#protoc \
#  -I api/proto \
#  -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.1 \
#  -I api/proto/third_party/googleapis \
#  -I api/proto/third_party/protoc-gen-openapiv2 \
#  --openapiv2_out=docs/swagger \
#  api/proto/user_auth/user_auth.proto

# Generate proto files
#proto:
#	protoc \
#	  --proto_path=shared/proto \
#	  --go_out=shared/gen/userpb --go_opt=paths=source_relative \
#	  --go-grpc_out=shared/gen/userpb --go-grpc_opt=paths=source_relative \
#	  shared/proto/user.proto
