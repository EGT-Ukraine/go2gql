package importer

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Import struct {
	Alias string
	Path  string
}
type Importer struct {
	CurrentPackage string
	imports        []Import
}

func (i *Importer) resolveImport(path string) (alias, importPath string) {
	paths := strings.Split(path, "/")
	if paths[len(paths)-1] == "" {
		// if Path is something like `a/b/c/`(ends with slash), Alias will be "c" Path will be "a/b/c"
		alias, importPath = paths[len(paths)-2], strings.Join(paths[:len(paths)-1], "/")
	} else {
		alias, importPath = paths[len(paths)-1], strings.Join(paths, "/")
	}
	alias = strings.NewReplacer("-", "_", " ", "_").Replace(alias)
	r, _ := utf8.DecodeRune([]byte(alias))
	if !unicode.IsLetter(r) {
		alias = "imp" + alias
	}
	return
}
func (i *Importer) findPath(path string) *Import {
	for _, imp := range i.imports {
		if imp.Path == path {
			return &imp
		}
	}
	return nil
}
func (i *Importer) aliasExists(alias string) bool {
	for _, imp := range i.imports {
		if imp.Alias == alias {
			return true
		}
	}
	return false
}
func (i *Importer) findAliasWithoutCollision(alias string) string {
	if !i.aliasExists(alias) {
		return alias
	}
	for j := 1; ; j++ {
		a := alias + "_" + strconv.Itoa(j)
		if !i.aliasExists(alias + "_" + strconv.Itoa(j)) {
			return a
		}
	}
}
func (i *Importer) New(path string) string {
	if i.CurrentPackage == path {
		return ""
	}
	alias, path := i.resolveImport(path)
	imp := i.findPath(path)
	if imp != nil {
		return imp.Alias
	}
	alias = i.findAliasWithoutCollision(alias)
	i.imports = append(i.imports, Import{
		Alias: alias,
		Path:  path,
	})
	return alias
}
func (i *Importer) Prefix(path string) string {
	if i.CurrentPackage == path || path == "" {
		return ""
	}
	return i.New(path) + "."
}

func (i *Importer) Imports() []Import {
	res := make([]Import, len(i.imports))
	copy(res, i.imports)
	return res
}
