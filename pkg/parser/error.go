package parser

type ErrorKind int

func (k ErrorKind) String() string {
	switch k {
	case ErrSyntax:
		return "syntax"
	case ErrSemantic:
		return "semantic"
	default:
		return "unknown"
	}
}

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
