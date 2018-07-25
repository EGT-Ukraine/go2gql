package multipart_file

import (
	"mime/multipart"
)

type MultipartFile struct {
	multipart.File
	Header *multipart.FileHeader
}

func (m MultipartFile) Name() string {
	if m.Header == nil {
		return ""
	}
	return m.Header.Filename
}
