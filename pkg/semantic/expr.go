package semantic

type (
	Expr interface {
		Node
		expr()
	}

	ExprNode struct {
		SyntaxNode
	}

	NumericLiteral struct {
		ExprNode
		Value int64
	}

	CursorAttribute struct {
		ExprNode
		Cursor string
		Attr   string
	}

	UnaryLogicalExpression struct {
		ExprNode
		Expr     Expr
		Operator string
		Not      bool
	}

	RelationalExpression struct {
		ExprNode
		Left     Expr
		Right    Expr
		Operator string
	}

	InExpression struct {
		ExprNode
		Expr  Expr
		Elems []Expr
	}

	LikeExpression struct {
		ExprNode
		Expr     Expr
		LikeExpr Expr
	}

	BetweenExpression struct {
		ExprNode
		Expr  Expr
		Elems []Expr
	}

	ExistsExpression struct {
		ExprNode
		Expr Expr
	}

	QueryExpression struct {
		ExprNode
		Query *SelectStatement
	}

	BinaryExpression struct {
		ExprNode
		Left     Expr
		Right    Expr
		Operator string
	}

	FunctionCallExpression struct {
		ExprNode
		Name Expr
		Args []Expr
	}

	DotExpression struct {
		ExprNode
		Name   Expr
		Parent Expr
	}

	NameExpression struct {
		ExprNode
		Name string
	}

	StringLiteral struct {
		ExprNode
		Value string
	}

	NullExpression struct {
		ExprNode
	}

	SignExpression struct {
		ExprNode
		Expr Expr
		Sign string
	}

	OuterJoinExpression struct {
		ExprNode
		Expr Expr
	}

	AliasExpression struct {
		ExprNode
		Expr  Expr
		Alias string
	}

	StatementExpression struct {
		ExprNode
		Stmt Statement
	}

	CastExpression struct {
		ExprNode
		Expr     Expr
		DataType string
	}

	BindNameExpression struct {
		ExprNode
		Name Expr
	}

	ForUpdateOptionsExpression struct {
		ExprNode
		SkipLocked bool
		NoWait     bool
		Wait       Expr
	}

	NamedArgumentExpression struct {
		ExprNode
		Name  Expr
		Value Expr
	}

	ListaggExpression struct {
		ExprNode
		Args   []Expr
		Within Expr
		Over   Expr
	}

	OrderByClause struct {
		ExprNode
		Siblings bool
		Elements []Expr
	}

	OrderByElement struct {
		ExprNode
		Desc bool
		Item Expr
	}

	ExprListExpression struct {
		ExprNode
		Exprs []Expr
	}

	UsingElement struct {
		ExprNode
		IsIn  bool
		IsOut bool
		Elem  Expr
	}

	CommonTableExpression struct {
		ExprNode
		Name        Expr
		Query       *StatementExpression
		ColNameList []Expr
		IsRecursive bool
	}
)

func (n *ExprNode) expr() {}
