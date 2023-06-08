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

	RelationalExpression struct {
		expression
		Left     Expr
		Right    Expr
		Operator string
	}

	BinaryExpression struct {
		expression
		Left     Expr
		Right    Expr
		Operator string
	}

	FunctionCallExpression struct {
		expression
		Name Expr
		Args []Expr
	}

	DotExpression struct {
		expression
		Name   string
		Parent Expr
	}

	NameExpression struct {
		expression
		Name string
	}

	StringLiteral struct {
		expression
		Value string
	}
)

func (n *expression) expr() {}
