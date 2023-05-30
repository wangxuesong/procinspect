package interp

import (
	"context"

	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	Interpreter struct {
		Source      string
		environment *Environment
		global      *Environment
		program     *Program

		semantic.StubExprVisitor
	}
)

func NewInterpreter() *Interpreter {
	env := NewGlobalEnvironment()
	return &Interpreter{
		environment: env,
		global:      env,
	}
}

func (i *Interpreter) LoadScript(src string) (*Program, error) {
	script, err := parser.ParseScript(src)
	if err != nil {
		return nil, err
	}

	return i.CompileAst(script)
}

func (i *Interpreter) CompileAst(script *semantic.Script) (*Program, error) {
	i.program = &Program{
		Script: script,
	}

	for _, stmt := range script.Statements {
		s := stmt.(semantic.Stmt)
		visitor := &resolver{interp: i}
		err := s.Accept(visitor)
		if err != nil {
			return nil, err
		}
	}

	return i.program, nil
}

func (i *Interpreter) Interpret(ctx context.Context, program *Program) (err error) {
	done := ctx.Done()
	for _, stmt := range program.Statements {
		select {
		case <-done:
			err = ctx.Err()
		default:
			// relax
		}
		err = i.execute(stmt)
		if err != nil {
			return
		}
	}

	return
}

func (i *Interpreter) beginScope() *Environment {
	env := NewChildEnvironment(i.environment)
	i.environment = env
	return env
}

func (i *Interpreter) endScope(env *Environment) {
	i.environment = env.parent
}

func (i *Interpreter) execute(stmt semantic.Stmt) (err error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitVariableDeclaration(s *semantic.VariableDeclaration) (err error) {
	var value any
	if s.Initialization != nil {
		value, err = s.Initialization.(semantic.Expression).Accept(i)
		if err != nil {
			return
		}
		err := i.environment.Assign(s.Name, value)
		if err != nil {
			return err
		}
	}
	i.environment.Define(s.Name, value)
	return
}

func (i *Interpreter) VisitBlockStatement(s *semantic.BlockStatement) (err error) {
	for _, decl := range s.Declarations {
		stmt := decl.(semantic.Stmt)
		err = stmt.Accept(i)
		if err != nil {
			return
		}
	}

	return s.Body.Accept(i)
}

func (i *Interpreter) VisitBody(s *semantic.Body) (err error) {
	for _, s := range s.Statements {
		stmt := s.(semantic.Stmt)
		err = stmt.Accept(i)
		if err != nil {
			return
		}
	}
	return
}

func (i *Interpreter) VisitAssignmentStatement(s *semantic.AssignmentStatement) (err error) {
	right := s.Right.(semantic.Expression)
	value, err := right.Accept(i)
	if err != nil {
		return err
	}
	err = i.environment.Assign(s.Left, value)
	return
}

func (i *Interpreter) VisitProcedureCall(s *semantic.ProcedureCall) (err error) {
	proc, err := i.environment.Get(s.Name)
	if err != nil {
		return err
	}

	// process arguments
	var arguments []any
	for _, arg := range s.Arguments {
		value, err := arg.(semantic.Expression).Accept(i)
		if err != nil {
			return err
		}
		arguments = append(arguments, value)
	}

	_, err = proc.(Callable).Call(i, arguments)
	if err != nil {
		return err
	}
	return
}

func (i *Interpreter) VisitNumericLiteral(s *semantic.NumericLiteral) (result any, err error) {
	number := &Number{}
	number.Value = s.Value
	return number, nil
}

func (i *Interpreter) VisitNameExpression(s *semantic.NameExpression) (result any, err error) {
	result, err = i.environment.Get(s.Name)
	return
}
