package parser

import (
	"procinspect/pkg/semantic"
	"testing"

	plsql "procinspect/pkg/parser/internal/plsql/parser"

	"github.com/stretchr/testify/assert"
)

func TestParseCreatePackage(t *testing.T) {
	text := `create or replace package test is
	end;`

	p := plsql.NewParser(text)
	root := p.Sql_script()
	assert.Nil(t, p.Error())
	assert.NotNil(t, root)
	assert.IsType(t, &plsql.Sql_scriptContext{}, root)
	_, ok := root.(*plsql.Sql_scriptContext)
	assert.True(t, ok)
}

func TestParseCreateProc(t *testing.T) {
	text := `create or replace procedure test is
	begin
		select 1 from dual;
	end;`

	p := plsql.NewParser(text)
	root := p.Sql_script()
	assert.Nil(t, p.Error())
	assert.NotNil(t, root)
	assert.IsType(t, &plsql.Sql_scriptContext{}, root)
	_, ok := root.(*plsql.Sql_scriptContext)
	assert.True(t, ok)
}

func TestParseSimple(t *testing.T) {
	text := `select * from dual, test;`

	p := plsql.NewParser(text)
	root := p.Sql_script()
	assert.Nil(t, p.Error())
	assert.NotNil(t, root)
	assert.IsType(t, &plsql.Sql_scriptContext{}, root)
	_, ok := root.(*plsql.Sql_scriptContext)
	assert.True(t, ok)
	node := GeneralScript(root)
	assert.NotNil(t, node)
	assert.Greater(t, len(node.Statements), 0)
	stmt, ok := node.Statements[0].(*semantic.SelectStatement)
	assert.True(t, ok)
	assert.NotNil(t, stmt)
	assert.Equal(t, len(stmt.Fields.Fields), 1)
	assert.Equal(t, stmt.Fields.Fields[0].WildCard.Table, "*")
	assert.Equal(t, len(stmt.From.TableRefs), 2)
	assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
	assert.Equal(t, stmt.From.TableRefs[1].Table, "test")
}
