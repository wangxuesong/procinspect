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

	InExpression struct {
		expression
		Expr  Expr
		Elems []Expr
	}

	LikeExpression struct {
		expression
		Expr     Expr
		LikeExpr Expr
	}

	BetweenExpression struct {
		expression
		Expr  Expr
		Elems []Expr
	}

	ExistsExpression struct {
		expression
		Expr Expr
	}

	QueryExpression struct {
		expression
		Query *SelectStatement
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
		Name   Expr
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

	NullExpression struct {
		expression
	}

	SignExpression struct {
		expression
		Expr Expr
		Sign string
	}

	OuterJoinExpression struct {
		expression
		Expr Expr
	}

	AliasExpression struct {
		expression
		Expr  Expr
		Alias string
	}

	StatementExpression struct {
		expression
		Stmt Statement
	}

	CastExpression struct {
		expression
		Expr     Expr
		DataType string
	}
)

func (n *expression) expr() {}
