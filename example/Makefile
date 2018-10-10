
default:
	cd ../ && make
	@protoc -I=${GOPATH}/src:proto:. --go_out=plugins=grpc:proto  proto/example.proto
	@protoc -I=${GOPATH}/src:proto:. --go_out=plugins=grpc:proto  proto/hello/hello.proto
	go2gql

.PHONY: proto
