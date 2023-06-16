package checker

import (
	"github.com/stretchr/testify/assert"
	"procinspect/pkg/semantic"
	"testing"
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
	var tests testSuite

	tests = append(tests, testCase{
		name: "create nest table",
		text: `CREATE OR REPLACE TYPE NTHIS."DATA_ROW" as TABLE OF data_object;`,
		Func: func(t *testing.T, src string) {
			script, err := LoadScript(src)
			assert.Nil(t, err)
			assert.NotNil(t, script)
			assert.Equal(t, len(script.Statements), 1)
			assert.IsType(t, script.Statements[0], &semantic.CreateNestTableStatement{})
			stmt := script.Statements[0].(*semantic.CreateNestTableStatement)
			assert.Equal(t, "NTHIS.\"DATA_ROW\"", stmt.Name)
		},
	})

	runTestSuite(t, tests)
}
