package semantic

type (
	StatementDepth interface {
		Get() int64
		Set(int64)
	}

	blockDepth struct {
		depth int64
	}

	AssignmentStatement struct {
		SyntaxNode
		Left  string
		Right Expr
	}

	BlockStatement struct {
		SyntaxNode
		Declarations []Declaration
		Body         *Body
	}

	IfStatement struct {
		SyntaxNode
		blockDepth
		Condition Expr
		ThenBlock []Statement
		ElseBlock []Statement
		ElseIfs   []*IfStatement
	}

	ElseBlock struct {
		SyntaxNode
		blockDepth
		Statements []Statement
	}

	LoopStatement struct {
		SyntaxNode
		blockDepth
		Statements []Statement
	}

	OpenStatement struct {
		SyntaxNode
		Name string
	}

	OpenForStatement struct {
		SyntaxNode
		Name  Expr
		For   Expr
		Using Expr
	}

	CloseStatement struct {
		SyntaxNode
		Name string
	}

	FetchStatement struct {
		SyntaxNode
		Cursor string
		Into   string
	}

	ExitStatement struct {
		SyntaxNode
		Condition Expr
	}

	ReturnStatement struct {
		SyntaxNode
		Name Expr
	}

	NullStatement struct {
		SyntaxNode
	}

	ProcedureCall struct {
		SyntaxNode
		Name      Expr
		Arguments []Expr
	}

	ExecuteImmediateStatement struct {
		SyntaxNode
		Sql   string
		Into  *IntoClause
		Using *UsingClause
	}

	IntoClause struct {
		SyntaxNode
		IsBulk bool
		Vars   []Expr
	}

	UsingClause struct {
		ExprNode
		WildCard *string
		Elems    []Expr
	}

	Declaration interface {
		Node
		declaration()
	}

	VariableDeclaration struct {
		SyntaxNode
		Name           string
		DataType       string
		Initialization Expr
	}

	ExceptionDeclaration struct {
		SyntaxNode
		Name string
	}

	CursorDeclaration struct {
		SyntaxNode
		Name        string
		Parameters  []*Parameter
		Stmt        Statement
		Return      string
		IsReference bool
	}

	NestTableTypeDeclaration struct {
		SyntaxNode
		Name string
	}

	FunctionDeclaration struct {
		SyntaxNode
		Name       string
		Parameters []*Parameter
	}

	Parameter struct {
		SyntaxNode
		Name     string
		DataType string
	}

	Argument struct {
		SyntaxNode
		Name string
	}

	AutonomousTransactionDeclaration struct {
		SyntaxNode
	}

	RaiseStatement struct {
		SyntaxNode
		Name string
	}

	GotoStatement struct {
		SyntaxNode
		Label string
	}

	LabelDeclaration struct {
		SyntaxNode
		Label string
	}
)

func (i *blockDepth) Get() int64 {
	return i.depth
}

func (i *blockDepth) Set(depth int64) {
	i.depth = depth
}

func (s *AssignmentStatement) Type() NodeType {
	return Assignment
}

func (s *AssignmentStatement) statement() {}

func (b *BlockStatement) statement() {}

func (d *VariableDeclaration) declaration() {}

func (d *ExceptionDeclaration) declaration() {}

func (d *CursorDeclaration) declaration() {}

func (d *NestTableTypeDeclaration) declaration() {}

func (d *FunctionDeclaration) declaration() {}

func (d *AutonomousTransactionDeclaration) declaration() {}

func (i *IfStatement) statement() {}

func (l *LoopStatement) statement() {}

func (o *OpenStatement) statement() {}

func (o *OpenForStatement) statement() {}

func (c *CloseStatement) statement() {}

func (s *FetchStatement) statement() {}

func (s *ExitStatement) statement() {}

func (s *ProcedureCall) statement() {}

func (s *ReturnStatement) statement() {}

func (s *NullStatement) statement() {}

func (s *ExecuteImmediateStatement) statement() {}

func (s *RaiseStatement) statement() {}

func (s *GotoStatement) statement() {}

func (s *LabelDeclaration) statement() {}
