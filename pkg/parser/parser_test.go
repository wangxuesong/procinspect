package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProc(t *testing.T) {
	text := `create or replace package test is
	end;`

	parser := &SqlParser{Buffer: text}
	err := parser.Init()
	assert.Nil(t, err)
	assert.Nil(t, parser.Parse())
	ast := parser.AST()
	node := ast.up
	for node != nil {
		switch node.pegRule {
		case ruleCreatePackageDeclaration:
			assert.Equal(t, 0, int(node.begin))
			parser.visitCreatePackageDeclaration(t, node)
			//return
		default:
			assert.Failf(t, "", "Unexpected rule: %d", node.pegRule)
		}
		node = node.next
	}
}
