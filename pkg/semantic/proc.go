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
		node
		Left  string
		Right Expr
	}

	BlockStatement struct {
		node
		Declarations []Declaration
		Body         *Body
	}

	IfStatement struct {
		node
		blockDepth
		Condition Expr
		ThenBlock []Statement
		ElseBlock []Statement
		ElseIfs   []*IfStatement
	}

	ElseBlock struct {
		node
		blockDepth
		Statements []Statement
	}

	LoopStatement struct {
		node
		blockDepth
		Statements []Statement
	}

	OpenStatement struct {
		node
		Name string
	}

	OpenForStatement struct {
		node
		Name  Expr
		For   Expr
		Using Expr
	}

	CloseStatement struct {
		node
		Name string
	}

	FetchStatement struct {
		node
		Cursor string
		Into   string
	}

	ExitStatement struct {
		node
		Condition Expr
	}

	ReturnStatement struct {
		node
		Name Expr
	}

	NullStatement struct {
		node
	}

	ProcedureCall struct {
		node
		Name      Expr
		Arguments []Expr
	}

	ExecuteImmediateStatement struct {
		node
		Sql string
	}

	Declaration interface {
		Node
		declaration()
	}

	VariableDeclaration struct {
		node
		Name           string
		DataType       string
		Initialization Expr
	}

	ExceptionDeclaration struct {
		node
		Name string
	}

	CursorDeclaration struct {
		node
		Name        string
		Parameters  []*Parameter
		Stmt        Statement
		Return      string
		IsReference bool
	}

	NestTableTypeDeclaration struct {
		node
		Name string
	}

	FunctionDeclaration struct {
		node
		Name       string
		Parameters []*Parameter
	}

	Parameter struct {
		node
		Name     string
		DataType string
	}

	Argument struct {
		node
		Name string
	}

	AutonomousTransactionDeclaration struct {
		node
	}

	RaiseStatement struct {
		node
		Name string
	}

	GotoStatement struct {
		node
		Label string
	}

	LabelDeclaration struct {
		node
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
