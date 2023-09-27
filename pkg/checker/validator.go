package checker

import (
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
		v.err = multierror.Append(v.err, e)
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
	// registerValidateRulesCtx(v, rules)
	return &SqlValidator{}
}
