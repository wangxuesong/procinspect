package semantic

type (
	CreatePackageStatement struct {
		node
		Name       string
		Procedures []*CreateProcedureStatement
	}

	CreatePackageBodyStatement struct {
		node
		Name       string
		Procedures []*CreateProcedureStatement
	}

	CreateProcedureStatement struct {
		node
		Name         string
		Parameters   []*Parameter
		Declarations []Declaration
		Body         *Body
		IsReplace    bool
	}

	Body struct {
		node
		Statements []Statement
	}
)

func (s *CreateProcedureStatement) Type() NodeType {
	return CreateProcedure
}

func (s *CreateProcedureStatement) statement() {}

func (s *CreatePackageStatement) statement() {}

func (s *CreatePackageBodyStatement) statement() {}
