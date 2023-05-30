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

func (p *Procedure) Arity() int {
	return len(p.Proc.Parameters)
}

func (p *Procedure) Call(i *Interpreter, arguments []any) (result any, err error) {
	err = p.Proc.Body.Accept(i)
	return
}

func (p *Procedure) String() string {
	return "Procedure " + p.Name
}
