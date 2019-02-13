package loaders

import (
	"context"

	proto "github.com/EGT-Ukraine/go2gql/example/proto"
)

type LoaderClients interface {
	GetServiceExampleClient() proto.ServiceExampleClient
}

type DataLoaders struct {
	SomeEntitiesByIDLoader  ListSomeEntitiesResponse_SomeEntityLoader
	SomeEntitiesByAIDLoader ListSomeEntitiesResponse_SomeEntitySliceLoader
}

type dataLoadersContextKeyType struct{}

var dataLoadersContextKey = dataLoadersContextKeyType{}

func GetContextWithLoaders(ctx context.Context, apiClients LoaderClients) context.Context {
	dataLoaders := &DataLoaders{
		SomeEntitiesByIDLoader:  createSomeEntitiesByID(ctx, apiClients.GetServiceExampleClient()),
		SomeEntitiesByAIDLoader: createSomeEntitiesByAID(ctx, apiClients.GetServiceExampleClient()),
	}

	return context.WithValue(ctx, dataLoadersContextKey, dataLoaders)
}

func createSomeEntitiesByID(ctx context.Context, client proto.ServiceExampleClient) ListSomeEntitiesResponse_SomeEntityLoader {
	return ListSomeEntitiesResponse_SomeEntityLoader{
		fetch: func(keys []string) ([]*proto.ListSomeEntitiesResponse_SomeEntity, []error) {
			response, err := client.ListSomeEntities(ctx, &proto.ListSomeEntitiesRequest{Filter: &proto.ListSomeEntitiesRequest_Filter{Ids: keys}})
			if err != nil {
				return nil, []error{err}
			}
			var result = make([]*proto.ListSomeEntitiesResponse_SomeEntity, len(keys))
			for i, key := range keys {
				for _, value := range response.Entities {
					if value.Id == key {
						result[i] = value
						break
					}
				}
			}

			return result, nil
		},
		wait: 10000000,
	}
}
func createSomeEntitiesByAID(ctx context.Context, client proto.ServiceExampleClient) ListSomeEntitiesResponse_SomeEntitySliceLoader {
	return ListSomeEntitiesResponse_SomeEntitySliceLoader{
		fetch: func(keys []string) ([][]*proto.ListSomeEntitiesResponse_SomeEntity, []error) {
			response, err := client.ListSomeEntities(ctx, &proto.ListSomeEntitiesRequest{Filter: &proto.ListSomeEntitiesRequest_Filter{AIds: keys}})
			if err != nil {
				return nil, []error{err}
			}
			var result = make([][]*proto.ListSomeEntitiesResponse_SomeEntity, len(keys))
			for i, key := range keys {
				for _, value := range response.Entities {
					if value.AId == key {
						result[i] = append(result[i], value)
					}
				}
			}

			return result, nil
		},
		wait: 10000000,
	}
}
func GetDataLoadersFromContext(ctx context.Context) *DataLoaders {
	val := ctx.Value(dataLoadersContextKey)

	if val == nil {
		return nil
	}

	return val.(*DataLoaders)
}
