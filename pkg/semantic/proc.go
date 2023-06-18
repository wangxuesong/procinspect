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

	ProcedureCall struct {
		node
		Name      Expr
		Arguments []Expr
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
		Name       string
		Parameters []*Parameter
		Stmt       Statement
	}

	NestTableTypeDeclaration struct {
		node
		Name string
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

func (i *IfStatement) statement() {}

func (l *LoopStatement) statement() {}

func (o *OpenStatement) statement() {}

func (c *CloseStatement) statement() {}

func (s *FetchStatement) statement() {}

func (s *ExitStatement) statement() {}

func (s *ProcedureCall) statement() {}
