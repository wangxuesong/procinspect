package interp

import "procinspect/pkg/semantic"

type (
	StmtVisitor struct {
		semantic.StubExprVisitor
		interp *Interpreter
	}
)

func (v *StmtVisitor) VisitCreateProcedureStatement(s *semantic.CreateProcedureStatement) (err error) {
	proc := &Procedure{Name: s.Name}
	v.interp.program.Procedures = append(v.interp.program.Procedures, proc)
	v.interp.environment.Define(s.Name, proc)
	return
}
