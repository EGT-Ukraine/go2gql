package swagger2gql

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

type parsedFile struct {
	File          *parser.File
	Config        *SwaggerFileConfig
	OutputPath    string
	OutputPkg     string
	OutputPkgName string
}

func (p *Plugin) fileOutputPath(cfg *SwaggerFileConfig) (string, error) {
	if cfg.GetOutputPath() == "" {
		if p.config.GetOutputPath() == "" {
			return "", errors.Errorf("need to specify global output_path")
		}
		absFilePath, err := filepath.Abs(cfg.Path)
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve file absolute path")
		}
		fileName := filepath.Base(cfg.Path)
		pkg, err := GoPackageByPath(filepath.Dir(absFilePath), p.generateConfig.VendorPath)
		var res string
		if err != nil {
			res, err = filepath.Abs(filepath.Join("./"+p.config.GetOutputPath()+"/", "./"+filepath.Dir(absFilePath), strings.TrimSuffix(fileName, ".json")+".go"))
		} else {
			res, err = filepath.Abs(filepath.Join("./"+p.config.GetOutputPath()+"/", "./"+pkg, strings.TrimSuffix(fileName, ".json")+".go"))
		}
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve absolute output path")
		}
		return res, nil
	}
	return filepath.Join(cfg.OutputPath, strings.TrimSuffix(filepath.Base(cfg.Path), ".json")+".go"), nil
}

func (p *Plugin) fileOutputPackage(cfg *SwaggerFileConfig) (name, pkg string, err error) {
	outPath, err := p.fileOutputPath(cfg)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve file output path")
	}
	pkg, err = GoPackageByPath(filepath.Dir(outPath), p.generateConfig.VendorPath)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve file go package")
	}
	return strings.Replace(filepath.Base(pkg), "-", "_", -1), pkg, nil
}
