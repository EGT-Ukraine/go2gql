package graphql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

const (
	ScalarsPkgPath       = "github.com/EGT-Ukraine/go2gql/api/scalars"
	MultipartFilePkgPath = "github.com/EGT-Ukraine/go2gql/api/multipartfile"
	InterceptorsPkgPath  = "github.com/EGT-Ukraine/go2gql/api/interceptors"
	GraphqlPkgPath       = "github.com/graphql-go/graphql"
	OpentracingPkgPath   = "github.com/opentracing/opentracing-go"
	ErrorsPkgPath        = "github.com/pkg/errors"
	LogPkg               = "github.com/opentracing/opentracing-go/log"
)

type typesGenerator struct {
	File                       *TypesFile
	tracerEnabled              bool
	imports                    *importer.Importer
	outputObjectFieldRenderers []OutputObjectFieldRender
}

func (g typesGenerator) importFunc(importPath string) func() string {
	return func() string {
		return g.imports.New(importPath)
	}
}

func (g typesGenerator) bodyTemplateContext() interface{} {
	return BodyContext{
		File:                 g.File,
		Importer:             g.imports,
		TracerEnabled:        g.tracerEnabled,
		OutputFieldRenderers: g.outputObjectFieldRenderers,
	}

}

func (g typesGenerator) goTypeStr(typ GoType) string {
	return typ.String(g.imports)
}

func (g typesGenerator) goTypeForNew(typ GoType) string {
	switch typ.Kind {
	case reflect.Ptr:
		return g.goTypeStr(*typ.ElemType)
	case reflect.Struct:
		return g.imports.Prefix(typ.Pkg) + typ.Name
	}
	panic("type " + typ.Kind.String() + " is not supported")
}

func (g typesGenerator) bodyTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"ctxPkg":          g.importFunc("context"),
		"debugPkg":        g.importFunc("runtime/debug"),
		"fmtPkg":          g.importFunc("fmt"),
		"logPkg":          g.importFunc(LogPkg),
		"errorsPkg":       g.importFunc(ErrorsPkgPath),
		"gqlPkg":          g.importFunc(GraphqlPkgPath),
		"scalarsPkg":      g.importFunc(ScalarsPkgPath),
		"interceptorsPkg": g.importFunc(InterceptorsPkgPath),
		"opentracingPkg":  g.importFunc(OpentracingPkgPath),
		"concat": func(st ...string) string {
			return strings.Join(st, "")
		},
		"isArray": func(typ GoType) bool {
			return typ.Kind == reflect.Slice
		},
		"goType":       g.goTypeStr,
		"goTypeForNew": g.goTypeForNew,
	}
}

func (g typesGenerator) headTemplateContext() map[string]interface{} {
	return map[string]interface{}{
		"imports": g.imports.Imports(),
		"package": g.File.PackageName,
	}

}

func (g typesGenerator) headTemplateFuncs() map[string]interface{} {
	return nil
}

func (g typesGenerator) generateBody() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesTypes_bodyGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get body template")
	}
	bodyTpl, err := template.New("body").Funcs(g.bodyTemplateFuncs()).Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}
	err = bodyTpl.Execute(buf, g.bodyTemplateContext())
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return buf.Bytes(), nil
}

func (g typesGenerator) generateServicesBody() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesTypes_serviceGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get services body template")
	}
	servicesTpl, err := template.New("service").Funcs(g.bodyTemplateFuncs()).Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	for _, service := range g.File.Services {
		queryServiceContext := ServiceContext{
			Service:        service,
			ServiceMethods: service.QueryMethods,
			FieldType:      "Query",
			TracerEnabled:  g.tracerEnabled,
			BodyContext:    g.bodyTemplateContext().(BodyContext),
		}

		err = servicesTpl.Execute(buf, queryServiceContext)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute template")
		}

		mutationServiceContext := ServiceContext{
			Service:        service,
			ServiceMethods: service.MutationMethods,
			FieldType:      "Mutation",
			TracerEnabled:  g.tracerEnabled,
			BodyContext:    g.bodyTemplateContext().(BodyContext),
		}

		err = servicesTpl.Execute(buf, mutationServiceContext)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute template")
		}
	}

	return buf.Bytes(), nil
}

func (g typesGenerator) generateHead() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesTypes_headGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get head template")
	}
	bodyTpl, err := template.New("head").Funcs(g.headTemplateFuncs()).Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}
	err = bodyTpl.Execute(buf, g.headTemplateContext())
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return buf.Bytes(), nil
}

func (g typesGenerator) generate(out io.Writer) error {
	fieldsRenderer := &fieldsRenderer{
		templateFuncs: g.bodyTemplateFuncs(),
	}

	mapFieldsRenderer := &mapFieldsRenderer{
		templateFuncs: g.bodyTemplateFuncs(),
	}

	g.outputObjectFieldRenderers = append(g.outputObjectFieldRenderers, fieldsRenderer, mapFieldsRenderer)

	body, err := g.generateBody()
	if err != nil {
		return errors.Wrap(err, "failed to generate body")
	}
	servicesBody, err := g.generateServicesBody()
	if err != nil {
		return errors.Wrap(err, "failed to generate services body")
	}
	head, err := g.generateHead()
	if err != nil {
		return errors.Wrap(err, "failed to generate head")
	}
	r := bytes.Join([][]byte{
		head,
		body,
		servicesBody,
	}, nil)

	res, err := imports.Process("file", r, &imports.Options{
		Comments: true,
	})
	// TODO: fix this
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error(), string(r))
	} else {
		r = res
	}
	_, err = out.Write(r)
	if err != nil {
		return errors.Wrap(err, "failed to write  output")
	}

	return nil
}
