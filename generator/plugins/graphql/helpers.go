package graphql

import (
	"go/build"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func typeIsScalar(p GoType) bool {
	switch p.Kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		return true
	}
	return false
}

func ResolverCall(resolverPkg, resolverFuncName string) ValueResolver {
	return func(arg string, ctx BodyContext) string {
		if ctx.TracerEnabled {
			return ctx.Importer.Prefix(resolverPkg) + resolverFuncName + "(tr, " + ctx.Importer.New(OpentracingPkgPath) + ".ContextWithSpan(ctx, span), " + arg + ")"
		}
		return ctx.Importer.Prefix(resolverPkg) + resolverFuncName + "(ctx, " + arg + ")"
	}
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

func IdentAccessValueResolver(ident string) ValueResolver {
	return func(arg string, ctx BodyContext) string {
		return arg + "." + ident
	}
}
