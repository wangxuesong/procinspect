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
		node

		Table  string
		Schema string
	}

	SelectField struct {
		node
		WildCard *WildCardField
		Expr     Expr
	}

	FieldList struct {
		node
		Fields []*SelectField
	}

	TableRef struct {
		node
		Table string
	}

	FromClause struct {
		node
		TableRefs []*TableRef
	}

	ForUpdateClause struct {
		node
		Options Expr
	}

	Statement interface {
		Node
		statement()
	}

	SelectStatement struct {
		node
		Fields      *FieldList
		From        *FromClause
		Where       Expr
		ForUpdate   *ForUpdateClause
		SetOperator *SetOperator
	}

	SetOperationStatement struct {
		node
		SelectList []Statement
	}

	CreateTypeStatement struct {
		node
		Name string
	}

	CreateNestTableStatement struct {
		node
		CreateTypeStatement
	}

	CreateSynonymStatement struct {
		node
		Synonym  Expr
		Original Expr
	}

	CaseWhenStatement struct {
		node
		Expr        Expr
		WhenClauses []*CaseWhenBlock
		ElseClause  *CaseWhenBlock
	}

	CaseWhenBlock struct {
		node
		Condition Expr
		Expr      Expr
		Stmts     []Statement
	}

	CommitStatement struct {
		node
	}

	RollbackStatement struct {
		node
	}

	ContinueStatement struct {
		node
	}

	DeleteStatement struct {
		node
		Table Expr
		Where Expr
	}

	UpdateStatement struct {
		node
		Table    Expr
		Where    Expr
		SetExprs []Expr
		SetValue Expr
	}

	InsertStatement struct {
		node
		AllInto []*InsertIntoClause
		Select  *SelectStatement
	}

	InsertIntoClause struct {
		node
		Table   *TableRef
		Columns []Expr
		Values  []Expr
	}

	MergeStatement struct {
		node
		Table       *TableRef
		Using       Expr
		OnCondition Expr
		MergeUpdate *MergeUpdateStatement
		MergeInsert *MergeInsertStatement
	}

	MergeUpdateStatement struct {
		node
		SetElems []Expr
		Where    Expr
		Delete   Expr
	}

	MergeInsertStatement struct {
		node
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
