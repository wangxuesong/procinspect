package checker

import (
	"errors"
	"reflect"

	"github.com/hashicorp/go-multierror"

	"procinspect/pkg/semantic"
)

type (
	SqlValidator struct {
		err     *multierror.Error
		ruleMap map[reflect.Type]Rule
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
	if r, ok := v.ruleMap[t]; ok {
		n := node.(semantic.Node)
		e := r.CheckFunc(r, n)
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
	v := &SqlValidator{
		ruleMap: make(map[reflect.Type]Rule),
	}
	v.RegisterValidateRules(rules)
	return v
}

func (v *SqlValidator) RegisterValidateRules(rs []Rule) {
	for _, r := range rs {
		t := reflect.TypeOf(r.Target)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		v.ruleMap[t] = r
	}
}
