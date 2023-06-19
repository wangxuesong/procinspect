package interp

import (
	"fmt"
)

// Environment is the environment of the interpreter.
// storage of variables.
type Environment struct {
	values map[string]interface{}
	parent *Environment
}

func NewGlobalEnvironment() (env *Environment) {
	env = &Environment{
		values: make(map[string]interface{}),
	}

	return
}

func NewChildEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]interface{}),
		parent: parent,
	}
}

// Define a variable in the environment.
func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

// Get a variable from the environment.
func (e *Environment) Get(name string) (value any, err error) {
	value, ok := e.values[name]
	if ok {
		return
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	err = fmt.Errorf("undefined variable %q", name)
	return
}

// Assign a variable in the environment with the given value.
func (e *Environment) Assign(name string, value interface{}) (err error) {
	_, ok := e.values[name]
	if ok {
		e.values[name] = value
		return
	}

	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	err = fmt.Errorf("can not assign undecleared variable %q", name)
	return
}

func (e *Environment) GetAt(distance int, name string) (result interface{}, err error) {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
		if env == nil {
			err = fmt.Errorf("non-existed env parent, searching for variable %q, want distance %d, current distance %d", name, distance, i)
			return
		}
	}

	result, ok := env.values[name]
	if !ok {
		err = fmt.Errorf("unexpected resolving on variable %q", name)
		return
	}

	return
}

func (e *Environment) AssignAt(distance int, name string, result interface{}) (err error) {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
		if env == nil {
			err = fmt.Errorf("non-existed env parent, searching for variable %q, want distance %d, current distance %d", name, distance, i)
			return
		}
	}

	env.values[name] = result
	return
}
