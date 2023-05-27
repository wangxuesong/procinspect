package interp

import "procinspect/pkg/semantic"

type (
	resolver struct {
		semantic.StubExprVisitor
		interp *Interpreter
	}
)

func (v *resolver) VisitCreateProcedureStatement(s *semantic.CreateProcedureStatement) (err error) {
	proc := &Procedure{Name: s.Name, Proc: s}
	v.interp.program.Procedures = append(v.interp.program.Procedures, proc)
	v.interp.environment.Define(s.Name, proc)
	return
}

func (v *resolver) VisitBlockStatement(s *semantic.BlockStatement) (err error) {
	v.interp.program.Statements = append(v.interp.program.Statements, s)
	for _, decl := range s.Declarations {
		stmt := decl.(semantic.Stmt)
		err = stmt.Accept(v)
		if err != nil {
			return
		}
	}
	return
}

func (v *resolver) VisitVariableDeclaration(s *semantic.VariableDeclaration) (err error) {
	var value any
	if s.Initialization != nil {
		// TODO: support initializer
	}
	v.interp.environment.Define(s.Name, value)
	return
}
