default: proto

proto:
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:${GOPATH}/src  test.proto
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:${GOPATH}/src  common/common.proto
	@go run ../cmd/proto2gql/main.go

.PHONY: proto
