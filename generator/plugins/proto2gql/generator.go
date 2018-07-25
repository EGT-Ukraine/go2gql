package proto2gql

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

var scalarsResolvers = map[string]graphql.TypeResolver{
	"double": graphql.GqlFloat64TypeResolver,
	"float":  graphql.GqlFloat32TypeResolver,
	"bool":   graphql.GqlBoolTypeResolver,
	"string": graphql.GqlStringTypeResolver,

	"int64":    graphql.GqlInt64TypeResolver,
	"sfixed64": graphql.GqlInt64TypeResolver,
	"sint64":   graphql.GqlInt64TypeResolver,

	"int32":    graphql.GqlInt32TypeResolver,
	"sfixed32": graphql.GqlInt32TypeResolver,
	"sint32":   graphql.GqlInt32TypeResolver,

	"uint32":  graphql.GqlUInt32TypeResolver,
	"fixed32": graphql.GqlUInt32TypeResolver,

	"uint64":  graphql.GqlUInt64TypeResolver,
	"fixed64": graphql.GqlUInt64TypeResolver,
}

type parsedFile struct {
	File           *parser.File
	Config         *ProtoFileConfig
	OutputPath     string
	OutputPkg      string
	OutputPkgName  string
	GRPCSourcesPkg string
}
type Proto2GraphQL struct {
	VendorPath      string
	GenerateTracers bool
	parser          parser.Parser
	ParsedFiles     []*parsedFile
}

func (g *Proto2GraphQL) parsedFile(file *parser.File) (*parsedFile, error) {
	for _, f := range g.ParsedFiles {
		if f.File == file {
			return f, nil
		}
	}
	for _, f := range g.ParsedFiles {
		if f.File.FilePath == file.FilePath {
			return f, nil
		}
	}
	outPath, err := g.fileOutputPath(nil, file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve file '%s' output path", file.FilePath)
	}
	outPkgName, outPkg, err := g.fileOutputPackage(nil, file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve file '%s' output Go package", file.FilePath)
	}
	grpcPkg, err := g.fileGRPCSourcesPackage(nil, file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve file '%s' GRPC sources Go package", file.FilePath)
	}
	res := &parsedFile{
		File:           file,
		Config:         nil,
		OutputPath:     outPath,
		OutputPkg:      outPkg,
		OutputPkgName:  outPkgName,
		GRPCSourcesPkg: grpcPkg,
	}
	g.ParsedFiles = append(g.ParsedFiles, res)
	return res, nil

}
func (g *Proto2GraphQL) prepareFile(file *parsedFile) (*graphql.TypesFile, error) {
	enums, err := g.prepareFileEnums(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve file enums")
	}
	inputs, err := g.fileInputObjects(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file input objects")
	}
	mapInputs, err := g.fileMapInputObjects(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map input objects")
	}
	mapOutputs, err := g.fileMapOutputObjects(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map output objects")
	}
	mapResolvers, err := g.fileInputMapResolvers(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map resolvers")
	}
	outputMessages, err := g.fileOutputMessages(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file output messages")
	}
	messagesResolvers, err := g.fileInputMessagesResolvers(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file messages resolvers")
	}
	services, err := g.fileServices(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file services")
	}
	res := &graphql.TypesFile{
		PackageName:             file.OutputPkgName,
		Package:                 file.OutputPkg,
		Enums:                   enums,
		InputObjects:            inputs,
		InputObjectResolvers:    messagesResolvers,
		OutputObjects:           outputMessages,
		MapInputObjects:         mapInputs,
		MapInputObjectResolvers: mapResolvers,
		MapOutputObjects:        mapOutputs,
		Services:                services,
	}
	return res, nil
}

func (g *Proto2GraphQL) fileConfig(file *parser.File) *ProtoFileConfig {
	for _, f := range g.ParsedFiles {
		if f.File == file {
			return f.Config
		}
	}
	return nil
}

// fileGRPCSourcesPackage returns golang package of protobuf golang sources
func (g *Proto2GraphQL) fileGRPCSourcesPackage(cfg *ProtoFileConfig, file *parser.File) (string, error) {
	cfgGoPkg := cfg.GetGoPackage()

	if cfgGoPkg != "" {
		return cfgGoPkg, nil
	}

	if file.GoPackage != "" {
		return file.GoPackage, nil
	}
	fileDir := filepath.Dir(file.FilePath)
	pkg, err := GoPackageByPath(fileDir, g.VendorPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to resolve resolve go package of '%s'", fileDir)
	}
	return pkg, nil
}

func (g *Proto2GraphQL) fileOutputPath(cfg *ProtoFileConfig, file *parser.File) (string, error) {
	if cfg.GetOutputPath() == "" {
		absFilePath, err := filepath.Abs(file.FilePath)
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve file absolute path")
		}
		fileName := filepath.Base(file.FilePath)
		pkg, err := GoPackageByPath(filepath.Dir(absFilePath), g.VendorPath)
		var res string
		if err != nil {
			res, err = filepath.Abs(filepath.Join("./out/", "./"+filepath.Dir(absFilePath), strings.TrimSuffix(fileName, ".proto")+".go"))
		} else {
			res, err = filepath.Abs(filepath.Join("./out/", "./"+pkg, strings.TrimSuffix(fileName, ".proto")+".go"))
		}
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve absolute output path")
		}
		return res, nil
	}
	return filepath.Join(cfg.OutputPath, strings.TrimSuffix(filepath.Base(file.FilePath), ".proto")+".go"), nil
}

func (g *Proto2GraphQL) fileOutputPackage(cfg *ProtoFileConfig, file *parser.File) (name, pkg string, err error) {
	outPath, err := g.fileOutputPath(cfg, file)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve file output path")
	}
	pkg, err = GoPackageByPath(filepath.Dir(outPath), g.VendorPath)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve file go package")
	}
	return strings.Replace(filepath.Base(pkg), "-", "_", -1), pkg, nil
}

// AddSourceByConfig parse source proto files according to config definition
func (g *Proto2GraphQL) AddSourceByConfig(config *ProtoFileConfig) error {
	// here we start parse files by absolute path
	file, err := g.parser.Parse(config.ProtoPath, config.ImportsAliases, config.Paths)
	if err != nil {
		return errors.Wrap(err, "failed to parse proto file")
	}
	outPath, err := g.fileOutputPath(config, file)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve file '%s' output path", file.FilePath)
	}
	outPkgName, outPkg, err := g.fileOutputPackage(config, file)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve file '%s' output Go package", file.FilePath)
	}
	grpcPkg, err := g.fileGRPCSourcesPackage(config, file)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve file '%s' GRPC sources Go package", file.FilePath)
	}
	g.ParsedFiles = append(g.ParsedFiles, &parsedFile{
		File:           file,
		Config:         config,
		OutputPath:     outPath,
		OutputPkg:      outPkg,
		OutputPkgName:  outPkgName,
		GRPCSourcesPkg: grpcPkg,
	})
	return nil
}
