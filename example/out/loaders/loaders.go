package loaders

import (
	"context"

	proto "github.com/EGT-Ukraine/go2gql/example/proto"
)

type LoaderClients interface {
	GetServiceExampleClient() proto.ServiceExampleClient
}

type DataLoaders struct {
	SomeEntitiesByAIDLoader ListSomeEntitiesResponse_SomeEntitySliceLoader
	SomeEntitiesByIDLoader  ListSomeEntitiesResponse_SomeEntityLoader
}

type dataLoadersContextKeyType struct{}

var dataLoadersContextKey = dataLoadersContextKeyType{}

func GetContextWithLoaders(ctx context.Context, apiClients LoaderClients) context.Context {
	dataLoaders := &DataLoaders{
		SomeEntitiesByAIDLoader: createSomeEntitiesByAID(ctx, apiClients.GetServiceExampleClient()),
		SomeEntitiesByIDLoader:  createSomeEntitiesByID(ctx, apiClients.GetServiceExampleClient()),
	}

	return context.WithValue(ctx, dataLoadersContextKey, dataLoaders)
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
		wait: 10,
	}
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
		wait: 10,
	}
}
func GetDataLoadersFromContext(ctx context.Context) *DataLoaders {
	val := ctx.Value(dataLoadersContextKey)

	if val == nil {
		return nil
	}

	return val.(*DataLoaders)
}
