package dataloader

import (
	"time"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

type DataLoader struct {
	OutputPath string
	Pkg        string
	Loaders    map[string]LoaderModel
}

type Service struct {
	Name          string
	CallInterface graphql.GoType
}

type LoaderModel struct {
	Name                  string
	WaitDuration          time.Duration
	Service               *Service
	Method                *graphql.Method
	InputGoType           graphql.GoType
	OutputGoType          graphql.GoType
	OutputGraphqlType     graphql.TypeResolver
	OutputGraphqlTypeName string
	FetchCode             func(importer *importer.Importer) string
	Slice                 bool
}

func (p *Plugin) createDataLoader(config *DataLoadersConfig, vendorPath string) (*DataLoader, error) {
	if config == nil {
		return nil, nil
	}

	goPkg, err := graphql.GoPackageByPath(config.OutputPath, vendorPath)

	if err != nil {
		return nil, errors.New("failed to get go package by path " + goPkg)
	}

	return &DataLoader{
		OutputPath: config.OutputPath,
		Pkg:        goPkg,
		Loaders:    p.loaders,
	}, nil
}
