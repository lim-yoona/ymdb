service.pb.go: service.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    service.proto
force:
	rm -f service.pb.go
	make service.pb.go
