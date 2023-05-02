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
		Right string
	}

	IfStatement struct {
		node
		blockDepth
		Condition string
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

	Declaration interface {
		Node
		declaration()
	}

	VariableDeclaration struct {
		node
		Name     string
		DataType string
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

	Parameter struct {
		node
		Name     string
		DataType string
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

func (d *VariableDeclaration) declaration() {}

func (d *ExceptionDeclaration) declaration() {}

func (d *CursorDeclaration) declaration() {}

func (i *IfStatement) statement() {}

func (l *LoopStatement) statement() {}

func (o *OpenStatement) statement() {}

func (c *CloseStatement) statement() {}
