package interp

import (
	"errors"
	"procinspect/pkg/semantic"
)

type (
	Callable interface {
		Arity() int
		Call(i *Interpreter, arguments []any) (result any, err error)
		String() string
	}

	Gettable interface {
		Get(name string) (value any, err error)
	}

	Program struct {
		Script     *semantic.Script
		Procedures []*Procedure
		Statements []semantic.Stmt
		Packages   []*Package
	}

	Procedure struct {
		Name string
		Proc *semantic.CreateProcedureStatement
	}

	Package struct {
		Name    string
		Package *semantic.CreatePackageStatement
		Body    *semantic.CreatePackageBodyStatement

		procedures map[string]*Procedure
	}
)

func (p *Package) Get(name string) (any, error) {
	proc, ok := p.procedures[name]
	if ok {
		return proc, nil
	}

	if p.procedures == nil {
		p.procedures = make(map[string]*Procedure)
	}

	if p.Body != nil {
		for _, procedure := range p.Body.Procedures {
			if procedure.Name == name {
				proc = &Procedure{Name: name, Proc: procedure}
				p.procedures[name] = proc
				return proc, nil
			}
		}
	}
	return nil, errors.New("procedure " + name + " not found")
}

func (p *Procedure) Arity() int {
	return len(p.Proc.Parameters)
}

func (p *Procedure) Call(i *Interpreter, arguments []any) (result any, err error) {
	env := i.beginScope()
	defer i.endScope(env)

	for i, param := range p.Proc.Parameters {
		env.Define(param.Name, arguments[i])
	}

	// define variables
	for _, decl := range p.Proc.Declarations {
		switch decl.(type) {
		case *semantic.VariableDeclaration:
			v := decl.(*semantic.VariableDeclaration)
			var value any
			value, err = v.Initialization.(semantic.Expression).Accept(i)
			env.Define(v.Name, value)
		}
	}

	err = p.Proc.Body.Accept(i)
	return
}

func (p *Procedure) String() string {
	return "Procedure " + p.Name
}
