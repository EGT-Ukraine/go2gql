default: proto generate

proto:
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:.  proto/example.proto

generate:
	cd ../ && $(MAKE) build_templates
	go run ../cmd/go2gql/

.PHONY: proto
