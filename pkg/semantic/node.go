package semantic

type (
	NodeType int

	Span struct {
		Start, End int
	}

	Node interface {
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

	node struct {
		line   int
		column int
		span   Span
	}

	Script struct {
		node
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

func (node) Type() NodeType {
	return NilNode
}

func (n node) Line() int {
	return n.line
}

func (n node) Column() int {
	return n.column
}

func (n node) Span() Span {
	return n.span
}

func (n *node) SetLine(line int) {
	n.line = line
}

func (n *node) SetColumn(column int) {
	n.column = column + 1
}

func (n *node) SetSpan(span Span) {
	n.span = span
}

func (*Script) Type() NodeType {
	return ScriptNode
}
