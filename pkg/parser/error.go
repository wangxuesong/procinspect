package parser

type ErrorKind int

type ParseError struct {
	Kind   ErrorKind
	Line   int
	Column int
	Msg    string
}

func (p ParseError) Error() string {
	return p.Msg
}

const (
	ErrSyntax   ErrorKind = iota
	ErrSemantic ErrorKind = iota
)
