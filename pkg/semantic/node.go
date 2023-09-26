package semantic

type (
	NodeType int

	Span struct {
		Start, End int
	}

	Node interface {
		AstNode
		Type() NodeType
		Line() int
		Column() int
		Span() Span
	}

	SetPosition interface {
		SetLine(int)
		SetColumn(int)
		SetSpan(Span)
	}

	SyntaxNode struct {
		SourceLine int
		SourceCol  int
		SourceSpan Span
	}

	Script struct {
		SyntaxNode
		Statements []Statement
	}
)

const (
	NilNode NodeType = iota
	ScriptNode
	CreateProcedure
	StatementSelect
	Assignment
)

func (SyntaxNode) Type() NodeType {
	return NilNode
}

func (n SyntaxNode) Line() int {
	return n.SourceLine
}

func (n SyntaxNode) Column() int {
	return n.SourceCol
}

func (n SyntaxNode) Span() Span {
	return n.SourceSpan
}

func (n *SyntaxNode) SetLine(line int) {
	n.SourceLine = line
}

func (n *SyntaxNode) SetColumn(column int) {
	n.SourceCol = column + 1
}

func (n *SyntaxNode) SetSpan(span Span) {
	n.SourceSpan = span
}

func (*Script) Type() NodeType {
	return ScriptNode
}
