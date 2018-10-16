package parser

import "strings"

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
	GetFullName() string
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

func (t ScalarType) GetFullName() string {
	return ""
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

func (t EnumType) GetFullName() string {
	if t.file.PkgName != "" {
		return t.file.PkgName + "." + strings.Join(t.Enum.TypeName, ".")
	} else {
		return strings.Join(t.Enum.TypeName, ".")
	}
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

func (t MapType) GetFullName() string {
	return ""
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

func (t MessageType) GetFullName() string {
	parentMessage := t.Message.parentMsg

	if parentMessage != nil {
		parentMessageName := parentMessage.Type.GetFullName()

		return parentMessageName + "." + t.Message.Name
	} else {
		if t.file.PkgName != "" {
			return t.file.PkgName + "." + t.Message.Name
		} else {
			return t.Message.Name
		}
	}
}
