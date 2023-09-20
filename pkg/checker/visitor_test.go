package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			text: `CREATE OR REPLACE TYPE NTHIS."DATA_ROW" as TABLE OF data_object;`,
			Func: func(t *testing.T, src string) {
				script, err := LoadScript(src)
				assert.Nil(t, err)
				assert.NotNil(t, script)
				v := NewValidVisitor()
				err = script.Accept(v)
				assert.NotNil(t, err)
				var target = &SqlValidationErrors{}
				assert.ErrorAs(t, err, target)
				sqlErr := (*target)[0]
				assert.ErrorIs(t, sqlErr, SqlValidationError{
					Line: 1,
					Msg:  "unsupported: nest table type declaration",
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.Func(t, test.text)
		})
	}
}
