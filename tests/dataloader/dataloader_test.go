package dataloader_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/EGT-Ukraine/go2gql/api/interceptors"
	"github.com/EGT-Ukraine/go2gql/tests"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/models"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/schema"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/schema/loaders"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/mock"
)

//go:generate mockgen -destination=mock/category.go -package=mock github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis CategoryServiceClient
//go:generate mockgen -destination=mock/item.go -package=mock github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis ItemsServiceClient
//go:generate mockgen -destination=mock/comments.go -package=mock github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/client/comments_controller IClient
//go:generate mockgen -destination=mock/reviews.go -package=mock github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis ItemsReviewServiceClient
//go:generate mockgen -destination=mock/user.go -package=mock github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/apis UserServiceClient

func TestDataLoader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.ItemListResponse{
		Items: []*apis.Item{
			{
				Name:       "item 1",
				CategoryId: 12,
			},
			{
				Name:       "item 2",
				CategoryId: 11,
			},
		},
	}, nil).AnyTimes()

	categoryClient := mock.NewMockCategoryServiceClient(mockCtrl)

	categoryClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.CategoryListResponse{
		Categories: []*apis.Category{
			{
				Name: "category 12",
			},
			{
				Name: "category 11",
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CategoryClient: categoryClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					name
					category {
						name
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"list": [
						{
							"name": "item 1",
							"category": {
								"name": "category 12"
							}
						},
						{
							"name": "item 2",
							"category": {
								"name": "category 11"
							}
						}
				]
			}
		}
	}`, response)
}

func TestDataLoaderServiceMakeOnlyOneRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.ItemListResponse{
		Items: []*apis.Item{
			{
				Name:       "item 1",
				CategoryId: 12,
			},
			{
				Name:       "item 2",
				CategoryId: 11,
			},
		},
	}, nil).AnyTimes()

	categoryClient := mock.NewMockCategoryServiceClient(mockCtrl)

	categoryClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.CategoryListResponse{
		Categories: []*apis.Category{
			{
				Name: "category 12",
			},
			{
				Name: "category 11",
			},
		},
	}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CategoryClient: categoryClient,
	}

	makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					name
					category {
						name
					}
				}
			}
		}`,
	})
}

func TestDataLoaderGetOne(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(&apis.Item{
		Name:       "item 1",
		CategoryId: 12,
	}, nil).AnyTimes()

	categoryClient := mock.NewMockCategoryServiceClient(mockCtrl)

	categoryClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.CategoryListResponse{
		Categories: []*apis.Category{
			{
				Name: "category 1",
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CategoryClient: categoryClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				GetOne {
					category {
						name
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"GetOne": {
					"category": {
						"name": "category 1"
					}
				}
			}
		}
	}`, response)
}

func TestDataLoaderWithSwagger(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.ItemListResponse{
		Items: []*apis.Item{
			{
				Id:         44,
				Name:       "item 1",
				CategoryId: 12,
			},
		},
	}, nil).AnyTimes()

	commentsClient := mock.NewMockIClient(mockCtrl)
	commentsClient.EXPECT().ItemsComments(gomock.Any()).Return(&comments_controller.ItemsCommentsOK{
		Payload: [][]*models.ItemComment{
			{
				{
					ID:     111,
					UserID: 54,
					Text:   "test comment",
				},
			},
		},
	}, nil).AnyTimes()

	userClient := mock.NewMockUserServiceClient(mockCtrl)
	userClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.UserListResponse{
		Users: []*apis.User{
			{
				Id:   54,
				Name: "Test User",
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CommentsClient: commentsClient,
		UserClient:     userClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					comments {
						id
						text
						user {
							name
						}
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"list": [
					{
						"comments": [
							{
								"id": 111,
								"text": "test comment",
								"user": {
	                  				"name": "Test User"
	                			}
							}
						]
					}
				]
			}
		}
	}`, response)
}

func TestDataLoaderWithProtoFieldUnwrapping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(&apis.Item{
		Name:       "item 1",
		CategoryId: 12,
	}, nil).AnyTimes()

	reviewsClient := mock.NewMockItemsReviewServiceClient(mockCtrl)
	reviewsClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.ListResponse{
		ItemReviews: []*apis.ItemReviews{
			{
				ItemReview: []*apis.Review{
					{
						Id:   456,
						Text: "excellent item",
					},
				},
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient:   itemsClient,
		ReviewsClient: reviewsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				GetOne {
					reviews {
						text
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"GetOne": {
					"reviews": [
						{
							"text": "excellent item"
						}
					]
				}
			}
		}
	}`, response)
}

func makeRequest(t *testing.T, clients *mock.Clients, opts *handler.RequestOptions) *graphql.Result {
	schemaClients := schema.APISchemaClients{
		ItemsServiceClient: clients.ItemsClient,
	}

	apiSchema, err := schema.GetAPISchema(schemaClients, &interceptors.InterceptorHandler{})

	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx := context.Background()

	loadersAPIClients := mock.NewLoaderClients(clients)

	ctx = loaders.GetContextWithLoaders(ctx, loadersAPIClients)

	return tests.MakeRequest(apiSchema, opts, ctx)
}
