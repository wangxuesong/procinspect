package semantic

type (
	Expr interface {
		Node
		expr()
	}

	expression struct {
		node
	}

	NumericLiteral struct {
		expression
		Value int64
	}
)

func (n *expression) expr() {}
