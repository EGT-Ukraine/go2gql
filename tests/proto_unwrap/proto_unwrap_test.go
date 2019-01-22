package proto_unwrap

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/EGT-Ukraine/go2gql/api/interceptors"
	"github.com/EGT-Ukraine/go2gql/tests"
	"github.com/EGT-Ukraine/go2gql/tests/proto_unwrap/generated/clients/apis"
	"github.com/EGT-Ukraine/go2gql/tests/proto_unwrap/generated/schema"
	"github.com/EGT-Ukraine/go2gql/tests/proto_unwrap/mock"
)

//go:generate mockgen -destination=mock/item.go -package=mock github.com/EGT-Ukraine/go2gql/tests/proto_unwrap/generated/clients/apis ItemsServiceClient

func TestProtoResponseFieldUnwrapping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().List(gomock.Any(), gomock.Any()).Return(&apis.ItemListResponse{
		Items: []*apis.Item{
			{
				Id:   11,
				Name: "item 1",
			},
			{
				Id:   12,
				Name: "item 2",
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				list {
					id
					name
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"list": [
					{
						"id": 11,
						"name": "item 1"
					},
					{
						"id": 12,
						"name": "item 2"
					}
				]
			}
		}
	}`, response)
}

func TestProtoResponseScalarFieldUnwrapping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().Activated(gomock.Any(), gomock.Any()).Return(&apis.ActivatedResponse{
		Activated: true,
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				activated
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"activated": true
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

	return tests.MakeRequest(apiSchema, opts, ctx)
}
