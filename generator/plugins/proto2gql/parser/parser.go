package parser

import (
	"os"
	"path/filepath"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

type Parser struct {
	parsedFiles []*File
}

func (p *Parser) ParsedFiles() []*File {
	return p.parsedFiles
}
func (p *Parser) parsedFile(filePath string) (*File, bool) {
	for _, f := range p.parsedFiles {
		if f.FilePath == filePath {
			return f, true
		}
	}
	return nil, false
}

func (p *Parser) importFilePath(filename string, importsAliases []map[string]string, paths []string) (filePath string, err error) {
	for _, aliases := range importsAliases {
		if v, ok := aliases[filename]; ok {
			filename = v
			break
		}
	}

	for _, path := range paths {
		p := filepath.Join(path, filename)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", errors.Errorf("can't find import %s in any of %s", filename, paths)
}

func (p *Parser) parseFileImports(file *File, importsAliases []map[string]string, paths []string) error {
	for _, v := range file.protoFile.Elements {
		imprt, ok := v.(*proto.Import)
		if !ok {
			continue
		}
		imprtPath, err := p.importFilePath(imprt.Filename, importsAliases, paths)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve import(%s) File path", imprt.Filename)
		}
		absImprtPath, err := filepath.Abs(imprtPath)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve import(%s) absolute File path", imprt.Filename)
		}
		if fl, ok := p.parsedFile(absImprtPath); ok {
			file.Imports = append(file.Imports, fl)
			continue
		}
		importFile, err := p.Parse(absImprtPath, importsAliases, paths)
		if err != nil {
			return errors.Wrapf(err, "can't parse import %s", imprtPath)
		}
		file.Imports = append(file.Imports, importFile)
	}
	return nil
}

func (p *Parser) Parse(path string, importAliases []map[string]string, paths []string) (*File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve File absolute path")
	}
	if pf, ok := p.parsedFile(absPath); ok {
		return pf, nil
	}
	file, err := os.Open(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open File")
	}

	f, err := proto.NewParser(file).Parse()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse File")
	}
	result := &File{
		FilePath:    absPath,
		protoFile:   f,
		PkgName:     resolveFilePkgName(f),
		Descriptors: map[string]Type{},
	}
	result.parseGoPackage()
	err = p.parseFileImports(result, importAliases, paths)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse File imports")
	}
	result.parseMessages()
	result.parseEnums()
	err = result.parseServices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse File services")
	}
	err = result.parseMessagesFields()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse messages fields")
	}
	p.parsedFiles = append(p.parsedFiles, result)
	return result, nil
}
