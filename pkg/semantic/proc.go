package semantic

type (
	AssignmentStatement struct {
		node
		Left  string
		Right string
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

func (s *AssignmentStatement) Type() NodeType {
	return Assignment
}

func (s *AssignmentStatement) statement() {}

func (d *VariableDeclaration) declaration() {}

func (d *ExceptionDeclaration) declaration() {}

func (d *CursorDeclaration) declaration() {}
