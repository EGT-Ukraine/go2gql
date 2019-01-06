package mock

import (
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller"
)

type LoaderClients struct {
	clients *Clients
}

type Clients struct {
	ItemsClient    apis.ItemsServiceClient
	CategoryClient apis.CategoryServiceClient
	CommentsClient comments_controller.IClient
	UserClient     apis.UserServiceClient
	ReviewsClient  apis.ItemsReviewServiceClient
}

func NewLoaderClients(clients *Clients) *LoaderClients {
	return &LoaderClients{
		clients: clients,
	}
}

func (l *LoaderClients) GetCategoryServiceClient() apis.CategoryServiceClient {
	return l.clients.CategoryClient
}

func (l *LoaderClients) GetCommentsServiceClient() comments_controller.IClient {
	return l.clients.CommentsClient
}

func (l *LoaderClients) GetUserServiceClient() apis.UserServiceClient {
	return l.clients.UserClient
}

func (l *LoaderClients) GetItemsReviewServiceClient() apis.ItemsReviewServiceClient {
	return l.clients.ReviewsClient
}
