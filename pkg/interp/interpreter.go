package interp

import (
	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	Interpreter struct {
		Source      string
		environment *Environment
		global      *Environment
		program     *Program
	}
)

func NewInterpreter() *Interpreter {
	env := NewGlobalEnvironment()
	return &Interpreter{
		environment: env,
		global:      env,
	}
}

func (interp *Interpreter) LoadScript(src string) (*Program, error) {
	script, err := parser.ParseScript(src)
	if err != nil {
		return nil, err
	}

	return interp.CompileAst(script)
}

func (interp *Interpreter) CompileAst(script *semantic.Script) (*Program, error) {
	interp.program = &Program{
		Script: script,
	}

	for _, stmt := range script.Statements {
		s := stmt.(semantic.Stmt)
		visitor := &StmtVisitor{interp: interp}
		err := s.Accept(visitor)
		if err != nil {
			return nil, err
		}
	}

	return interp.program, nil
}

func (interp *Interpreter) compileCreateProcedure(s *semantic.CreateProcedureStatement) Procedure {
	procedure := Procedure{
		Name: s.Name,
	}
	// compile declaration

	// compile body
	return procedure
}
