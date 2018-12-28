package graphql

import (
	"reflect"

	"github.com/pkg/errors"
)

type DataLoader struct {
	OutputPath string
	Pkg        string
	Loaders    map[string]LoaderModel
}

type LoaderModel struct {
	Service           *Service
	Method            *Method
	InputGoType       GoType
	OutputGoType      GoType
	OutputGraphqlType TypeResolver
	Config            DataLoaderConfig
}

func getServiceByName(files map[string]*TypesFile, name string) *Service {
	for _, file := range files {
		for _, svc := range file.Services {
			if svc.Name == name {
				return &svc
			}
		}
	}

	return nil
}

func createLoaders(config *DataLoadersConfig, files map[string]*TypesFile) (map[string]LoaderModel, error) {
	loaders := make(map[string]LoaderModel)

	for _, dataLoaderConfig := range config.Loaders {
		svc := getServiceByName(files, dataLoaderConfig.ServiceName)

		if svc == nil {
			return nil, errors.Errorf("Failed to found service with name %s for dataloader", dataLoaderConfig.ServiceName)
		}

		method := svc.FindMethodByName(dataLoaderConfig.MethodName)

		if method == nil {
			return nil, errors.Errorf(
				"Failed to found method with name %s in service %s for dataloader",
				dataLoaderConfig.MethodName,
				dataLoaderConfig.ServiceName,
			)
		}

		if len(method.Arguments) != 1 {
			return nil, errors.Errorf(
				"Method `%s` in service `%s` must have 1 input argument",
				dataLoaderConfig.MethodName,
				dataLoaderConfig.ServiceName,
			)
		}

		inputArgument := method.Arguments[0]

		if inputArgument.GoType.Kind != reflect.Slice {
			return nil, errors.Errorf(
				"Argument `%s` in method `%s` in service `%s` must be slice",
				inputArgument.Name,
				dataLoaderConfig.MethodName,
				dataLoaderConfig.ServiceName,
			)
		}

		inputArgumentElementType := inputArgument.GoType.ElemType

		if !inputArgumentElementType.Scalar {
			return nil, errors.Errorf(
				"Argument `%s` in method `%s` in service `%s` must slice of scalars",
				inputArgument.Name,
				dataLoaderConfig.MethodName,
				dataLoaderConfig.ServiceName,
			)
		}

		responseGoType, err := method.DataLoaderResponseType()

		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve output go type")
		}

		loaderMethod := LoaderModel{
			Service:           svc,
			Method:            method,
			InputGoType:       inputArgument.GoType,
			OutputGoType:      responseGoType,
			Config:            dataLoaderConfig,
			OutputGraphqlType: method.GraphQLOutputDataLoaderType,
		}

		loaders[dataLoaderConfig.Name] = loaderMethod
	}

	return loaders, nil
}

func CreateDataLoader(config *DataLoadersConfig, vendorPath string, files map[string]*TypesFile) (*DataLoader, error) {
	if config == nil {
		return nil, nil
	}

	loaders, err := createLoaders(config, files)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create loaders")
	}

	goPkg, err := GoPackageByPath(config.OutputPath, vendorPath)

	if err != nil {
		return nil, errors.New("failed to get go package by path " + goPkg)
	}

	return &DataLoader{
		OutputPath: config.OutputPath,
		Pkg:        goPkg,
		Loaders:    loaders,
	}, nil
}
