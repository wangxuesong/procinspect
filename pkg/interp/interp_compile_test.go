package interp

import (
	"github.com/stretchr/testify/assert"
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
