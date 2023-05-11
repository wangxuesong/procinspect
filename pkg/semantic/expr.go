package semantic

type (
	Expr interface {
		Node
		expr()
	}

	expression struct {
		node
	}

	NumericLiteral struct {
		expression
		Value int64
	}

	CursorAttribute struct {
		expression
		Cursor string
		Attr   string
	}

	UnaryLogicalExpression struct {
		expression
		Expr     Expr
		Operator string
		Not      bool
	}

	NameExpression struct {
		expression
		Name string
	}
)

func (n *expression) expr() {}
