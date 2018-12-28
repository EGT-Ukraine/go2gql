package mock

import (
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type Clients struct {
	ItemsClient    config.ItemsServiceClient
	CategoryClient config.CategoryServiceClient
	CommentsClient comments_controller.IClient
	UserClient     config.UserServiceClient
	ReviewsClient  config.ItemsReviewServiceClient
}
