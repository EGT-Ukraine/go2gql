package parser

type TypeKind byte

const (
	TypeUndefined TypeKind = iota
	TypeScalar
	TypeMessage
	TypeEnum
	TypeMap
)

type Type interface {
	String() string
	Kind() TypeKind
	File() *File
}
