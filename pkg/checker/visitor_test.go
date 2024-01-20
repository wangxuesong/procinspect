package checker

import (
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	testCase struct {
		name string
		text string
		Func testCaseFunc
	}

	testCaseFunc func(*testing.T, string)

	testSuite []testCase
)

func runTestSuite(t *testing.T, tests testSuite) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.Func(t, test.text)
		})
	}
}

func TestCheckCreateNestTable(t *testing.T) {
	tests := []testCase{
		{
			name: "create nest table",
			text: `CREATE OR REPLACE TYPE NTHIS."DATA_ROW" as TABLE OF data_object;
select * from NTHIS."DATA_ROW";
with sql1 as (
select * from test@dblink
)
select * from sql1;`,
			Func: func(t *testing.T, src string) {
				script, err := LoadScript(src)
				assert.Nil(t, err)
				assert.NotNil(t, script)
				v := NewValidVisitor()
				_ = script.Accept(v)
				err = v.Error()
				assert.NotNil(t, err)
				assert.IsType(t, &multierror.Error{}, err)
				errs := err.(*multierror.Error).Errors
				require.Equal(t, 1, len(errs))
				var target = SqlValidationError{}
				assert.ErrorAs(t, errs[0], &target)
				sqlErr := target
				assert.ErrorIs(t, sqlErr, SqlValidationError{
					Line: 1,
					Msg:  "unsupported: nest table type declaration",
				})
				// assert.ErrorAs(t, errs[1], &target)
				// sqlErr = target
				// assert.ErrorIs(t, sqlErr, SqlValidationError{
				// 	Line: 0,
				// 	Msg:  "unsupported: select from dblink",
				// })
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.Func(t, test.text)
		})
	}
}
