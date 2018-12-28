package mock

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type ItemsServiceClient struct {
	GetOneResponse *config.Item
	ListResponse   *config.ItemListResponse
}

func (c *ItemsServiceClient) List(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*config.ItemListResponse, error) {
	return c.ListResponse, nil
}

func (c *ItemsServiceClient) GetOne(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*config.Item, error) {
	return c.GetOneResponse, nil
}
