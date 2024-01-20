package semantic

type SetOperator int

const (
	Union SetOperator = iota
	UnionAll
	Intersect
	Minus
)

type (
	WildCardField struct {
		SyntaxNode

		Table  string
		Schema string
	}

	SelectField struct {
		SyntaxNode
		WildCard *WildCardField
		Expr     Expr
	}

	FieldList struct {
		SyntaxNode
		Fields []*SelectField
	}

	TableRef struct {
		SyntaxNode
		Table string
	}

	FromClause struct {
		SyntaxNode
		TableRefs []*TableRef
	}

	ForUpdateClause struct {
		SyntaxNode
		Options Expr
	}

	Statement interface {
		Node
		statement()
	}

	SelectStatement struct {
		SyntaxNode
		Fields      *FieldList
		From        *FromClause
		Where       Expr
		ForUpdate   *ForUpdateClause
		SetOperator *SetOperator
		With        *WithClause
	}

	WithClause struct {
		SyntaxNode
		IsRecursive bool
		CTEs        []*CommonTableExpression
	}

	SetOperationStatement struct {
		SyntaxNode
		SelectList []Statement
	}

	CreateTypeStatement struct {
		SyntaxNode
		Name string
	}

	CreateNestTableStatement struct {
		SyntaxNode
		CreateTypeStatement
	}

	CreateSynonymStatement struct {
		SyntaxNode
		Synonym  Expr
		Original Expr
	}

	CaseWhenStatement struct {
		SyntaxNode
		Expr        Expr
		WhenClauses []*CaseWhenBlock
		ElseClause  *CaseWhenBlock
	}

	CaseWhenBlock struct {
		SyntaxNode
		Condition Expr
		Expr      Expr
		Stmts     []Statement
	}

	CommitStatement struct {
		SyntaxNode
	}

	RollbackStatement struct {
		SyntaxNode
	}

	ContinueStatement struct {
		SyntaxNode
	}

	DeleteStatement struct {
		SyntaxNode
		Table Expr
		Where Expr
	}

	UpdateStatement struct {
		SyntaxNode
		Table    Expr
		Where    Expr
		SetExprs []Expr
		SetValue Expr
	}

	InsertStatement struct {
		SyntaxNode
		AllInto []*InsertIntoClause
		Select  *SelectStatement
	}

	InsertIntoClause struct {
		SyntaxNode
		Table   *TableRef
		Columns []Expr
		Values  []Expr
	}

	MergeStatement struct {
		SyntaxNode
		Table       *TableRef
		Using       Expr
		OnCondition Expr
		MergeUpdate *MergeUpdateStatement
		MergeInsert *MergeInsertStatement
	}

	MergeUpdateStatement struct {
		SyntaxNode
		SetElems []Expr
		Where    Expr
		Delete   Expr
	}

	MergeInsertStatement struct {
		SyntaxNode
	}
)

func (s *SelectStatement) Type() NodeType {
	return StatementSelect
}

func (s *SelectStatement) statement() {}

func (s *SetOperationStatement) statement() {}

func (s *CreateTypeStatement) statement() {}

func (s *CreateNestTableStatement) statement() {}

func (s *CreateSynonymStatement) statement() {}

func (s *CaseWhenStatement) statement() {}

func (s *CommitStatement) statement() {}

func (s *RollbackStatement) statement() {}

func (s *ContinueStatement) statement() {}

func (s *DeleteStatement) statement() {}

func (s *UpdateStatement) statement() {}

func (s *InsertStatement) statement() {}

func (s *MergeStatement) statement() {}

func (s *MergeUpdateStatement) statement() {}

func (s *MergeInsertStatement) statement() {}
