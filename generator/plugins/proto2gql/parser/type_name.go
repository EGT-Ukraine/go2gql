package parser

type TypeName []string

func (t TypeName) NewSubTypeName(subName string) TypeName {
	var res = make(TypeName, len(t)+1)
	copy(res, t)
	res[len(t)] = subName
	return res
}
func (t TypeName) Equal(t2 TypeName) bool {
	if len(t) != len(t2) {
		return false
	}
	for i, v := range t {
		if t2[i] != v {
			return false
		}
	}
	return true
}
