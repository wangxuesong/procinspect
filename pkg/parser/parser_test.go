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

func TestCreateProcedure(t *testing.T) {
	text := `create or replace procedure test is
	begin
		select 1 from dual;
		select 2 from t;
	end;`

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

	{ // assert that the statement is a CreateProcedureStatement
		assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
		stmt := node.Statements[0].(*semantic.CreateProcedureStatement)

		assert.True(t, stmt.IsReplace)

		assert.Equal(t, stmt.Name, "test")

		assert.NotNil(t, stmt.Body)
		assert.Equal(t, len(stmt.Body.Statements), 2)

		assert.IsType(t, &semantic.SelectStatement{}, stmt.Body.Statements[0])
		select1 := stmt.Body.Statements[0].(*semantic.SelectStatement)
		assert.Equal(t, select1.From.TableRefs[0].Table, "dual")

		assert.IsType(t, &semantic.SelectStatement{}, stmt.Body.Statements[1])
		select2 := stmt.Body.Statements[1].(*semantic.SelectStatement)
		assert.Equal(t, select2.From.TableRefs[0].Table, "t")
	}
}
