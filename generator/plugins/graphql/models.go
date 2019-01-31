package graphql

import (
	"reflect"
	"strings"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

const (
	SchemaNodeTypeObject  = "OBJECT"
	SchemaNodeTypeService = "SERVICE"
)

type TypeResolver func(ctx BodyContext) string

type ValueResolver func(arg string, ctx BodyContext) string

type AssigningWrapper func(arg string, ctx BodyContext) string

type PayloadErrorChecker func(arg string) string

type PayloadErrorAccessor func(arg string) string

type ClientMethodCaller func(client string, req string, ctx BodyContext) string

const KindBytes = reflect.Kind(255)

type GoType struct {
	Scalar    bool
	Kind      reflect.Kind
	ElemType  *GoType
	Elem2Type *GoType
	Name      string
	Pkg       string
}

func (g GoType) String(i *importer.Importer) string {
	if typeIsScalar(g) && g.Name == "" {
		return g.Kind.String()
	}
	switch g.Kind {
	case reflect.Slice:
		return "[]" + g.ElemType.String(i)
	case reflect.Ptr:
		return "*" + g.ElemType.String(i)
	case reflect.Struct, reflect.Interface:
		return i.Prefix(g.Pkg) + g.Name
	case reflect.Map:
		return "map[" + g.ElemType.String(i) + "]" + g.Elem2Type.String(i)
	case KindBytes:
		return "[]byte"
	}
	if g.Name != "" {
		return i.Prefix(g.Pkg) + g.Name
	}
	panic("type " + g.Kind.String() + " is not supported")
}

type InputObjectResolver struct {
	FunctionName string
	OutputGoType GoType
	Fields       []InputObjectResolverField
	OneOfFields  []InputObjectResolverOneOf
}

type InputObjectResolverOneOf struct {
	OutputFieldName string
	Fields          []InputObjectResolverOneOfField
}

type InputObjectResolverOneOfField struct {
	GraphQLInputFieldName string
	ValueResolver         ValueResolver
	ResolverWithError     bool
	AssigningWrapper      AssigningWrapper
}

type InputObjectResolverField struct {
	OutputFieldName       string
	GraphQLInputFieldName string
	GoType                GoType
	ValueResolver         ValueResolver
	ResolverWithError     bool
	IsFromArgs            bool
}

type InputObject struct {
	VariableName string
	GraphQLName  string
	Fields       []ObjectField
}

type ObjectField struct {
	Name          string
	Type          TypeResolver
	GoType        GoType
	QuotedComment string
	Value         ValueResolver
	NeedCast      bool
	CastTo        GoType
}

type DataLoaderField struct {
	Name                         string
	Type                         string
	IsSlice                      bool
	ParentKeyFieldName           string
	KeyFieldSlice                bool
	NormalizedParentKeyFieldName string
	DataLoaderName               string
}

type OutputObject struct {
	VariableName     string
	GraphQLName      string
	GoType           GoType
	Fields           []ObjectField
	DataLoaderFields []*DataLoaderField // TODO: move to dataloader plugin
	MapFields        []ObjectField
}

func (s *OutputObject) FindFieldByName(name string) *ObjectField {
	searchName := strings.ToLower(name)

	for _, field := range s.Fields {
		if strings.ToLower(field.Name) == searchName {
			return &field
		}
	}

	return nil
}

type Enum struct {
	VariableName string
	GraphQLName  string
	Comment      string
	Values       []EnumValue
}

type EnumValue struct {
	Name    string
	Value   int
	Comment string
}

type MapInputObject struct {
	VariableName    string
	GraphQLName     string
	KeyObjectType   TypeResolver
	ValueObjectType TypeResolver
}

type MapInputObjectResolver struct {
	FunctionName           string
	KeyGoType              GoType
	ValueGoType            GoType
	KeyResolver            ValueResolver
	KeyResolverWithError   bool
	ValueResolver          ValueResolver
	ValueResolverWithError bool
}

type MapOutputObject struct {
	VariableName    string
	GraphQLName     string
	KeyObjectType   TypeResolver
	ValueObjectType TypeResolver
}

type Service struct {
	Name            string
	QuotedComment   string
	CallInterface   GoType
	QueryMethods    []Method
	MutationMethods []Method
}

type Method struct {
	Name                   string
	QuotedComment          string
	GraphQLOutputType      TypeResolver
	Arguments              []MethodArgument
	RequestResolver        ValueResolver
	RequestResolverWithErr bool
	ClientMethodCaller     ClientMethodCaller
	RequestType            GoType
	PayloadErrorChecker    PayloadErrorChecker
	PayloadErrorAccessor   PayloadErrorAccessor
}

type MethodArgument struct {
	Name          string
	Type          TypeResolver
	QuotedComment string
}

type TypesFile struct {
	PackageName             string
	Package                 string
	Enums                   []Enum
	OutputObjects           []OutputObject
	InputObjects            []InputObject
	InputObjectResolvers    []InputObjectResolver
	MapInputObjects         []MapInputObject
	MapInputObjectResolvers []MapInputObjectResolver
	MapOutputObjects        []MapOutputObject
	Services                []Service
}

type BodyContext struct {
	File                 *TypesFile
	Importer             *importer.Importer
	TracerEnabled        bool
	OutputFieldRenderers []OutputObjectFieldRender
}

type ServiceContext struct {
	Service        Service
	TracerEnabled  bool
	ServiceMethods []Method
	FieldType      string
	BodyContext    BodyContext
}

type SchemaBodyContext struct {
	File           SchemaConfig
	Importer       *importer.Importer
	SchemaName     string
	Services       []SchemaService
	QueryObject    string
	MutationObject string
	Objects        []*gqlObject
	TracerEnabled  bool
}

type SchemaService struct {
	Name            string
	ConstructorName string
	Fields          []string
	Pkg             string
	ClientGoType    GoType
}

type fieldConfig struct {
	QuotedComment string
	Name          string
	Service       *SchemaService
	Object        *gqlObject
}

type gqlObject struct {
	QueryObject   bool
	Name          string
	QuotedComment string
	Fields        []fieldConfig
}

func (gqlObject *gqlObject) TypeName() string {
	if gqlObject.QueryObject {
		return "Query"
	}

	return "Mutation"
}
