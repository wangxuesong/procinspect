package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (c *SqlParser) visitCreatePackageDeclaration(t *testing.T, node *node32) {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case ruleIdentifier:
			assert.Equal(t, "test", strings.TrimSpace(string(c.buffer[node.begin:node.end])))
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
			assert.Failf(t, "", "Unexpected rule: %d", node.pegRule)
		}
		node = node.next
	}
}
