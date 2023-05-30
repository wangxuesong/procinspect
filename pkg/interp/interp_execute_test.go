package interp

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterpreter_ExecuteBlock(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "execute block",
		text: `
DECLARE
	a NUMBER;
	b NUMBER := 2;
BEGIN
	a:=1;
END`,
		Func: func(t *testing.T, i *Interpreter) {
			script, err := compileBlock(i.Source)
			assert.Nil(t, err)
			program, err := i.CompileAst(script)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Statements), 1)
			assert.Equal(t, len(i.global.values), 2)
			v, err := i.global.Get("a")
			assert.Nil(t, err)
			assert.Nil(t, v)
			v, err = i.global.Get("b")
			assert.Nil(t, err)
			assert.Nil(t, v)

			ctx := context.Background()
			err = i.Interpret(ctx, program)
			assert.Nil(t, err)
			assert.Equal(t, len(i.global.values), 2)
			v, err = i.global.Get("a")
			assert.Nil(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, v, &Number{Value: 1})
			v, err = i.global.Get("b")
			assert.Nil(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, v, &Number{Value: 2})

		},
	})

	runTestSuite(t, tests)
}

func TestInterpreter_ExecuteAnonymousBlock(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "execute anonymous block",
		text: `
DECLARE
	a NUMBER;
	b NUMBER := 2;
BEGIN
	a:=1;
END;`,
		Func: func(t *testing.T, i *Interpreter) {
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Statements), 1)
			assert.Equal(t, len(i.global.values), 2)
			v, err := i.global.Get("a")
			assert.Nil(t, err)
			assert.Nil(t, v)
			v, err = i.global.Get("b")
			assert.Nil(t, err)
			assert.Nil(t, v)

			ctx := context.Background()
			err = i.Interpret(ctx, program)
			assert.Nil(t, err)
			assert.Equal(t, len(i.global.values), 2)
			v, err = i.global.Get("a")
			assert.Nil(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, v, &Number{Value: 1})
			v, err = i.global.Get("b")
			assert.Nil(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, v, &Number{Value: 2})

		},
	})

	runTestSuite(t, tests)
}

type fooProcedure struct {
	value int64
}

func (f *fooProcedure) Arity() int {
	return 1
}

func (f *fooProcedure) Call(i *Interpreter, arguments []any) (result any, err error) {
	f.value = arguments[0].(*Number).Value
	result = f.value
	return
}

func (f *fooProcedure) String() string {
	return "foo procedure for testing"
}

func TestInterpreter_ExecuteCallProcedure(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "call native procedure",
		text: `
BEGIN
	foo(2);
END;`,
		Func: func(t *testing.T, i *Interpreter) {
			i.environment.Define("foo", &fooProcedure{})
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Statements), 1)
			assert.Equal(t, len(i.global.values), 1)

			ctx := context.Background()
			err = i.Interpret(ctx, program)
			assert.Nil(t, err)
			assert.Equal(t, len(i.global.values), 1)
			foo, ok := i.global.values["foo"].(*fooProcedure)
			assert.True(t, ok)
			assert.Equal(t, foo.value, int64(2))

		},
	})

	tests = append(tests, testCase{
		name: "call procedure",
		text: `
create procedure swth as
BEGIN
	foo(2);
END;
BEGIN
	swth();
END;`,
		Func: func(t *testing.T, i *Interpreter) {
			i.environment.Define("foo", &fooProcedure{})
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Statements), 1)
			assert.Equal(t, len(i.global.values), 2)

			ctx := context.Background()
			err = i.Interpret(ctx, program)
			assert.Nil(t, err)
			assert.Equal(t, len(i.global.values), 2)
			foo, ok := i.global.values["foo"].(*fooProcedure)
			assert.True(t, ok)
			assert.Equal(t, foo.value, int64(2))

		},
	})

	tests = append(tests, testCase{
		name: "call procedure with arguments",
		text: `
create procedure swth(a NUMBER) as
BEGIN
	foo(a);
END;
BEGIN
	swth(11);
END;`,
		Func: func(t *testing.T, i *Interpreter) {
			i.environment.Define("foo", &fooProcedure{})
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Statements), 1)
			assert.Equal(t, len(i.global.values), 2)

			ctx := context.Background()
			err = i.Interpret(ctx, program)
			assert.Nil(t, err)
			assert.Equal(t, len(i.global.values), 2)
			foo, ok := i.global.values["foo"].(*fooProcedure)
			assert.True(t, ok)
			assert.Equal(t, foo.value, int64(11))

		},
	})

	runTestSuite(t, tests)
}
