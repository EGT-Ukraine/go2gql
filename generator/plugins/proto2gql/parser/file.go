package parser

import (
	"strings"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

type File struct {
	GoPackage   string
	FilePath    string
	protoFile   *proto.Proto
	PkgName     string
	Services    []*Service
	Messages    []*Message
	Enums       []*Enum
	Imports     []*File
	Descriptors map[string]Type
}

func (f *File) parseGoPackage() {
	for _, el := range f.protoFile.Elements {
		option, ok := el.(*proto.Option)
		if !ok {
			continue
		}
		if option.Name == "go_package" {
			parts := strings.Split(option.Constant.Source, ";")
			f.GoPackage = parts[0]
		}
	}
}

func (f *File) findTypeInMessage(msg *Message, typ string) (Type, bool) {
	if typeIsScalar(typ) {
		return &Scalar{ScalarName: typ, file: f}, true
	}

	return f.findType(typ, msg.GetFullName())
}

func (f *File) findSymbol(fullName string) (Type, bool) {
	symbol, ok := f.Descriptors[fullName]

	if ok {
		return symbol, ok
	}

	for _, importedFile := range f.Imports {
		symbol, ok := importedFile.findSymbol(fullName)

		if ok {
			return symbol, ok
		}
	}

	return nil, false
}

func (f *File) findType(name string, relativeToFullName string) (Type, bool) {
	var result Type

	var ok bool

	if strings.HasPrefix(name, ".") {
		// Fully-qualified name.
		result, ok = f.findSymbol(name[1:])
	} else {
		// We will search each parent scope of "relativeTo" looking for the
		// symbol.
		var scopeToTry = relativeToFullName + "."

		for {
			// Chop off the last component of the scope.
			var dotpos = strings.LastIndex(scopeToTry, ".")

			if dotpos == -1 {
				result, ok = f.findSymbol(name)
				break
			} else {
				scopeToTryRunes := []rune(scopeToTry)

				scopeToTry = string(scopeToTryRunes[0 : dotpos+1])

				// Append name and try to find.
				scopeToTry += name

				result, ok = f.findSymbol(scopeToTry)

				if ok {
					break
				}

				// Not found.  Remove the name so we can try again.
				scopeToTry = string(scopeToTryRunes[0:dotpos])
			}
		}
	}

	return result, ok
}

func (f *File) parseServices() error {
	for _, el := range f.protoFile.Elements {
		service, ok := el.(*proto.Service)
		if !ok {
			continue
		}
		srv := &Service{
			Name:          service.Name,
			QuotedComment: quoteComment(service.Comment, nil),
		}
		for _, el := range service.Elements {
			method, ok := el.(*proto.RPC)
			if !ok || method.StreamsRequest || method.StreamsReturns {
				continue
			}
			reqTyp, ok := f.findType(method.RequestType, f.PkgName)
			if !ok {
				return errors.Errorf("can't find request message %s", method.RequestType)
			}
			retTyp, ok := f.findType(method.ReturnsType, f.PkgName)
			if !ok {
				return errors.Errorf("can't find request message %s", method.RequestType)
			}
			mtd := &Method{
				Name:          method.Name,
				QuotedComment: quoteComment(method.Comment, method.InlineComment),
				InputMessage:  reqTyp.(*Message),
				OutputMessage: retTyp.(*Message),
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
				fl := &NormalField{
					Name:          fld.Name,
					QuotedComment: quoteComment(fld.Comment, fld.InlineComment),
					Repeated:      fld.Repeated,
					Optional:      fld.Optional,
					Required:      fld.Required,
					descriptor:    fld.Field,
					Type:          typ,
				}
				msg.NormalFields = append(msg.NormalFields, fl)
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
					file:      f,
				}
				mf := &MapField{
					Name:          fld.Name,
					QuotedComment: quoteComment(fld.Comment, fld.InlineComment),
					descriptor:    fld,
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
					of.Fields = append(of.Fields, &NormalField{
						Name:          fld.Name,
						QuotedComment: quoteComment(fld.Comment, fld.InlineComment),
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
		// TODO: implement.
		if msg.IsExtend {
			continue
		}
		m := message(f, msg, TypeName{msg.Name}, nil)
		f.Messages = append(f.Messages, m)
		f.parseMessagesInMessage(TypeName{msg.Name}, m)
		f.Descriptors[m.GetFullName()] = m
	}
}

func (f *File) parseMessagesInMessage(msgTypeName TypeName, msg *Message) {
	for _, el := range msg.Descriptor.Elements {
		elv, ok := el.(*proto.Message)

		if ok {
			tn := msgTypeName.NewSubTypeName(elv.Name)
			m := message(f, elv, tn, msg)
			f.Messages = append(f.Messages, m)
			f.Descriptors[m.GetFullName()] = m
			f.parseMessagesInMessage(tn, m)
		}
	}
}

func (f *File) parseEnums() {
	for _, el := range f.protoFile.Elements {
		switch val := el.(type) {
		case *proto.Enum:
			var enum = newEnum(f, val, TypeName{val.Name})
			f.Enums = append(f.Enums, enum)
			f.Descriptors[enum.GetFullName()] = enum
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
			var enum = newEnum(f, elv, msgTypeName.NewSubTypeName(elv.Name))
			f.Enums = append(f.Enums, enum)
			f.Descriptors[enum.GetFullName()] = enum
		}
	}
}
