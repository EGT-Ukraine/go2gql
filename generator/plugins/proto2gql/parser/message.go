package parser

import "github.com/emicklei/proto"

type Messages []*Message

func (m Messages) Copy() Messages {
	result := make(Messages, len(m))
	copy(result, m)

	return result
}

func (m Messages) Contains(msg *Message) bool {
	for _, value := range m {
		if value == msg {
			return true
		}
	}

	return false
}

type Message struct {
	Name          string
	QuotedComment string
	NormalFields  []*NormalField
	MapFields     []*MapField
	OneOffs       []*OneOf
	Descriptor    *proto.Message
	TypeName      TypeName
	file          *File
	parentMsg     *Message
}

func (m Message) GetFields() []Field {
	var res []Field
	for _, field := range m.NormalFields {
		res = append(res, field)
	}
	for _, field := range m.MapFields {
		res = append(res, field)
	}
	for _, oneOff := range m.OneOffs {
		for _, field := range oneOff.Fields {
			res = append(res, field)
		}
	}

	return res
}

func (m Message) HaveFields() bool {
	if len(m.NormalFields) > 0 || len(m.MapFields) > 0 {
		return true
	}
	for _, of := range m.OneOffs {
		if len(of.Fields) > 0 {
			return true
		}
	}

	return false
}

func (m Message) GetFieldByName(name string) (Field, bool) {
	for _, field := range m.GetFields() {
		if field.GetName() == name {
			return field, true
		}
	}

	return nil, false
}

func (m Message) HaveFieldsExcept(field string) bool {
	for _, f := range m.NormalFields {
		if f.Name != field {
			return true
		}
	}
	for _, f := range m.MapFields {
		if f.Name != field {
			return true
		}
	}
	for _, of := range m.OneOffs {
		for _, f := range of.Fields {
			if f.Name != field {
				return true
			}
		}
	}

	return false
}

func (m Message) File() *File {
	return m.file
}

func (m Message) Kind() TypeKind {
	return TypeMessage
}

func (m Message) String() string {
	return m.Name + " message"
}

func (m Message) GetFullName() string {
	parentMessage := m.parentMsg

	if parentMessage != nil {
		parentMessageName := parentMessage.GetFullName()

		return parentMessageName + "." + m.Name
	}

	if m.file.PkgName != "" {
		return m.file.PkgName + "." + m.Name
	}

	return m.Name
}

type Field interface {
	GetName() string
	GetType() Type
}

type NormalField struct {
	Name          string
	QuotedComment string
	Repeated      bool
	descriptor    *proto.Field
	Type          Type
	Optional      bool
	Required      bool
}

func (n *NormalField) GetName() string {
	return n.Name
}

func (n *NormalField) GetType() Type {
	return n.Type
}

type MapField struct {
	Name          string
	QuotedComment string
	descriptor    *proto.MapField
	Map           *Map
}

func (n *MapField) GetName() string {
	return n.Name
}

func (n *MapField) GetType() Type {
	return n.Map
}

type OneOf struct {
	Name   string
	Fields []*NormalField
}

type Map struct {
	Message   *Message
	KeyType   Type
	ValueType Type
	Field     *proto.MapField
	file      *File
}

func (m Map) File() *File {
	return m.file
}

func (m Map) Kind() TypeKind {
	return TypeMap
}

func (m Map) String() string {
	return m.Message.Name + "." + m.Field.Name + " map"
}
