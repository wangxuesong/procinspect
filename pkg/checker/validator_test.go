package checker

import (
	"testing"

	"github.com/hashicorp/go-multierror"
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
				assert.Nil(t, err)
				assert.NotNil(t, vv.Error())
				err = vv.Error()
				assert.IsType(t, &multierror.Error{}, err)
				var err1 *multierror.Error
				assert.ErrorAs(t, err, &err1)
				assert.NotNil(t, err1)
				assert.Equal(t, 1, len(err1.Errors))
				err2 := err1.Errors[0].(SqlValidationError)
				assert.Equal(t, 3, err2.Line)
				err = vv.Validate(stmt.Declarations[1])
				assert.Nil(t, err)
				assert.NotNil(t, stmt.Body)
				{
					assert.Equal(t, 2, len(stmt.Body.Statements))
					err = vv.Validate(stmt.Body.Statements[0])
					assert.Nil(t, err)
					err = vv.Validate(stmt.Body.Statements[1])
					errs := vv.Error().(*multierror.Error)
					assert.NotNil(t, errs)
					assert.IsType(t, SqlValidationError{}, errs.Errors[1])
					e1 := errs.Errors[1].(SqlValidationError)
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

func TestRuleEngine(t *testing.T) {
	tests := []validTestCase{
		{
			"rule engine",
			`select * from test@dblink;`,
			func(t *testing.T, root any) {
				require.IsType(t, &semantic.Script{}, root)
				node := root.(*semantic.Script)
				assert.Equal(t, len(node.Statements), 1)

				v := NewValidVisitor()
				node.Accept(v)
				assert.Nil(t, v.Error())
				r := rule{
					name:      "select from dblink",
					target:    &semantic.SelectStatement{},
					checkFunc: validDblinkFunc(`indexOf(node.From.TableRefs[0].Table, "@") > 0`),
					message:   "unsupported: select from dblink",
				}
				addRule(r)
				e := node.Accept(v)
				assert.Nil(t, e, func() string {
					if e != nil {
						return e.Error()
					}
					return ""
				}())
				assert.NotNil(t, v.Error())
				err := v.Error()
				assert.IsType(t, &multierror.Error{}, err)
				var err1 *multierror.Error
				assert.ErrorAs(t, err, &err1)
				assert.NotNil(t, err1)
				assert.Equal(t, 1, len(err1.Errors))
				err2 := err1.Errors[0].(SqlValidationError)
				assert.Equal(t, 1, err2.Line)
				assert.Equal(t, "unsupported: select from dblink", err2.Error())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runTest(t, test.text, test.Func)
		})
	}
}
