package semantic

type (
	NodeType int

	Node interface {
		Type() NodeType
	}

	node struct {
	}

	Script struct {
		node
		Statements []Statement
	}
)

const (
	NilNode NodeType = iota
	ScriptNodeT
	StatementSelect
)

func (*node) Type() NodeType {
	return NilNode
}

func (*Script) Type() NodeType {
	return ScriptNodeT
}
