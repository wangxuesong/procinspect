package parser

import (
	"fmt"
)

type (
	sqlVisitor interface {
		VisitStatement(node *node32) error
		VisitCreatePackageDeclaration(node *node32) error
		VisitCreateProcedureDeclaration(node *node32) error
	}

	acceptor interface {
		Accept(visitor sqlVisitor) error
	}
)

func (n *node32) Accept(visitor sqlVisitor) error {
	switch n.pegRule {
	case ruleStatement:
		return visitor.VisitStatement(n)
	case ruleCreatePackageDeclaration:
		return visitor.VisitCreatePackageDeclaration(n)
	case ruleCreateProcedureDeclaration:
		return visitor.VisitCreateProcedureDeclaration(n)
	default:
		return fmt.Errorf("unexpected rule: %d", n.pegRule)
	}
}
