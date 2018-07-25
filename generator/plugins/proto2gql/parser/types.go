package parser

type Kind byte

const (
	TypeUndefined Kind = iota
	TypeScalar
	TypeMessage
	TypeEnum
	TypeMap
)

type Type interface {
	Kind() Kind
	String() string
	File() *File
}

type ScalarType struct {
	file       *File
	ScalarName string
}

func (t ScalarType) Kind() Kind {
	return TypeScalar
}

func (t ScalarType) String() string {
	return t.ScalarName
}

func (t ScalarType) File() *File {
	return t.file
}

type EnumType struct {
	file *File
	Enum *Enum
}

func (t EnumType) Kind() Kind {
	return TypeEnum
}

func (t EnumType) String() string {
	return t.Enum.Name + " enum"
}

func (t EnumType) File() *File {
	return t.file
}

type MapType struct {
	file *File
	Map  *Map
}

func (t MapType) Kind() Kind {
	return TypeMap
}

func (t MapType) String() string {
	return t.Map.Message.Name + "." + t.Map.Field.Name + " map"
}

func (t MapType) File() *File {
	return t.file
}

type MessageType struct {
	file    *File
	Message *Message
}

func (t MessageType) Kind() Kind {
	return TypeMessage
}

func (t MessageType) String() string {
	return t.Message.Name + " message"
}

func (t MessageType) File() *File {
	return t.file
}
