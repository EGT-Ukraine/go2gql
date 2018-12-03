package parser

import (
	"strings"

	"github.com/emicklei/proto"
)

type Enum struct {
	Name          string
	QuotedComment string
	Values        []*EnumValue
	file          *File
	TypeName      TypeName
	Descriptor    *proto.Enum
}

type EnumValue struct {
	Name          string
	Value         int
	QuotedComment string
}

func newEnum(file *File, enum *proto.Enum, typeName []string) *Enum {
	m := &Enum{
		Name:          enum.Name,
		QuotedComment: quoteComment(enum.Comment, nil),
		Descriptor:    enum,
		TypeName:      typeName,
		file:          file,
	}
	for _, v := range enum.Elements {
		value, ok := v.(*proto.EnumField)
		if !ok {
			continue
		}
		m.Values = append(m.Values, &EnumValue{
			Name:          value.Name,
			Value:         value.Integer,
			QuotedComment: quoteComment(value.Comment, value.InlineComment),
		})
	}

	return m
}

func (e Enum) Kind() TypeKind {
	return TypeEnum
}

func (e Enum) String() string {
	return e.Name + " enum"
}

func (e Enum) File() *File {
	return e.file
}

func (e Enum) GetFullName() string {
	if e.file.PkgName == "" {
		return strings.Join(e.TypeName, ".")
	}

	return e.file.PkgName + "." + strings.Join(e.TypeName, ".")
}
