package parser

import (
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
	text := `select 1 from dual`

	p := plsql.NewParser(text)
	root := p.Sql_script()
	assert.Nil(t, p.Error())
	assert.NotNil(t, root)
	assert.IsType(t, &plsql.Sql_scriptContext{}, root)
	_, ok := root.(*plsql.Sql_scriptContext)
	assert.True(t, ok)
}
