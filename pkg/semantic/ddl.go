package semantic

type (
	CreatePackageStatement struct {
		SyntaxNode
		Name       string
		Procedures []*CreateProcedureStatement
		Types      []Declaration
		Variables  []Declaration
	}

	CreatePackageBodyStatement struct {
		SyntaxNode
		Name       string
		Procedures []*CreateProcedureStatement
		Functions  []*CreateFunctionStatement
	}

	CreateProcedureStatement struct {
		SyntaxNode
		Name         string
		Parameters   []*Parameter
		Declarations []Declaration
		Body         *Body
		IsReplace    bool
	}

	CreateFunctionStatement struct {
		SyntaxNode
		Name         string
		Parameters   []*Parameter
		Return       string
		Declarations []Declaration
		Body         *Body
		IsReplace    bool
	}

	Body struct {
		SyntaxNode
		Statements []Statement
	}

	CreateTriggerStatement struct {
		SyntaxNode
		Name        string
		TriggerBody TriggerBody
	}

	TriggerBody interface {
		triggerBody()
	}

	TriggerBlock struct {
		SyntaxNode
		Declarations []Declaration
		Body         *Body
	}

	CreateSimpleDmlTriggerStatement struct {
		CreateTriggerStatement
		IsBefore   bool
		ForEachRow bool
		Events     []*TriggerEvent
		TableView  string
	}

	CreateCompoundDmlTriggerStatement struct {
		CreateTriggerStatement
		Events    []*TriggerEvent
		TableView string
	}

	TriggerEvent struct {
		SyntaxNode
		Name   string
		Column string
	}

	CompoundTriggerBlock struct {
		SyntaxNode
		Declarations []Declaration
		TimingPoints []*TimingPoint
	}

	TimingPoint struct {
		SyntaxNode
		IsBefore   bool
		ForEachRow bool
		Body       *Body
	}

	DropFunctionStatement struct {
		SyntaxNode
		Name string
	}

	DropProcedureStatement struct {
		SyntaxNode
		Name string
	}

	DropPackageStatement struct {
		SyntaxNode
		Name   string
		Schema string
		IsBody bool
	}

	DropTriggerStatement struct {
		SyntaxNode
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

func (s *CreateTriggerStatement) statement() {}

func (s *CreateSimpleDmlTriggerStatement) statement() {}

func (s *TriggerBlock) statement() {}

func (s *TriggerBlock) triggerBody() {}

func (s *CreateCompoundDmlTriggerStatement) statement() {}

func (s *CompoundTriggerBlock) statement() {}

func (s *CompoundTriggerBlock) triggerBody() {}

func (s *TimingPoint) statement() {}
