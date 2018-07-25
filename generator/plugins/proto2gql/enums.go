package proto2gql

import (
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) enumTypeResolver(enumFile *parsedFile, enum *parser.Enum) (graphql.TypeResolver, error) {
	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(enumFile.OutputPkg) + g.enumVariable(enumFile, enum)
	}, nil
}

func (g *Proto2GraphQL) enumGraphQLName(enumFile *parsedFile, enum *parser.Enum) string {
	return enumFile.Config.GetGQLEnumsPrefix() + camelCaseSlice(enum.TypeName)
}

func (g *Proto2GraphQL) enumVariable(enumFile *parsedFile, enum *parser.Enum) string {
	return enumFile.Config.GetGQLEnumsPrefix() + camelCaseSlice(enum.TypeName)
}

func (g *Proto2GraphQL) prepareFileEnums(file *parsedFile) ([]graphql.Enum, error) {
	var res []graphql.Enum
	for _, enum := range file.File.Enums {
		vals := make([]graphql.EnumValue, len(enum.Values))
		for i, value := range enum.Values {
			vals[i] = graphql.EnumValue{
				Name:    value.Name,
				Value:   value.Value,
				Comment: value.QuotedComment,
			}
		}
		res = append(res, graphql.Enum{
			VariableName: g.enumVariable(file, enum),
			GraphQLName:  g.enumGraphQLName(file, enum),
			Comment:      enum.QuotedComment,
			Values:       vals,
		})
	}
	return res, nil
}
