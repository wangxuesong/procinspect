package interp

import (
	"context"
	"fmt"

	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	Interpreter struct {
		Source      string
		environment *Environment
		global      *Environment
		program     *Program

		semantic.StubStmtVisitor
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
		err := s.StmtAccept(visitor)
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
	return stmt.StmtAccept(i)
}

func (i *Interpreter) evaluate(expr semantic.Expression) (any, error) {
	return expr.ExprAccept(i)
}

func (i *Interpreter) VisitVariableDeclaration(s *semantic.VariableDeclaration) (err error) {
	var value any
	if s.Initialization != nil {
		value, err = s.Initialization.(semantic.Expression).ExprAccept(i)
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
		err = stmt.StmtAccept(i)
		if err != nil {
			return
		}
	}

	return s.Body.StmtAccept(i)
}

func (i *Interpreter) VisitBody(s *semantic.Body) (err error) {
	for _, s := range s.Statements {
		stmt := s.(semantic.Stmt)
		err = stmt.StmtAccept(i)
		if err != nil {
			return
		}
	}
	return
}

func (i *Interpreter) VisitAssignmentStatement(s *semantic.AssignmentStatement) (err error) {
	right := s.Right.(semantic.Expression)
	value, err := right.ExprAccept(i)
	if err != nil {
		return err
	}
	err = i.environment.Assign(s.Left, value)
	return
}

func (i *Interpreter) VisitProcedureCall(s *semantic.ProcedureCall) (err error) {
	callee, err := i.evaluate(s.Name.(semantic.Expression))
	if err != nil {
		return err
	}

	// process arguments
	var arguments []any
	for _, arg := range s.Arguments {
		value, err := arg.(semantic.Expression).ExprAccept(i)
		if err != nil {
			return err
		}
		arguments = append(arguments, value)
	}

	callable := callee.(Callable)
	if want, got := callable.Arity(), len(arguments); want != got {
		err = fmt.Errorf("function expected %d arguments but got %d, at line %d", want, got, s.Line())
		return
	}

	_, err = callable.Call(i, arguments)
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

func (i *Interpreter) VisitDotExpression(s *semantic.DotExpression) (result any, err error) {
	expr := s.Parent
	if expr == nil {
		expr = s.Name
		instance, err := i.evaluate(expr.(semantic.Expression))
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
	instance, err := i.evaluate(expr.(semantic.Expression))
	if err != nil {
		return nil, err
	}

	gettable, ok := instance.(Gettable)
	if !ok {
		return nil, fmt.Errorf("instance is not gettable")
	}

	name := s.Name.(*semantic.NameExpression).Name
	return gettable.Get(name)
}
