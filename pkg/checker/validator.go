package checker

import (
	"errors"
	"reflect"

	"github.com/hashicorp/go-multierror"

	"procinspect/pkg/semantic"
)

type (
	SqlValidator struct {
		err *multierror.Error
	}

	Validator interface {
		Validate() error
	}

	ValidateFunc func() error
)

func (fn ValidateFunc) Validate() error {
	return fn()
}

func (v *SqlValidator) Validate(node semantic.AstNode) error {
	t := reflect.TypeOf(node)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if r, ok := ruleMap[t]; ok {
		n := node.(semantic.Node)
		e := r.checkFunc(r, n)
		var verr SqlValidationError
		if errors.As(e, &verr) {
			v.err = multierror.Append(v.err, verr)
		} else {
			return e
		}
	}

	return nil
}

func (v *SqlValidator) Error() error {
	return v.err
}

// NewSqlValidator creates a new instance of SqlValidator.
//
// It initializes a validator and registers the validate rules context.
// It returns a pointer to the newly created SqlValidator.
func NewSqlValidator() *SqlValidator {
	registerValidateRules(rules)
	return &SqlValidator{}
}

func registerValidateRules(rs []rule) {
	ruleMap = make(map[reflect.Type]rule)
	for _, r := range rs {
		addRule(r)
	}
}

func addRule(r rule) {
	t := reflect.TypeOf(r.target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ruleMap[t] = r
}
