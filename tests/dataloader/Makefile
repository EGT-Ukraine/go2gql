dataloader_tests:
	rm -rf generated/*
	mkdir -p generated/clients
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/category.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/items.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/user.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/reviews.proto

	swagger generate client --template-dir=../swagger_templates/ -f apis/swagger.json -t generated/clients

	go run ../../cmd/go2gql/main.go ../../cmd/go2gql/basic_plugins.go
