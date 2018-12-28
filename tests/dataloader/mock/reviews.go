package mock

import (
	"context"

	"google.golang.org/grpc"

	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type ReviewsClient struct {
	ListResponse *config.ListResponse
}

func (c *ReviewsClient) List(ctx context.Context, in *config.ListRequest, opts ...grpc.CallOption) (*config.ListResponse, error) {
	return c.ListResponse, nil
}
