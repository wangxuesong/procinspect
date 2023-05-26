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
