package parser

import (
	"strings"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

type File struct {
	GoPackage string
	FilePath  string
	protoFile *proto.Proto
	PkgName   string
	Services  []*Service
	Messages  []*Message
	Enums     []*Enum
	Imports   []*File
}

func (f *File) parseGoPackage() {
	for _, el := range f.protoFile.Elements {
		option, ok := el.(*proto.Option)
		if !ok {
			continue
		}
		if option.Name == "go_package" {
			f.GoPackage = option.Constant.Source
		}
	}
}

func (f *File) messageByTypeName(typeName TypeName) (*Message, bool) {
	for _, msg := range f.Messages {
		if msg.TypeName.Equal(typeName) {
			return msg, true
		}
	}
	return nil, false
}

func (f *File) enumByTypeName(typeName TypeName) (*Enum, bool) {
	for _, e := range f.Enums {
		if e.TypeName.Equal(typeName) {
			return e, true
		}
	}
	return nil, false
}
func (f *File) findTypeInMessage(msg *Message, typ string) (Type, bool) {
	if typeIsScalar(typ) {
		return &ScalarType{ScalarName: typ, file: f}, true
	}
	ms, ok := f.messageByTypeName(msg.TypeName.NewSubTypeName(typ))
	if ok {
		return ms.Type, true
	}

	enum, ok := f.enumByTypeName(msg.TypeName.NewSubTypeName(typ))
	if ok {
		return enum.Type, true
	}
	if msg.parentMsg != nil {
		return f.findTypeInMessage(msg.parentMsg, typ)
	}
	return f.findType(typ)
}

func (f *File) findType(typ string) (Type, bool) {
	if typeIsScalar(typ) {
		return &ScalarType{ScalarName: typ, file: f}, true
	}
	parts := strings.Split(typ, ".")
	msg, ok := f.messageByTypeName(parts)
	if ok {
		return msg.Type, true
	}
	en, ok := f.enumByTypeName(parts)
	if ok {
		return en.Type, true
	}
	for _, imp := range f.Imports {
		if imp.PkgName == f.PkgName {
			it, ok := imp.findType(typ)
			if ok {
				return it, true
			}
		}
	}
	for i := 0; i < len(parts)-1; i++ {
		pkg, typ := strings.Join(parts[:i+1], "."), strings.Join(parts[i+1:], ".")
		for _, imp := range f.Imports {
			if imp.PkgName == pkg {
				t, ok := imp.findType(typ)
				if ok {
					return t, ok
				}
			}
		}
	}
	return nil, false
}

func (f *File) parseServices() error {
	for _, el := range f.protoFile.Elements {
		service, ok := el.(*proto.Service)
		if !ok {
			continue
		}
		srv := &Service{
			Name:          service.Name,
			QuotedComment: quoteComment(service.Comment),
		}
		for _, el := range service.Elements {
			method, ok := el.(*proto.RPC)
			if !ok || method.StreamsRequest || method.StreamsReturns {
				continue
			}
			reqTyp, ok := f.findType(method.RequestType)
			if !ok {
				return errors.Errorf("can't find request message %s", method.RequestType)
			}
			retTyp, ok := f.findType(method.ReturnsType)
			if !ok {
				return errors.Errorf("can't find request message %s", method.RequestType)
			}
			mtd := &Method{
				Name:          method.Name,
				QuotedComment: quoteComment(method.Comment),
				InputMessage:  reqTyp.(*MessageType).Message,
				OutputMessage: retTyp.(*MessageType).Message,
				Service:       srv,
			}
			srv.Methods = append(srv.Methods, mtd)
		}
		f.Services = append(f.Services, srv)
	}
	return nil
}

func (f *File) parseMessagesFields() error {
	for _, msg := range f.Messages {
		for _, el := range msg.Descriptor.Elements {
			switch fld := el.(type) {
			case *proto.NormalField:
				typ, ok := f.findTypeInMessage(msg, fld.Type)
				if !ok {
					return errors.Errorf("failed to find message %s field %s type", strings.Join(msg.TypeName, "."), fld.Name)
				}
				fl := &Field{
					Name:          fld.Name,
					QuotedComment: quoteComment(fld.Comment),
					Repeated:      fld.Repeated,
					descriptor:    fld.Field,
					Type:          typ,
				}
				msg.Fields = append(msg.Fields, fl)
			case *proto.MapField:
				ktyp, ok := f.findTypeInMessage(msg, fld.KeyType)
				if !ok {
					return errors.Errorf("failed to find message %s field %s type", strings.Join(msg.TypeName, "."), fld.Name)
				}
				vtyp, ok := f.findTypeInMessage(msg, fld.Type)
				if !ok {
					return errors.Errorf("failed to find message %s field %s type", strings.Join(msg.TypeName, "."), fld.Name)
				}
				mp := &Map{
					Message:   msg,
					KeyType:   ktyp,
					ValueType: vtyp,
					Field:     fld,
					File:      f,
				}
				t := &MapType{Map: mp, file: f}
				mp.Type = t
				mf := &MapField{
					Name:          fld.Name,
					QuotedComment: quoteComment(fld.Comment),
					descriptor:    fld,
					Type:          t,
					Map:           mp,
				}
				msg.MapFields = append(msg.MapFields, mf)
			case *proto.Oneof:
				of := &OneOf{
					Name: fld.Name,
				}
				for _, el := range fld.Elements {
					fld, ok := el.(*proto.OneOfField)
					if !ok {
						continue
					}
					typ, ok := f.findTypeInMessage(msg, fld.Type)
					if !ok {
						return errors.Errorf("failed to find message %s field %s type", strings.Join(msg.TypeName, "."), fld.Name)
					}
					of.Fields = append(of.Fields, &Field{
						Name:          fld.Name,
						QuotedComment: quoteComment(fld.Comment),
						Repeated:      false,
						descriptor:    fld.Field,
						Type:          typ,
					})
				}
				msg.OneOffs = append(msg.OneOffs, of)
			}
		}
	}
	return nil
}

func (f *File) parseMessages() {
	for _, el := range f.protoFile.Elements {
		msg, ok := el.(*proto.Message)
		if !ok {
			continue
		}
		m := message(f, msg, TypeName{msg.Name}, nil)
		f.Messages = append(f.Messages, m)
		f.parseMessagesInMessage(TypeName{msg.Name}, m)
	}
}

func (f *File) parseMessagesInMessage(msgTypeName TypeName, msg *Message) {
	for _, el := range msg.Descriptor.Elements {
		switch elv := el.(type) {
		case *proto.Message:
			tn := msgTypeName.NewSubTypeName(elv.Name)
			m := message(f, elv, tn, msg)
			f.Messages = append(f.Messages, m)
			f.parseMessagesInMessage(tn, m)
		}
	}
}

func (f *File) parseEnums() {
	for _, el := range f.protoFile.Elements {
		switch val := el.(type) {
		case *proto.Enum:
			f.Enums = append(f.Enums, newEnum(f, val, TypeName{val.Name}))
		case *proto.Message:
			f.parseEnumsInMessage(TypeName{val.Name}, val)
		}

	}
}

func (f *File) parseEnumsInMessage(msgTypeName TypeName, msg *proto.Message) {
	for _, el := range msg.Elements {
		switch elv := el.(type) {
		case *proto.Message:
			f.parseEnumsInMessage(msgTypeName.NewSubTypeName(elv.Name), elv)
		case *proto.Enum:
			f.Enums = append(f.Enums, newEnum(f, elv, msgTypeName.NewSubTypeName(elv.Name)))
		}
	}
}
