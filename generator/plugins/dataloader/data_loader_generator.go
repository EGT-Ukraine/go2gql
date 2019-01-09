package dataloader

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/EGT-Ukraine/dataloaden/pkg/generator"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

const DefaultWaitDurationMs = 10

type LoadersHeadContext struct {
	Imports []importer.Import
}

type LoadersBodyContext struct {
	Loaders []Loader
}

type Loader struct {
	LoaderTypeName    string
	Service           Service
	FetchCode         string
	RequestGoType     graphql.GoType
	ResponseGoType    graphql.GoType
	OutputGraphqlType graphql.TypeResolver
	Config            DataLoaderProviderConfig
}

type LoaderGenerator struct {
	dataLoader *DataLoader
	importer   *importer.Importer
}

func NewLoaderGenerator(dataLoader *DataLoader) *LoaderGenerator {
	return &LoaderGenerator{dataLoader: dataLoader, importer: &importer.Importer{}}
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

func (p *LoaderGenerator) generateLoaders(requestGoType graphql.GoType, responseGoType graphql.GoType) error {
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

func (p *LoaderGenerator) generateBody() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := templatesLoaders_bodyGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loaders template")
	}

	importFunc := func(importPath string) func() string {
		return func() string {
			return p.importer.New(importPath)
		}
	}

	templateFuncs := map[string]interface{}{
		"timePkg": importFunc("time"),
		"goType": func(typ graphql.GoType) string {
			return typ.String(p.importer)
		},
		"duration": func(duration int) int {
			if duration == 0 {
				return DefaultWaitDurationMs
			}

			return duration
		},
	}

	servicesTpl, err := template.New("loaders_body").Funcs(templateFuncs).Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	var loaders []Loader

	for _, dataLoaderModel := range p.dataLoader.Loaders {
		service := dataLoaderModel.Service

		requestGoType := dataLoaderModel.InputGoType

		responseGoType := dataLoaderModel.OutputGoType

		loaderTypeName := responseGoType.ElemType.Name

		if responseGoType.Kind == reflect.Slice {
			loaderTypeName = responseGoType.ElemType.ElemType.Name + "Slice"
		}

		loaders = append(loaders, Loader{
			LoaderTypeName: loaderTypeName,
			Service:        *service,
			FetchCode:      dataLoaderModel.FetchCode(p.importer),
			RequestGoType:  requestGoType,
			ResponseGoType: responseGoType,
			Config:         dataLoaderModel.Config,
		})
	}

	context := LoadersBodyContext{
		Loaders: loaders,
	}

	err = servicesTpl.Execute(buf, context)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return buf.Bytes(), nil
}

func (p *LoaderGenerator) generateHead() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesLoaders_headGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get head template")
	}
	bodyTpl, err := template.New("loaders_head").Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	context := LoadersHeadContext{
		Imports: p.importer.Imports(),
	}

	err = bodyTpl.Execute(buf, context)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return buf.Bytes(), nil
}

func (p *LoaderGenerator) renderLoaders(out *os.File) error {
	body, err := p.generateBody()
	if err != nil {
		return errors.Wrap(err, "failed to generate body")
	}
	head, err := p.generateHead()
	if err != nil {
		return errors.Wrap(err, "failed to generate head")
	}
	r := bytes.Join([][]byte{
		head,
		body,
	}, nil)

	res, err := imports.Process("file", r, &imports.Options{
		Comments: true,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		r = res
	}
	_, err = out.Write(r)
	if err != nil {
		return errors.Wrap(err, "failed to write  output")
	}

	return nil
}
