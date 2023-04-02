package semantic

type (
	WildCardField struct {
		node

		Table  string
		Schema string
	}

	SelectField struct {
		node
		WildCard *WildCardField
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

	Statement interface {
		Node
		statement()
	}

	SelectStatement struct {
		node
		Fields *FieldList
		From   *FromClause
	}
)

func (s *SelectStatement) Type() NodeType {
	return StatementSelect
}

func (s *SelectStatement) statement() {}
