package checker

import (
	"errors"

	"github.com/hashicorp/go-multierror"

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
		return script, err
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
	var result *multierror.Error
	for _, child := range semantic.GetChildren(node) {
		e := child.Accept(v)
		if e != nil {
			result = multierror.Append(result, e)
		}
	}
	var errs *multierror.Error
	err = v.v.Validate(node)
	if errors.As(err, &errs) {
		result = multierror.Append(result, errs.ErrorOrNil())
	} else {
		result = multierror.Append(result, err)
	}
	return result.ErrorOrNil()
}

func (v *ValidVisitor) VisitScript(node *semantic.Script) error {
	var errs *multierror.Error
	for _, stmt := range node.Statements {
		s := stmt.(semantic.AstNode)
		err := s.Accept(v)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

func (v *ValidVisitor) Error() error {
	return v.v.Error()
}

func (v *ValidVisitor) RegisterValidateRules(r []Rule) {
	v.v.RegisterValidateRules(r)
}
