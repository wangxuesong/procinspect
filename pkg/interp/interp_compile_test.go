package interp

import (
	"github.com/stretchr/testify/assert"
	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
	"testing"
)

type (
	testCase struct {
		name string
		text string
		Func testCaseFunc
	}

	testCaseFunc func(*testing.T, *Interpreter)

	testSuite []testCase
)

func runTestSuite(t *testing.T, tests testSuite) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := NewInterpreter()
			i.Source = test.text
			test.Func(t, i)
		})
	}
}

func TestInterpreter_LoadScript(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "create simple procedure",
		text: `create or replace procedure test is
			t integer;
		begin
			t := 1;
			t := t + 1;
		end;`,
		Func: func(t *testing.T, i *Interpreter) {
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Procedures), 1)
			assert.Equal(t, program.Procedures[0].Name, "test")
			assert.Equal(t, len(i.global.values), 1)
			v, err := i.global.Get("test")
			assert.Nil(t, err)
			assert.IsType(t, v, &Procedure{})
		},
	})

	runTestSuite(t, tests)
}

func compileBlock(src string) (*semantic.Script, error) {
	script, err := parser.ParseBlock(src)
	if err != nil {
		return nil, err
	}
	return script, nil
}

func TestInterpreter_CompileBlock(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "compile block",
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
		},
	})

	runTestSuite(t, tests)
}

func TestInterpreter_CompileProcedure(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "compile simple procedure",
		text: `create or replace procedure test is
			t integer;
		begin
			t := 1;
			t := t + 1;
		end;`,
		Func: func(t *testing.T, i *Interpreter) {
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Procedures), 1)
			assert.Equal(t, program.Procedures[0].Name, "test")
			assert.Equal(t, len(i.global.values), 1)
			v, err := i.global.Get("test")
			assert.Nil(t, err)
			assert.IsType(t, v, &Procedure{})
			assert.NotNil(t, v.(*Procedure).Proc)
		},
	})

	tests = append(tests, testCase{
		name: "compile procedure with parameters",
		text: `create or replace procedure test (a NUMBER, b NUMBER default 100) is
			t integer;
		begin
			t := 1;
			t := t + 1;
		end;`,
		Func: func(t *testing.T, i *Interpreter) {
			program, err := i.LoadScript(i.Source)
			assert.Nil(t, err)
			assert.NotNil(t, program)
			assert.Equal(t, len(program.Procedures), 1)
			assert.Equal(t, program.Procedures[0].Name, "test")
			assert.Equal(t, len(i.global.values), 1)
			v, err := i.global.Get("test")
			assert.Nil(t, err)
			assert.IsType(t, v, &Procedure{})
			assert.NotNil(t, v.(*Procedure).Proc)
		},
	})

	runTestSuite(t, tests)

}
