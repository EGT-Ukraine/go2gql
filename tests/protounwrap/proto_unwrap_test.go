package protounwrap

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/EGT-Ukraine/go2gql/api/interceptors"
	"github.com/EGT-Ukraine/go2gql/tests"
	"github.com/EGT-Ukraine/go2gql/tests/protounwrap/generated/clients/apis"
	"github.com/EGT-Ukraine/go2gql/tests/protounwrap/generated/schema"
	"github.com/EGT-Ukraine/go2gql/tests/protounwrap/mock"
)

//go:generate mockgen -destination=mock/item.go -package=mock github.com/EGT-Ukraine/go2gql/tests/protounwrap/generated/clients/apis ItemsServiceClient

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

func TestProtoRequestFieldUnwrappingResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().TestRequestUnwrap(gomock.Any(), gomock.Any()).Return(&empty.Empty{}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrap(name: "username")
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrap": null
			}
		}
	}`, response)
}

func TestProtoRequestFieldUnwrappingCorrectRequestFormed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)

	itemsClient.EXPECT().TestRequestUnwrap(gomock.Any(), &apis.TestRequestUnwrapRequest{
		Name: &wrappers.StringValue{
			Value: "username",
		},
	}).Return(&empty.Empty{}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrap(name: "username")
			}
		}`,
	})
}

func TestProtoRequestFieldUnwrappingNestedMessageResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().TestRequestUnwrapInnerMessage(gomock.Any(), gomock.Any()).Return(&apis.TestRequestUnwrapInnerMessageResponse{}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapInnerMessage(payload: [{names: ["username"], activated: true}]){
					names
					activated
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrapInnerMessage": null
			}
		}
	}`, response)
}

func TestProtoRequestFieldUnwrappingNestedMessageCorrectRequestFormed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)

	itemsClient.EXPECT().TestRequestUnwrapInnerMessage(gomock.Any(), &apis.TestRequestUnwrapInnerMessageRequest{
		Payload: []*apis.TestRequestUnwrapInnerMessageRequestPayload{
			{
				Names: []*wrappers.StringValue{
					{Value: "username1"},
					{Value: "username2"},
				},
				Activated: true,
			},
		},
	}).Return(&apis.TestRequestUnwrapInnerMessageResponse{
		Payload: &apis.TestRequestUnwrapInnerMessageResponsePayload{
			Data: &apis.TestRequestUnwrapInnerMessageResponsePayloadData{
				Names: []*wrappers.StringValue{
					{Value: "a"},
					{Value: "b"},
					{Value: "c"},
				},
				Activated: true,
			},
		},
	}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapInnerMessage(payload: [{names: ["username1", "username2"], activated: true}]){
					activated
					names
				}
			}
		}`,
	})
	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrapInnerMessage": {
					"names": ["a","b","c"],
					"activated": true
				}
			}
		}
	}`, response)
}

func TestRequestUnwrapRepeatedMessageResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().TestRequestUnwrapRepeatedMessage(gomock.Any(), gomock.Any()).Return(&apis.TestRequestUnwrapRepeatedMessageResponse{}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapRepeatedMessage(payload: ["username1", "username2"]){
					names
					activated
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrapRepeatedMessage": []
			}
		}
	}`, response)
}

func TestRequestUnwrapRepeatedMessageCorrectRequestFormed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)

	itemsClient.EXPECT().TestRequestUnwrapRepeatedMessage(gomock.Any(), &apis.TestRequestUnwrapRepeatedMessageRequest{
		Payload: &apis.TestRequestUnwrapRepeatedMessageRequestPayload{
			Test: []string{"username1", "username2"},
		},
	}).Return(&apis.TestRequestUnwrapRepeatedMessageResponse{
		Payload: []*apis.TestRequestUnwrapRepeatedMessageResponsePayload{
			{
				Data: []*apis.TestRequestUnwrapRepeatedMessageResponsePayloadData{
					{
						Names: []*wrappers.StringValue{
							{Value: "a"},
							{Value: "b"},
						},
						Activated: true,
					},
					{
						Names: []*wrappers.StringValue{
							{Value: "c"},
							{Value: "d"},
						},
						Activated: false,
					},
				},
			},
		},
	}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapRepeatedMessage(payload: ["username1", "username2"]){
					names
					activated
				}
			}
		}`,
	})
	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrapRepeatedMessage": [
					[
						{
							"names": ["a", "b"],
							"activated": true
						},
						{
							"names": ["c", "d"],
							"activated": false
						}
					]
				]
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

	return tests.MakeRequest(ctx, apiSchema, opts)
}
