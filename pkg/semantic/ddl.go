package semantic

type (
	CreatePackageStatement struct {
		node
		Name       string
		Procedures []*CreateProcedureStatement
		Types      []Declaration
		Variables  []Declaration
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

	DropFunctionStatement struct {
		node
		Name string
	}

	DropProcedureStatement struct {
		node
		Name string
	}

	DropPackageStatement struct {
		node
		Name   string
		Schema string
		IsBody bool
	}

	DropTriggerStatement struct {
		node
		Name string
	}
)

func (s *CreateProcedureStatement) Type() NodeType {
	return CreateProcedure
}

func (s *CreateProcedureStatement) statement() {}

func (s *CreatePackageStatement) statement() {}

func (s *CreatePackageBodyStatement) statement() {}

func (s *CreateFunctionStatement) statement() {}

func (s *Body) statement() {}

func (s *DropFunctionStatement) statement() {}

func (s *DropProcedureStatement) statement() {}

func (s *DropPackageStatement) statement() {}

func (s *DropTriggerStatement) statement() {}
