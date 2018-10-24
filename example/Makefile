default: proto

proto:
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:.  proto/example.proto
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:.  import/proto/import.proto

.PHONY: proto
