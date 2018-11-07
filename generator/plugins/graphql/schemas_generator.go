package graphql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

type schemaGenerator struct {
	tracerEnabled bool
	schemaCfg     SchemaConfig
	goPkg         string
	parser        *schemaParser
	imports       *importer.Importer
}

func (g schemaGenerator) importFunc(importPath string) func() string {
	return func() string {
		return g.imports.New(importPath)
	}
}

func (g schemaGenerator) bodyTemplateContext() (interface{}, error) {
	schemaObjects, err := g.parser.SchemaObjects()

	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve objects to generate")
	}
	return SchemaBodyContext{
		File:           g.schemaCfg,
		Importer:       g.imports,
		SchemaName:     g.schemaCfg.Name,
		QueryObject:    schemaObjects.QueryObject,
		MutationObject: schemaObjects.MutationObject,
		Objects:        schemaObjects.Objects,
		Services:       schemaObjects.Services,
		TracerEnabled:  g.tracerEnabled,
	}, nil

}
func (g schemaGenerator) goTypeStr(typ GoType) string {
	return typ.String(g.imports)
}

func (g schemaGenerator) goTypeForNew(typ GoType) string {
	switch typ.Kind {
	case reflect.Ptr:
		return g.goTypeStr(*typ.ElemType)
	case reflect.Struct:
		return g.imports.Prefix(typ.Pkg) + typ.Name
	}
	panic("type " + typ.Kind.String() + " is not supported")
}

func (g schemaGenerator) bodyTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"ctxPkg":          g.importFunc("context"),
		"debugPkg":        g.importFunc("runtime/debug"),
		"fmtPkg":          g.importFunc("fmt"),
		"errorsPkg":       g.importFunc(ErrorsPkgPath),
		"gqlPkg":          g.importFunc(GraphqlPkgPath),
		"scalarsPkg":      g.importFunc(ScalarsPkgPath),
		"interceptorsPkg": g.importFunc(InterceptorsPkgPath),
		"opentracingPkg":  g.importFunc(OpentracingPkgPath),
		"tracerPkg":       g.importFunc(TracerPkgPath),
		"concat": func(st ...string) string {
			return strings.Join(st, "")
		},
		"isArray": func(typ GoType) bool {
			return typ.Kind == reflect.Slice
		},
		"goType":       g.goTypeStr,
		"goTypeForNew": g.goTypeForNew,

		"serviceConstructor": func(filedType string, service SchemaService, ctx SchemaBodyContext) string {
			return ctx.Importer.Prefix(service.Pkg) + "Get" + service.Name + "Service" + filedType + "Methods"
		},
	}
}

func (g schemaGenerator) headTemplateContext() map[string]interface{} {
	return map[string]interface{}{
		"imports": g.imports.Imports(),
		"package": g.schemaCfg.OutputPackage,
	}

}
func (g schemaGenerator) headTemplateFuncs() map[string]interface{} {
	return nil
}
func (g schemaGenerator) generateBody() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesSchemas_bodyGohtmlBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get head template")
	}
	bodyTpl, err := template.New("body").Funcs(g.bodyTemplateFuncs()).Parse(string(tmpl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}
	bodyCtx, err := g.bodyTemplateContext()
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare body context")
	}
	err = bodyTpl.Execute(buf, bodyCtx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}
	return buf.Bytes(), nil
}

func (g schemaGenerator) generateHead() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesSchemas_headGohtmlBytes()
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
func (g schemaGenerator) generate(out io.Writer) error {
	body, err := g.generateBody()
	if err != nil {
		return errors.Wrap(err, "failed to generate body")
	}
	head, err := g.generateHead()
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
	// TODO: fix this
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

func (p *Plugin) generateSchemas() error {
	for _, schema := range p.schemaConfigs {
		pkg, err := GoPackageByPath(filepath.Dir(schema.OutputPath), p.generateCfg.VendorPath)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve schema %s output go package", schema.Name)
		}

		parser := NewSchemaParser(schema, p.files)

		g := schemaGenerator{
			parser:        parser,
			tracerEnabled: p.generateCfg.GenerateTraces,
			schemaCfg:     schema,
			goPkg:         pkg,
			imports: &importer.Importer{
				CurrentPackage: pkg,
			},
		}
		file, err := os.OpenFile(schema.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open schema %s output file for write", schema.OutputPath)
		}
		err = g.generate(file)
		if err != nil {
			if cerr := file.Close(); cerr != nil {
				err = errors.Wrap(err, cerr.Error())
			}
			return errors.Wrapf(err, "failed to generate types file %s", schema.OutputPath)
		}
		if file.Close(); err != nil {
			return errors.Wrapf(err, "failed to close generated schema %s file", schema.OutputPath)
		}
	}
	return nil
}
