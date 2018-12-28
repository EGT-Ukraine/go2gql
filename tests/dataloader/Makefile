dataloader_tests:
	mkdir -p generated/clients
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients config/category.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients config/items.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients config/user.proto
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients config/reviews.proto

	swagger generate client --template-dir=../swagger_templates/ -f config/swagger.json -t generated/clients

	go run ../../cmd/go2gql/main.go ../../cmd/go2gql/basic_plugins.go
