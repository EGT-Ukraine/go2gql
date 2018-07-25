package graphql

func GqlStringTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(GraphqlPkgPath) + ".String"
}
func GqlIntTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(GraphqlPkgPath) + ".Int"
}
func GqlFloatTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(GraphqlPkgPath) + ".Float"
}
func GqlBoolTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(GraphqlPkgPath) + ".Boolean"
}
func GqlDateTimeTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(GraphqlPkgPath) + ".DateTime"
}
func GqlListTypeResolver(r TypeResolver) TypeResolver {
	return func(ctx BodyContext) string {
		return ctx.Importer.New(GraphqlPkgPath) + ".NewList(" + r(ctx) + ")"
	}

}
func GqlNonNullTypeResolver(r TypeResolver) TypeResolver {
	return func(ctx BodyContext) string {
		return ctx.Importer.New(GraphqlPkgPath) + ".NewNonNull(" + r(ctx) + ")"
	}
}
func GqlNoDataTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".NoDataScalar"
}
func GqlInt64TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLInt64Scalar"
}
func GqlInt32TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLInt32Scalar"
}
func GqlUInt64TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLUInt64Scalar"
}
func GqlUInt32TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLUInt32Scalar"
}
func GqlFloat32TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLFloat32Scalar"
}
func GqlFloat64TypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".GraphQLFloat64Scalar"
}
func GqlMultipartFileTypeResolver(ctx BodyContext) string {
	return ctx.Importer.New(ScalarsPkgPath) + ".MultipartFile"
}
