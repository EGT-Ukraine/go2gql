package mock

import (
	"context"

	"google.golang.org/grpc"

	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type CategoryServiceClient struct {
	ListResponse *config.CategoryListResponse
}

func (c *CategoryServiceClient) List(ctx context.Context, in *config.CategoryListRequest, opts ...grpc.CallOption) (*config.CategoryListResponse, error) {
	return c.ListResponse, nil
}
