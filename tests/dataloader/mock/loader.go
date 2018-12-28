package mock

import (
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
)

type LoaderClients struct {
	clients *Clients
}

func NewLoaderClients(clients *Clients) *LoaderClients {
	return &LoaderClients{
		clients: clients,
	}
}

func (l *LoaderClients) GetCategoryServiceClient() config.CategoryServiceClient {
	return l.clients.CategoryClient
}

func (l *LoaderClients) GetCommentsServiceClient() comments_controller.IClient {
	return l.clients.CommentsClient
}

func (l *LoaderClients) GetUserServiceClient() config.UserServiceClient {
	return l.clients.UserClient
}

func (l *LoaderClients) GetItemsReviewServiceClient() config.ItemsReviewServiceClient {
	return l.clients.ReviewsClient
}
