package graphql

import (
	"os"
	"reflect"
	"text/template"

	"github.com/EGT-Ukraine/dataloaden/pkg/generator"
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

const DefaultWaitDurationMs = 10

type LoadersContext struct {
	Imports []importer.Import
	Loaders []Loader
}

type Loader struct {
	LoaderTypeName string
	Service        Service
	FetchCode      string
	RequestGoType  GoType
	ResponseGoType GoType
	Config         DataLoaderConfig
}

type LoaderGenerator struct {
	dataLoader *DataLoader
}

func NewLoaderGenerator(dataLoader *DataLoader) *LoaderGenerator {
	return &LoaderGenerator{dataLoader: dataLoader}
}

func (p *LoaderGenerator) GenerateDataLoaders() error {
	if err := os.MkdirAll(p.dataLoader.OutputPath, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create output path dir "+p.dataLoader.OutputPath)
	}

	if err := p.generateSchemaLoaders(); err != nil {
		return err
	}

	for _, dataLoader := range p.dataLoader.Loaders {
		if err := p.generateLoaders(dataLoader.InputGoType, dataLoader.OutputGoType); err != nil {
			return err
		}
	}

	return nil
}

func (p *LoaderGenerator) generateLoaders(requestGoType GoType, responseGoType GoType) error {
	keyType := requestGoType.ElemType.Kind.String()

	var typeName string

	slice := responseGoType.Kind == reflect.Slice

	if slice {
		typeName = responseGoType.ElemType.ElemType.Pkg + "." + responseGoType.ElemType.ElemType.Name
	} else {
		typeName = responseGoType.ElemType.Pkg + "." + responseGoType.ElemType.Name
	}

	if err := generator.Generate(typeName, keyType, slice, true, p.dataLoader.OutputPath); err != nil {
		return errors.Wrapf(err, "Failed to generate loader for '%s'", typeName)
	}

	return nil
}

func (p *LoaderGenerator) generateSchemaLoaders() error {
	path := p.dataLoader.OutputPath + "/loaders.go"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		return errors.Wrapf(err, "failed to open loaders %s output file for write", path)
	}

	err = p.renderLoaders(file)

	if err != nil {
		if cerr := file.Close(); cerr != nil {
			err = errors.Wrap(err, cerr.Error())
		}

		return errors.Wrapf(err, "failed to generate loaders file %s", path)
	}

	if file.Close(); err != nil {
		return errors.Wrapf(err, "failed to close generated loaders %s file", path)
	}

	return nil
}

func (p *LoaderGenerator) renderLoaders(out *os.File) error {
	tmpl, err := templatesLoadersGohtmlBytes()
	if err != nil {
		return errors.Wrap(err, "failed to get loaders template")
	}

	fileImporter := &importer.Importer{}

	templateFuncs := map[string]interface{}{
		"goType": func(typ GoType) string {
			return typ.String(fileImporter)
		},
		"duration": func(duration int) int {
			if duration == 0 {
				return DefaultWaitDurationMs
			}

			return duration
		},
	}

	servicesTpl, err := template.New("config").Funcs(templateFuncs).Parse(string(tmpl))
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	var loaders []Loader

	for _, dataLoaderModel := range p.dataLoader.Loaders {
		service := dataLoaderModel.Service

		fileImporter.New(service.CallInterface.Pkg)

		method := dataLoaderModel.Method

		requestGoType := dataLoaderModel.InputGoType

		responseGoType := dataLoaderModel.OutputGoType

		fileImporter.New(requestGoType.Pkg)

		loaderTypeName := responseGoType.ElemType.Name

		if responseGoType.Kind == reflect.Slice {
			loaderTypeName = responseGoType.ElemType.ElemType.Name + "Slice"
			fileImporter.New(responseGoType.ElemType.ElemType.Pkg)
		} else {
			fileImporter.New(responseGoType.Pkg)
		}

		loaders = append(loaders, Loader{
			LoaderTypeName: loaderTypeName,
			Service:        *service,
			FetchCode:      method.DataLoaderFetch(fileImporter),
			RequestGoType:  requestGoType,
			ResponseGoType: responseGoType,
			Config:         dataLoaderModel.Config,
		})
	}

	configContext := LoadersContext{
		Imports: fileImporter.Imports(),
		Loaders: loaders,
	}

	err = servicesTpl.Execute(out, configContext)
	if err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}
