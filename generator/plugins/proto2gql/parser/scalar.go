package parser

type Scalar struct {
	file       *File
	ScalarName string
}

func (s Scalar) String() string {
	return s.ScalarName
}

func (s Scalar) Kind() TypeKind {
	return TypeScalar
}

func (s Scalar) File() *File {
	return s.file
}
