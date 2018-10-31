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
	types         map[string]*TypesFile
	imports       *importer.Importer
}

func (g schemaGenerator) importFunc(importPath string) func() string {
	return func() string {
		return g.imports.New(importPath)
	}
}

func (g *schemaGenerator) resolveObjectFields(nodeCfg SchemaNodeConfig, object *gqlObject) (services []SchemaService, newObjectx []*gqlObject, err error) {
	var newObjs []*gqlObject
	var newServices []SchemaService
	switch nodeCfg.Type {
	case SchemaNodeTypeObject:
		for _, fld := range nodeCfg.Fields {
			fldObj := &gqlObject{
				QueryObject: object.QueryObject,
				Name:        strings.Replace(fld.ObjectName, " ", "_", -1),
			}
			services, subObjs, err := g.resolveObjectFields(fld, fldObj)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "can't resolve field %s object fields", fld.Field)
			}
			if len(fldObj.Fields) > 0 {
				object.Fields = append(object.Fields, fieldConfig{
					Name:   fld.Field,
					Object: fldObj,
				})
				newServices = append(newServices, services...)
				newObjs = append(newObjs, subObjs...)
				newObjs = append(newObjs, fldObj)
			}
		}
		return newServices, newObjs, nil
	case SchemaNodeTypeService:
		for _, typesFile := range g.types {
			for _, service := range typesFile.Services {
				if service.Name != nodeCfg.Service {
					continue
				}

				var serviceMethods []Method

				if object.QueryObject == true {
					serviceMethods = service.QueryMethods
				} else {
					serviceMethods = service.MutationMethods
				}

				meths := make([]string, len(serviceMethods))

				for i, meth := range serviceMethods {
					meths[i] = meth.Name
				}
				fields := g.filterMethods(meths, nodeCfg.FilterMethods, nodeCfg.ExcludeMethods)
				srv := SchemaService{
					Name:         service.Name,
					ClientGoType: service.CallInterface,
					Pkg:          typesFile.Package,
				}
				newServices = append(newServices, srv)
				for _, fld := range fields {
					object.Fields = append(object.Fields, fieldConfig{
						Name:    fld,
						Service: &srv,
					})

				}
				return newServices, nil, nil
			}
		}
		return nil, nil, errors.Errorf("service '%s' not found", nodeCfg.Service)

	default:
		return nil, nil, errors.Errorf("unknown type %s", nodeCfg.Type)
	}
}
func (g *schemaGenerator) filterMethods(methods []string, filter, exclude []string) []string {
	var res []string
	var filteredMethods = make(map[string]struct{})
	for _, f := range filter {
		filteredMethods[f] = struct{}{}
	}
	var excludedMethods = make(map[string]interface{})
	for _, f := range exclude {
		excludedMethods[f] = struct{}{}
	}
	for _, m := range methods {
		if len(excludedMethods) > 0 {
			if _, ok := excludedMethods[m]; ok {
				continue
			}
		}
		if len(filteredMethods) > 0 {
			if _, ok := filteredMethods[m]; !ok {
				continue
			}
		}
		res = append(res, m)
	}
	return res
}
func (g *schemaGenerator) resolveObjectsToGenerate() ([]SchemaService, []*gqlObject, string, string, error) {
	var objects []*gqlObject
	var services []SchemaService
	var queryObject, mutationObject string
	if g.schemaCfg.Queries != nil {
		var queryObj = &gqlObject{
			QueryObject: true,
			Name:        "Query",
		}
		newServices, newObjs, err := g.resolveObjectFields(*g.schemaCfg.Queries, queryObj)
		if err != nil {
			return nil, nil, "", "", errors.Wrap(err, "failed to resolve queries fields")
		}
		objects = append(objects, newObjs...)
		objects = append(objects, queryObj)
		services = append(services, newServices...)
		queryObject = queryObj.Name
	}
	if g.schemaCfg.Mutations != nil {
		var mutationObj = &gqlObject{
			QueryObject: false,
			Name:        "Mutation",
		}
		newServices, newObjs, err := g.resolveObjectFields(*g.schemaCfg.Mutations, mutationObj)

		if err != nil {
			return nil, nil, "", "", errors.Wrap(err, "failed to resolve mutations fields")
		}
		if len(mutationObj.Fields) > 0 {
			objects = append(objects, newObjs...)
			objects = append(objects, mutationObj)
			services = append(services, newServices...)
			mutationObject = mutationObj.Name
		}
	}
	var uniqueServices []SchemaService
	var handledService = map[string]struct{}{}
	for _, service := range services {
		if _, ok := handledService[service.Name]; !ok {
			uniqueServices = append(uniqueServices, service)
			handledService[service.Name] = struct{}{}
		}
	}
	return uniqueServices, objects, queryObject, mutationObject, nil
}
func (g schemaGenerator) bodyTemplateContext() (interface{}, error) {
	services, objects, queryObj, mutationsObj, err := g.resolveObjectsToGenerate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve objects to generate")
	}
	return SchemaBodyContext{
		File:           g.schemaCfg,
		Importer:       g.imports,
		SchemaName:     g.schemaCfg.Name,
		QueryObject:    queryObj,
		MutationObject: mutationsObj,
		Objects:        objects,
		Services:       services,
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
		g := schemaGenerator{
			types:         p.files,
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
