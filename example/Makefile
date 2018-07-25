default: proto

proto:
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:.  proto/example.proto

.PHONY: proto
