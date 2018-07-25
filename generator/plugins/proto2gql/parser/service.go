package parser

type Service struct {
	Name          string
	QuotedComment string
	Methods       []*Method
}

type Method struct {
	Name          string
	QuotedComment string
	InputMessage  *Message
	OutputMessage *Message
	Service       *Service
}
