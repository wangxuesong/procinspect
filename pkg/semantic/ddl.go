package semantic

type (
	CreatePackageStatement struct {
		node
		Name       string
		Procedures []*CreateProcedureStatement
		Types      []Declaration
	}

	CreatePackageBodyStatement struct {
		node
		Name       string
		Procedures []*CreateProcedureStatement
		Functions  []*CreateFunctionStatement
	}

	CreateProcedureStatement struct {
		node
		Name         string
		Parameters   []*Parameter
		Declarations []Declaration
		Body         *Body
		IsReplace    bool
	}

	CreateFunctionStatement struct {
		node
		Name         string
		Parameters   []*Parameter
		Return       string
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

func (s *CreateFunctionStatement) statement() {}
