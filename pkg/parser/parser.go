package parser

import (
	"fmt"
)

type (
	sqlVisitor interface {
		VisitCreatePackageDeclaration(node *node32) error
	}

	acceptor interface {
		Accept(visitor sqlVisitor) error
	}
)

func (n *node32) Accept(visitor sqlVisitor) error {
	switch n.pegRule {
	case ruleCreatePackageDeclaration:
		return visitor.VisitCreatePackageDeclaration(n)
	default:
		return fmt.Errorf("unexpected rule: %d", n.pegRule)
	}
}
