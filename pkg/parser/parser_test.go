package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testSqlParser struct {
		SqlParser
		t *testing.T
	}
)

func TestParseProc(t *testing.T) {
	text := `create or replace package test is
	end;`

	parser := &testSqlParser{
		SqlParser{Buffer: text},
		t,
	}
	err := parser.Init()
	assert.Nil(t, err)
	assert.Nil(t, parser.Parse())
	ast := parser.AST()
	node := ast.up
	for node != nil {
		err := node.Accept(parser)
		assert.Nil(t, err)
		node = node.next
	}
}

func (p *testSqlParser) VisitCreatePackageDeclaration(node *node32) error {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulePackageName:
			assert.Equal(p.t, "test", strings.TrimSpace(string(p.buffer[node.begin:node.end])))
			return nil
		case ruleCREATE:
			fallthrough
		case ruleOR:
			fallthrough
		case ruleREPLACE:
			fallthrough
		case rulePACKAGE:
			fallthrough
		case ruleIS:
			fallthrough
		case ruleEND:
			fallthrough
		case ruleSEMI:
		default:
			assert.Failf(p.t, "", "Unexpected rule: %d", node.pegRule)
			return fmt.Errorf("unexpected rule: %d", node.pegRule)
		}
		node = node.next
	}
	return nil
}
