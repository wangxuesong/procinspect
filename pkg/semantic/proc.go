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
)

func (s *AssignmentStatement) Type() NodeType {
	return Assignment
}

func (s *AssignmentStatement) statement() {}

func (d *VariableDeclaration) declaration() {}

func (d *ExceptionDeclaration) declaration() {}
