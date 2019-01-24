package swagger2gql

import (
	"go/build"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

const (
	strFmtPkg = "github.com/go-openapi/strfmt"
	timePkg   = "time"
)

var scalarsGoTypesNames = map[parser.Kind]string{
	parser.KindString:  "string",
	parser.KindFloat32: "float32",
	parser.KindFloat64: "float64",
	parser.KindInt64:   "int64",
	parser.KindInt32:   "int32",
	parser.KindBoolean: "bool",
	parser.KindFile:    "File",
}
var scalarsGoTypes = map[parser.Kind]graphql.GoType{
	parser.KindBoolean: {Scalar: true, Kind: reflect.Bool},
	parser.KindFloat64: {Scalar: true, Kind: reflect.Float64},
	parser.KindFloat32: {Scalar: true, Kind: reflect.Float32},
	parser.KindInt64:   {Scalar: true, Kind: reflect.Int64},
	parser.KindInt32:   {Scalar: true, Kind: reflect.Int32},
	parser.KindString:  {Scalar: true, Kind: reflect.String},
}

func GoPackageByPath(path, vendorPath string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve absolute filepath")
	}
	var prefixes []string
	if vendorPath != "" {
		absVendorPath, err := filepath.Abs(vendorPath)
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve absolute vendor path")
		}
		prefixes = append(prefixes, absVendorPath)
	}
	absGoPath, err := filepath.Abs(build.Default.GOPATH)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve absolute gopath")
	}
	prefixes = append(prefixes, filepath.Join(absGoPath, "src"))

	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return strings.TrimLeft(strings.TrimPrefix(path, prefix), " "+string(filepath.Separator)), nil
		}
	}
	return "", errors.Errorf("path '%s' is outside GOPATH or Vendor folder", path)
}
func (p *Plugin) goTypeByParserType(typeFile *parsedFile, typ parser.Type, ptrObj bool) (_ graphql.GoType, err error) {
	if typ == parser.ObjDateTime {
		t := graphql.GoType{
			Kind: reflect.Struct,
			Name: "DateTime",
			Pkg:  strFmtPkg,
		}
		if ptrObj {
			tcp := t
			t = graphql.GoType{
				Kind:     reflect.Ptr,
				ElemType: &tcp,
			}
		}
		return t, nil
	}
	switch t := typ.(type) {
	case *parser.Scalar:
		goTyp, ok := scalarsGoTypes[t.Kind()]
		if !ok {
			err = errors.Errorf("convertation of scalar %s to golang type is not implemented", typ.Kind())
			return
		}
		return goTyp, nil
	case *parser.Object:
		if ptrObj {
			return graphql.GoType{
				Kind: reflect.Ptr,
				ElemType: &graphql.GoType{
					Kind: reflect.Struct,
					Name: pascalize(camelCaseSlice(t.Route)),
					Pkg:  typeFile.Config.ModelsGoPath,
				},
			}, nil
		}
		return graphql.GoType{
			Kind: reflect.Struct,
			Name: pascalize(camelCaseSlice(t.Route)),
			Pkg:  typeFile.Config.ModelsGoPath,
		}, nil
	case *parser.Array:
		elemGoType, err := p.goTypeByParserType(typeFile, t.ElemType, ptrObj)
		if err != nil {
			err = errors.Wrap(err, "failed to resolve array element go type")
			return graphql.GoType{}, err
		}
		return graphql.GoType{
			Kind:     reflect.Slice,
			ElemType: &elemGoType,
		}, nil
	case *parser.Map:
		valueType, err := p.goTypeByParserType(typeFile, t.ElemType, true)
		if err != nil {
			return graphql.GoType{}, errors.Wrap(err, "failed to resolve map output type")
		}
		return graphql.GoType{
			Kind: reflect.Map,
			ElemType: &graphql.GoType{
				Kind: reflect.String,
			},
			Elem2Type: &valueType,
		}, nil
	}
	err = errors.Errorf("unknown type %v", typ.Kind().String())
	return
}

// camelCaseSlice is like camelCase, but the argument is a slice of strings to
// be joined with "_".
func camelCaseSlice(elem []string) string      { return pascalize(strings.Join(elem, "")) }
func snakeCamelCaseSlice(elem []string) string { return pascalize(strings.Join(elem, "_")) }
func pascalize(arg string) string {
	arg = strings.NewReplacer(">=", "Ge", "<=", "Le", ">", "Gt", "<", "Lt", "=", "Eq").Replace(arg)
	if len(arg) == 0 || arg[0] > '9' {
		return swag.ToGoName(arg)
	}
	if arg[0] == '+' {
		return swag.ToGoName("Plus " + arg[1:])
	}
	if arg[0] == '-' {
		return swag.ToGoName("Minus " + arg[1:])
	}

	return swag.ToGoName("Nr " + arg)
}
