package checker

import (
	"errors"

	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	ValidVisitor struct {
		semantic.StubNodeVisitor
		v *SqlValidator
	}
)

func LoadScript(src string) (*semantic.Script, error) {
	script, err := parser.ParseScript(src)
	if err != nil {
		return nil, err
	}

	return script, nil
}

func NewValidVisitor() *ValidVisitor {
	return &ValidVisitor{
		v: NewSqlValidator(),
	}
}

func (v *ValidVisitor) VisitScript(node *semantic.Script) error {
	var errs []error
	for _, stmt := range node.Statements {
		s := stmt.(semantic.AstNode)
		err := s.Accept(v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (v *ValidVisitor) VisitCreateNestTableStatement(node *semantic.CreateNestTableStatement) error {
	return v.v.Validate(node)
}

func (v *ValidVisitor) VisitSelectStatement(node *semantic.SelectStatement) error {
	return v.v.Validate(node)
}
