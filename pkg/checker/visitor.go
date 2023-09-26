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
	v := &ValidVisitor{
		v: NewSqlValidator(),
	}
	v.StubNodeVisitor.NodeVisitor = v
	return v
}

func (v *ValidVisitor) VisitChildren(node semantic.AstNode) (err error) {
	for _, child := range semantic.GetChildren(node) {
		e := child.Accept(v)
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return
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
	err := v.VisitChildren(node)
	return errors.Join(err, v.v.Validate(node))
}
