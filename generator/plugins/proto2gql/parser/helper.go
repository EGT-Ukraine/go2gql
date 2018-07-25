package parser

import (
	"strconv"
	"strings"

	"github.com/emicklei/proto"
)

func typeIsScalar(typ string) bool {
	switch typ {
	case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes":
		return true
	}
	return false
}

func quoteComment(comments *proto.Comment) string {
	if comments == nil {
		return `""`
	}
	return strconv.Quote(strings.TrimSpace(strings.Join(comments.Lines, "\n")))
}

func resolveFilePkgName(file *proto.Proto) string {
	for _, el := range file.Elements {
		if p, ok := el.(*proto.Package); ok {
			return p.Name
		}
	}
	return ""
}

func message(file *File, msg *proto.Message, typeName []string, parent *Message) *Message {
	m := &Message{
		Name:          msg.Name,
		QuotedComment: quoteComment(msg.Comment),
		Descriptor:    msg,
		TypeName:      typeName,
		File:          file,
		parentMsg:     parent,
	}
	m.Type = &MessageType{file: file, Message: m}
	return m
}
