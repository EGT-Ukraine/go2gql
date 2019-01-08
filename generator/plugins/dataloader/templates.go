// Code generated by go-bindata.
// sources:
// generator/plugins/dataloader/templates/loaders.gohtml
// generator/plugins/dataloader/templates/output_object_fields.gohtml
// DO NOT EDIT!

package dataloader

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesLoadersGohtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x54\xc1\x6a\x1b\x31\x10\x3d\x4b\x5f\x31\x98\x50\xbc\xc6\xd9\x85\x1e\x0d\x3e\x14\xa7\x31\xa1\x4d\x29\x6d\xda\x1c\x42\x0e\xaa\x76\xbc\x16\x91\x25\x47\x3b\xeb\x26\x08\xfd\x7b\xd1\xae\x1c\x6f\x92\x0d\xce\xa1\x27\x1b\xcd\xd3\xbc\x37\xf3\xde\xca\xfb\x53\x28\x26\x95\xa5\xc7\x2d\xce\xa0\x52\xb4\x6e\xfe\xe4\xd2\x6e\x8a\xcf\xcb\xab\xd3\x5f\x77\x4e\x28\x83\x45\x65\x3f\x56\xf7\xba\xa8\xd0\xa0\x13\x64\x5d\xb1\xd5\x4d\xa5\x4c\x5d\x94\x82\x84\xb6\xa2\x44\x97\x7f\x6d\x7f\xea\x85\x35\x84\x0f\x34\x29\xe0\x34\x04\xce\xb7\x42\xde\x89\x0a\xa1\x03\xd5\x9c\xab\xcd\xd6\x3a\x82\x31\x67\x23\xd9\x41\x47\x9c\x8d\x48\x6d\x70\xc4\x39\xf3\xde\x09\x53\x21\x9c\x24\xd8\x6c\x0e\x27\xf9\x45\xfb\xbf\x6e\x1b\x32\xe6\x7d\x2a\xe6\x9f\xb4\x12\x75\x08\x30\x3a\x1c\x7d\x17\xb4\x0e\x61\x14\x1b\xa1\x29\xdb\x1b\x19\xe7\x71\x36\xe8\xf4\x2d\xb4\x42\x43\x35\x28\x43\xe8\x56\x42\x22\xf8\x1e\x6b\xa7\xb2\x63\x4d\xf3\x24\xd6\x25\x92\xf7\xa9\x9e\xff\x44\xb7\x53\x12\xf3\x6f\x62\x83\x21\x74\x2d\xc7\x19\x78\x5f\xd9\xab\x48\xf5\x12\xb7\x10\x5a\x5f\xec\x09\x63\xbb\x83\xba\x90\xd4\x9d\x09\x12\x7b\xc6\x9a\x5c\x23\xe9\x9d\xc2\x0e\xaa\x16\xd6\xac\x54\x95\x44\x75\x20\x38\x54\xbb\x83\x28\xaf\x0f\x18\x92\x52\x1e\xa4\x24\x33\xbf\xe0\x63\x3b\x57\x27\xcc\x07\xce\x77\xc2\x0d\xe3\x60\xfe\xf6\xfd\x78\x71\xd5\x18\x09\x4b\xa4\x54\xb9\x56\xb4\x4e\xd8\xb1\xa4\x07\x48\x91\xc8\x53\x79\x0a\x62\xab\xf6\x96\x3d\x33\x30\x7b\x09\x8d\xeb\xea\x31\xc7\x55\x7d\xe8\x2d\xf5\x3f\x2c\x73\x06\xd2\xa1\x20\x7c\x03\x14\xf5\xf7\xf5\xe6\xef\x88\x4c\x36\xed\x1b\xc0\x02\xe7\xcc\x21\x35\xce\x3c\x4d\x17\x17\xf4\x5b\xe8\x06\xbb\xf6\x83\xbb\x7d\x76\x9c\x45\x1f\x8f\x8e\xda\xda\x70\x74\x9c\xd7\x76\xc8\x56\xf8\xbb\x93\x9e\x1d\x0d\x60\xb4\x2d\x8d\x7c\x0c\xea\x39\x63\x2b\x24\xb9\x9e\x41\x94\x3f\xbe\xc3\xc7\xfa\xb5\x92\x1f\x78\xdf\x60\x4d\xcb\xf6\x34\x2a\x18\xdf\xdc\x0e\x80\xea\xad\x35\x35\xee\x51\x53\xb8\xb9\x45\xe7\xac\xcb\xa2\x9e\x7e\x0e\xce\x23\xe1\xc2\x96\xed\x67\xcb\xc2\x94\x33\xf6\x57\x28\x9a\x81\xf7\x65\xe3\x04\x29\x6b\xe0\xc5\xfe\xae\x85\xa2\xb3\x54\xbb\x8c\xef\xd3\x04\xe2\xf3\x96\x5f\x2a\xad\x55\x8d\xd2\x9a\x72\x1a\xbd\x0e\xfc\x60\xfd\xd3\x77\xd1\x8b\xec\xb9\xb3\x9b\xb4\xf8\x21\x33\x32\x98\xf4\x1f\x0d\xcf\xd9\x4e\xe8\xe8\xb4\xa4\x87\xbc\x8b\xcc\x60\x5a\x32\xce\x99\x5a\x41\x04\xcf\xe7\x60\x94\x6e\x47\x4e\x1e\x18\xa5\x9f\xc5\x70\x27\x74\x3e\xee\xf3\xc4\x78\xfd\x0b\x00\x00\xff\xff\xbd\x30\xe8\xa7\x33\x06\x00\x00")

func templatesLoadersGohtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesLoadersGohtml,
		"templates/loaders.gohtml",
	)
}

func templatesLoadersGohtml() (*asset, error) {
	bytes, err := templatesLoadersGohtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/loaders.gohtml", size: 1587, mode: os.FileMode(436), modTime: time.Unix(1546955364, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesOutput_object_fieldsGohtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x53\xc1\x8e\xda\x30\x10\xbd\xef\x57\x8c\xd0\xaa\x0a\x08\x1c\xa9\x47\x24\x0e\x15\xdb\xe5\xd0\x8a\x45\x5b\xda\x9e\x4d\x32\x31\x2e\xc6\x0e\x8e\xd3\x76\x6b\xf9\xdf\xab\x71\x4c\x02\xe9\xa5\xbe\xc4\xce\xcc\x7b\xf3\xde\x73\xe2\xfd\x02\xf2\x99\x30\xee\xad\xc6\x25\x08\xe9\x8e\xed\x81\x15\xe6\x9c\x7f\xdc\xec\x17\x5f\x4f\x96\x4b\x8d\xb9\x30\xef\xc5\x45\xe5\x02\x35\x5a\xee\x8c\xcd\x6b\xd5\x0a\xa9\x9b\x5c\x58\x5e\x1f\x2f\x8a\xbd\xa2\x2e\xd1\x3e\x4b\x54\x65\xb3\x36\xda\xe1\x6f\x37\xcb\x61\x11\xc2\x83\xf7\x60\xb9\x16\x08\x8f\x15\x55\x61\xb9\x82\x47\xf6\xd2\xba\xba\x75\x2f\x87\x1f\x58\x38\xf6\xc4\x1d\xff\x6c\x78\x8f\x8f\x30\x00\x00\xef\x47\x9d\xdf\xb8\x95\xfc\xa0\x70\xcb\xcf\x18\x02\xfb\x50\x96\x11\xb1\x36\xba\x92\x22\x9b\x78\xdf\xcd\x60\x5d\x7d\x32\x87\x77\xde\x8b\x8b\xda\x9d\x44\x08\x2c\xb6\xfa\x48\x4c\x8b\x7a\x96\xd7\xc3\xbf\xd0\xbe\xef\x09\x9b\xc2\xca\xda\x49\xa3\x97\x30\xb9\x29\xec\x63\x60\x69\x79\x9f\x92\xe8\xf4\x76\x76\xa8\x83\xf8\xc8\x70\x34\x90\x92\x49\x51\xdc\x18\xef\xa6\x0e\xdc\xaf\xd8\x18\xf5\x13\x97\x50\xb5\xba\xc8\x6a\x18\x62\x8e\xef\x77\xdc\xf2\x73\x33\x85\x4c\x6a\x87\xb6\xe2\x05\xfa\x30\x07\xb4\xd6\xd8\x29\x0c\x16\x69\xd5\xdc\xa2\x76\x94\x7a\xcd\xbe\x98\xd6\x16\xc8\xb2\x99\xf7\xc2\x90\xb8\xf1\x4d\x6c\xe2\xdb\x10\xa6\x0f\x77\x1c\x2a\x6a\x6c\x88\xc4\xfb\x74\xe8\x32\xdd\xa0\x1b\x4c\x34\xcf\xd6\x9c\x93\xc5\xac\x66\x69\x37\xe2\x92\x55\x4f\xb7\x5a\x81\x96\x6a\xa4\x97\x96\x45\xd7\x5a\x4d\xc5\x64\xaa\x61\x5b\xfc\x95\x4d\x68\x54\x8f\xd6\xc6\x41\x65\x5a\x5d\x82\xd4\x50\x74\xb3\x18\xac\xb9\x52\xd7\x16\x52\x97\x44\x7c\x97\xee\x98\x44\x4e\xa6\x77\xf3\xc2\xbd\x3c\x77\x6c\xf5\x89\x8c\x5e\x39\xfa\x0f\x63\x7c\x59\xdd\x9e\xd1\x63\x4f\xa0\xac\x4b\x7a\x00\x6c\x8d\x3d\x73\x25\xff\x60\xb9\x8b\x95\x4f\xf8\x16\x3f\xc1\x0e\x3e\x8a\x25\x59\x8e\xd7\xfd\x5f\xf7\x7a\x03\x8a\x9a\xb3\x91\xad\x39\xc5\xf7\x30\x1c\xe3\x36\x4c\xe9\x77\x44\x5d\xd2\x2f\xf6\x37\x00\x00\xff\xff\x89\xe3\x1a\xa0\xf9\x03\x00\x00")

func templatesOutput_object_fieldsGohtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesOutput_object_fieldsGohtml,
		"templates/output_object_fields.gohtml",
	)
}

func templatesOutput_object_fieldsGohtml() (*asset, error) {
	bytes, err := templatesOutput_object_fieldsGohtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/output_object_fields.gohtml", size: 1017, mode: os.FileMode(436), modTime: time.Unix(1546955341, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/loaders.gohtml":              templatesLoadersGohtml,
	"templates/output_object_fields.gohtml": templatesOutput_object_fieldsGohtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"loaders.gohtml":              &bintree{templatesLoadersGohtml, map[string]*bintree{}},
		"output_object_fields.gohtml": &bintree{templatesOutput_object_fieldsGohtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}