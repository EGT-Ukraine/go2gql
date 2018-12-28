package mock

import (
	"context"

	"google.golang.org/grpc"

	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type UserClient struct {
	ListResponse *config.UserListResponse
}

func (c *UserClient) List(ctx context.Context, in *config.UserListRequest, opts ...grpc.CallOption) (*config.UserListResponse, error) {
	return c.ListResponse, nil
}
