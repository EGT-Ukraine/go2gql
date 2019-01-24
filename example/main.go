package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/graphql-go/graphql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/EGT-Ukraine/go2gql/example/out/schema"
	example "github.com/EGT-Ukraine/go2gql/example/proto"
)

type Client struct {
}

func (Client) GetQueryMethod(ctx context.Context, in *example.AOneOffs, opts ...grpc.CallOption) (*example.B, error) {
	return &example.B{
		RScalar: []int32{
			1, 2, 3, 4, 5,
		},
	}, nil
}

func (Client) MutationMethod(ctx context.Context, in *example.B, opts ...grpc.CallOption) (*example.A, error) {
	return &example.A{
		NREnum: example.A_Val5,
	}, nil
}

func (Client) QueryMethod(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*timestamp.Timestamp, error) {
	return &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Now().Nanosecond()),
	}, nil
}

func (Client) GetMutatuionMethod(ctx context.Context, in *example.MsgWithEmpty, opts ...grpc.CallOption) (*example.MsgWithEmpty, error) {
	return &example.MsgWithEmpty{}, nil
}

func (Client) GetEmptiesMsg(ctx context.Context, in *example.Empty, opts ...grpc.CallOption) (*example.Empty, error) {
	return &example.Empty{}, nil
}

func main() {
	schem, err := schema.GetExampleSchemaSchema(schema.ExampleSchemaSchemaClients{
		ServiceExampleClient: Client{},
	}, nil)
	if err != nil {
		panic(err)
	}
	spew.Dump(graphql.Do(graphql.Params{
		Schema: schem,
		RequestString: `
				{
					getQueryMethod{
						r_scalar
					}
					queryMethod{
						seconds
						nanos
					}
					getEmptiesMsg
				}
			`,
	}))
	spew.Dump(graphql.Do(graphql.Params{
		Schema: schem,
		RequestString: `
				mutation {
					mutationMethod{
						n_r_enum
					}
					getMutatuionMethod{
						empty_field
					}
				}
			`,
	}))
}
