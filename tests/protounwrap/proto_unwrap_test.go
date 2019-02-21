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

func TestProtoDeepResponseFieldUnwrapping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().GetDeep(gomock.Any(), gomock.Any()).Return(&apis.GetDeepResponse{
		Payload: &apis.GetDeepResponsePayload{
			Data: &apis.GetDeepResponsePayloadData{
				Id:   "123",
				Name: "456",
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				deep {
					id
					name
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"deep": {
					"id": "123",
					"name": "456"
				}
			}
		}
	}`, response)
}
func TestProtoListDeepRepeatedResponseFieldUnwrapping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().ListDeepRepeated(gomock.Any(), gomock.Any()).Return(&apis.ListDeepRepeatedResponse{
		Payload: []*apis.ListDeepRepeatedResponsePayload{
			{
				Data: []*apis.ListDeepRepeatedResponsePayloadData{
					{
						Id:   "123",
						Name: "456",
					},
					{
						Id:   "789",
						Name: "101112",
					},
				},
			},
			{
				Data: []*apis.ListDeepRepeatedResponsePayloadData{
					{
						Id:   "222",
						Name: "444",
					},
					{
						Id:   "555",
						Name: "666",
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
				listDeepRepeated {
					id
					name
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"listDeepRepeated": [
					[
						{
							"id": "123",
							"name": "456"
						},
						{
							"id": "789",
							"name": "101112"
						}
					],
					[
						{
							"id": "222",
							"name": "444"
						},
						{
							"id": "555",
							"name": "666"
						}
					]
				]
			}
		}
	}`, response)
}

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

func TestMapUnwrappingResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().MapUnwrap(gomock.Any(), gomock.Any()).Return(&apis.MapUnwrapResponse{
		Items: map[string]*wrappers.StringValue{
			"key": {Value: "value"},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				mapUnwrap {
					items {
						key
						value
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"mapUnwrap": {
					"items": [
						{
							"key": "key",
							"value": "value"
				  		}
					]
				}
			}
		}
	}`, response)
}
func TestMapRepeatedUnwrappingResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().MapRepeatedUnwrap(gomock.Any(), gomock.Any()).Return(&apis.MapRepeatedUnwrapResponse{
		Items: map[string]*apis.RepeatedUnwrapResponse{
			"key": {
				Items: []string{"1", "2", "3"},
			},
		},
	}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				mapRepeatedUnwrap {
					items {
						key
						value
					}
				}
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"mapRepeatedUnwrap": {
					"items": [
						{
							"key": "key",
							"value": ["1", "2", "3"]
				  		}
					]
				}
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
	itemsClient.EXPECT().TestRequestUnwrapInnerMessage(gomock.Any(), gomock.Any()).Return(&empty.Empty{}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapInnerMessage(payload: [{names: ["username"], activated: true}])
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
	}).Return(&empty.Empty{}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapInnerMessage(payload: [{names: ["username1", "username2"], activated: true}])
			}
		}`,
	})
}

func TestRequestUnwrapRepeatedMessageResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	itemsClient := mock.NewMockItemsServiceClient(mockCtrl)
	itemsClient.EXPECT().TestRequestUnwrapRepeatedMessage(gomock.Any(), gomock.Any()).Return(&empty.Empty{}, nil).AnyTimes()

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	response := makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapRepeatedMessage(payload: ["username1", "username2"])
			}
		}`,
	})

	tests.AssertJSON(t, `{
		"data": {
			"items": {
				"testRequestUnwrapRepeatedMessage": null
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
	}).Return(&empty.Empty{}, nil).Times(1)

	clients := &mock.Clients{
		ItemsClient: itemsClient,
	}

	makeRequest(t, clients, &handler.RequestOptions{
		Query: `{
			items {
				testRequestUnwrapRepeatedMessage(payload: ["username1", "username2"])
			}
		}`,
	})
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
