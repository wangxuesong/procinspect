package semantic

type (
	NodeType int

	Node interface {
		Type() NodeType
		Line() int
		Column() int
	}

	node struct {
		line   int
		column int
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

func (*node) Type() NodeType {
	return NilNode
}

func (n *node) Line() int {
	return n.line
}

func (n *node) Column() int {
	return n.column
}

func (n *node) SetLine(line int) {
	n.line = line
}

func (n *node) SetColumn(column int) {
	n.column = column + 1
}

func (*Script) Type() NodeType {
	return ScriptNode
}
