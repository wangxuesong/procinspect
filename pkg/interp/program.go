package interp

import "procinspect/pkg/semantic"

type (
	Callable interface {
		Arity() int
		Call(i *Interpreter, arguments []any) (result any, err error)
		String() string
	}

	Program struct {
		Script     *semantic.Script
		Procedures []*Procedure
		Statements []semantic.Stmt
	}

	Procedure struct {
		Name string
		Proc *semantic.CreateProcedureStatement
	}
)
