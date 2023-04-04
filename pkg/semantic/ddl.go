package semantic

type (
	CreateProcedureStatement struct {
		node
		Name         string
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
