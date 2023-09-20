package checker

import (
	"context"
	"fmt"
	"testing"

	validator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"procinspect/pkg/semantic"
)

type (
	validTestCase struct {
		name string
		text string
		Func validTestCaseFunc
	}

	validTestCaseFunc func(*testing.T, any)
)

// runTest is a function that runs a test for a given input and test function.
//
// It takes the following parameters:
//   - t: a *testing.T value, used for testing purposes.
//   - input: a string containing the input for the test.
//   - testFunc: a function that takes a *testing.T value and a node value of any type,
//     used to perform the test.
//   - rootFunc: an optional variadic parameter that represents a root function.
//
// The function does not return any value.
func runTest(t *testing.T, input string, testFunc func(t *testing.T, node any)) {
	node, err := LoadScript(input)
	assert.Nil(t, err)
	testFunc(t, node)
	return
}

func TestSimpleValidator(t *testing.T) {
	tests := []validTestCase{
		{
			name: "simple",
			text: `
declare
	type t is table of number index by binary_integer;
	a t;
begin
	a := 1;
    update t set (a,b) = (select 1, 2 from dual);
end;`,
			Func: func(t *testing.T, root any) {
				require.IsType(t, &semantic.Script{}, root)
				node := root.(*semantic.Script)
				assert.Greater(t, len(node.Statements), 0)
				assert.IsType(t, &semantic.BlockStatement{}, node.Statements[0])
				stmt := node.Statements[0].(*semantic.BlockStatement)
				assert.IsType(t, &semantic.NestTableTypeDeclaration{}, stmt.Declarations[0])
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[1])

				vv := NewSqlValidator()
				err := vv.Validate(stmt)
				assert.Nil(t, err)
				err = vv.Validate(stmt.Declarations[0])
				assert.NotNil(t, err)
				assert.IsType(t, SqlValidationErrors{}, err)
				err1 := err.(SqlValidationErrors)[0]
				assert.Equal(t, 3, err1.Line)
				err = vv.Validate(stmt.Declarations[1])
				assert.Nil(t, err)
				assert.NotNil(t, stmt.Body)
				{
					assert.Equal(t, 2, len(stmt.Body.Statements))
					err = vv.Validate(stmt.Body.Statements[0])
					assert.Nil(t, err)
					err = vv.Validate(stmt.Body.Statements[1])
					assert.NotNil(t, err)
					assert.IsType(t, SqlValidationErrors{}, err)
					e1 := err.(SqlValidationErrors)[0]
					assert.Error(t, e1)
					assert.Equal(t, 7, e1.Line)
					assert.Equal(t, "unsupported: update set multiple columns with select", e1.Error())

				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runTest(t, test.text, test.Func)
		})
	}

}

func ValidateSelectStatement(ctx context.Context, sl validator.StructLevel) {
	stmt, ok := sl.Current().Interface().(semantic.FromClause)
	if !ok {
		errors := ctx.Value("errors").(map[string]error)
		errors["err"] = fmt.Errorf("is not a SelectStatement")
		return
	}
	from := stmt
	if len(from.TableRefs) > 1 {
		for _, table := range from.TableRefs {
			if table.Table == "dual" {
				sl.ReportError(stmt, "From", "From", "fromdualwithtable", "fromdualwithtable")
			}
		}
	}
}
