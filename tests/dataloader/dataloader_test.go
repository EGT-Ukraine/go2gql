package dataloader_test

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/EGT-Ukraine/go2gql/api/interceptors"
	"github.com/EGT-Ukraine/go2gql/tests"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/config"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/clients/models"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/schema"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/generated/schema/loaders"
	"github.com/EGT-Ukraine/go2gql/tests/dataloader/mock"
)

func TestDataLoader(t *testing.T) {
	itemsListResponse := &config.ItemListResponse{
		Items: []*config.Item{
			{
				Name:       "item 1",
				CategoryId: 12,
			},
		},
	}

	itemsClient := &mock.ItemsServiceClient{ListResponse: itemsListResponse}

	categoryListResponse := &config.CategoryListResponse{
		Categories: []*config.Category{
			{
				Name: "category 1",
			},
		},
	}

	categoryClient := &mock.CategoryServiceClient{ListResponse: categoryListResponse}

	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CategoryClient: categoryClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					items {
						category {
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
				"list": {
					"items": [
						{
							"category": {
								"name": "category 1"
							}
						}
					]
				}
			}
		}
	}`, response)
}

func TestDataLoaderGetOne(t *testing.T) {
	itemsClient := &mock.ItemsServiceClient{
		GetOneResponse: &config.Item{
			Name:       "item 1",
			CategoryId: 12,
		},
	}

	categoryListResponse := &config.CategoryListResponse{
		Categories: []*config.Category{
			{
				Name: "category 1",
			},
		},
	}

	categoryClient := &mock.CategoryServiceClient{ListResponse: categoryListResponse}

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
	itemsListResponse := &config.ItemListResponse{
		Items: []*config.Item{
			{
				Id:         44,
				Name:       "item 1",
				CategoryId: 12,
			},
		},
	}

	itemsClient := &mock.ItemsServiceClient{ListResponse: itemsListResponse}

	comments := [][]*models.ItemComment{
		{
			{
				ID:     111,
				UserID: 54,
				Text:   "test comment",
			},
		},
	}

	commentsClient := &mock.CommentsClient{Comments: comments}

	userListResponse := &config.UserListResponse{
		Users: []*config.User{
			{
				Id:   54,
				Name: "Test User",
			},
		},
	}

	userClient := &mock.UserClient{ListResponse: userListResponse}
	clients := &mock.Clients{
		ItemsClient:    itemsClient,
		CommentsClient: commentsClient,
		UserClient:     userClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					items {
						comments {
							id
							text
							user {
								name
							}
						}
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"list": {
					"items": [
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
		}
	}`, response)
}

func TestDataLoaderGetOneWithGRPCArrayLoader(t *testing.T) {
	itemsClient := &mock.ItemsServiceClient{
		GetOneResponse: &config.Item{
			Name:       "item 1",
			CategoryId: 12,
		},
	}

	listResponse := &config.ListResponse{
		ItemReviews: []*config.ItemReviews{
			{
				ItemReview: []*config.Review{
					{
						Id:   456,
						Text: "excellent item",
					},
				},
			},
		},
	}

	reviewsClient := &mock.ReviewsClient{ListResponse: listResponse}

	clients := &mock.Clients{
		ItemsClient:   itemsClient,
		ReviewsClient: reviewsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				GetOne {
					reviews {
						item_review {
							text
						}
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"GetOne": {
					"reviews": {
						"item_review": [
							{
								"text": "excellent item"
							}
						]
					}
				}
			}
		}
	}`, response)

	// TODO: after proto field unwrapping implementation should be:
	//tests.AssertJSON(t, `{
	//	"data": {
	//		"items": {
	//			"GetOne": {
	//				"reviews": [
	//					{
	//						"text": "excellent item"
	//					}
	//				]
	//			}
	//		}
	//	}
	//}`, response)
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
