package interp

import "procinspect/pkg/semantic"

type (
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
