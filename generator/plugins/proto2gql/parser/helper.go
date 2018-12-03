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

func quoteComment(comment *proto.Comment, inlineComment *proto.Comment) string {
	var lines []string

	if comment != nil {
		lines = append(lines, comment.Lines...)
	}

	if inlineComment != nil {
		lines = append(lines, inlineComment.Lines...)
	}

	if len(lines) == 0 {
		return `""`
	}

	return strconv.Quote(strings.TrimSpace(strings.Join(lines, "\n")))
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
		QuotedComment: quoteComment(msg.Comment, nil),
		Descriptor:    msg,
		TypeName:      typeName,
		file:          file,
		parentMsg:     parent,
	}

	return m
}
