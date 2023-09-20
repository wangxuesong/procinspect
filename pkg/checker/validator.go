package checker

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"

	"procinspect/pkg/semantic"
)

type (
	SqlValidator struct {
		v   *validator.Validate
		err error
	}
)

func (v *SqlValidator) Validate(node semantic.Node) error {
	err := v.v.StructCtx(context.Background(), node)
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		var errs SqlValidationErrors
		for _, er := range vErrs {
			node := er.Value().(semantic.Node)
			sqlErr := SqlValidationError{
				Line: node.Line(),
				Msg:  er.Param(),
			}
			errs = append(errs, sqlErr)
		}
		v.err = errs
	} else {
		v.err = err
	}
	return v.err
}

// NewSqlValidator creates a new instance of SqlValidator.
//
// It initializes a validator and registers the validate rules context.
// It returns a pointer to the newly created SqlValidator.
func NewSqlValidator() *SqlValidator {
	v := validator.New()
	registerValidateRulesCtx(v, rules)
	return &SqlValidator{
		v: v,
	}
}

// registerValidateRulesCtx registers the validation rules with the given validator instance.
//
// It takes a pointer to a validator.Validate struct and a slice of rule structs as parameters.
// Each rule struct contains a validation function and a target struct.
// The function iterates over the rules and registers the validation function with the target struct.
func registerValidateRulesCtx(v *validator.Validate, rules []rule) {
	for _, rule := range rules {
		v.RegisterStructValidationCtx(rule.validFunc(rule), rule.target)
	}
}
