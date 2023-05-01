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
	assert.Equal(t, 1, stmt.Line())
	assert.Equal(t, 1, stmt.Column())
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

	node := getRoot(t, text)

	{ // assert that the statement is a CreateProcedureStatement
		assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
		stmt := node.Statements[0].(*semantic.CreateProcedureStatement)

		assert.Equal(t, 1, stmt.Line())
		assert.Equal(t, 1, stmt.Column())

		assert.True(t, stmt.IsReplace)

		assert.Equal(t, stmt.Name, "test")

		assert.NotNil(t, stmt.Body)
		assert.Equal(t, len(stmt.Body.Statements), 2)

		assert.IsType(t, &semantic.SelectStatement{}, stmt.Body.Statements[0])
		select1 := stmt.Body.Statements[0].(*semantic.SelectStatement)
		assert.Equal(t, select1.From.TableRefs[0].Table, "dual")
		// assert line & column
		assert.Equal(t, 3, select1.Line())
		assert.Equal(t, 3, select1.Column())

		assert.IsType(t, &semantic.SelectStatement{}, stmt.Body.Statements[1])
		select2 := stmt.Body.Statements[1].(*semantic.SelectStatement)
		assert.Equal(t, select2.From.TableRefs[0].Table, "t")
		// assert line & column
		assert.Equal(t, 4, select2.Line())
		assert.Equal(t, 3, select2.Column())
	}

	text = `CREATE OR REPLACE PROCEDURE PROC(PARAM NUMBER)
IS
LOCAL_PARAM NUMBER;
USER_EXCEPTION EXCEPTION;
BEGIN
LOCAL_PARAM:=0;
END;`

	node = getRoot(t, text)

	{ // assert that the statement is a CreateProcedureStatement
		assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
		stmt := node.Statements[0].(*semantic.CreateProcedureStatement)

		// assert line & column
		assert.Equal(t, 1, stmt.Line())
		assert.Equal(t, 1, stmt.Column())

		assert.True(t, stmt.IsReplace)

		assert.Equal(t, stmt.Name, "PROC")

		assert.Equal(t, len(stmt.Parameters), 1)
		assert.Equal(t, &semantic.Parameter{Name: "PARAM", DataType: "NUMBER"}, stmt.Parameters[0])

		assert.NotNil(t, stmt.Declarations)
		assert.Equal(t, len(stmt.Declarations), 2)

		assert.IsType(t, &semantic.VariableDeclaration{Name: "LOCAL_PARAM", DataType: "NUMBER"}, stmt.Declarations[0])
		// assert line & column
		assert.Equal(t, 3, stmt.Declarations[0].Line())
		assert.Equal(t, 1, stmt.Declarations[0].Column())
		decl1 := stmt.Declarations[0].(*semantic.VariableDeclaration)
		assert.Equal(t, decl1.Name, "LOCAL_PARAM")
		assert.Equal(t, decl1.DataType, "NUMBER")
		assert.IsType(t, &semantic.ExceptionDeclaration{Name: "USER_EXCEPTION"}, stmt.Declarations[1])
		// assert line & column
		assert.Equal(t, 4, stmt.Declarations[1].Line())
		assert.Equal(t, 1, stmt.Declarations[1].Column())
		decl2 := stmt.Declarations[1].(*semantic.ExceptionDeclaration)
		assert.Equal(t, decl2.Name, "USER_EXCEPTION")

		assert.NotNil(t, stmt.Body)
		assert.Equal(t, len(stmt.Body.Statements), 1)

		assert.IsType(t, &semantic.AssignmentStatement{}, stmt.Body.Statements[0])
		stmt1 := stmt.Body.Statements[0].(*semantic.AssignmentStatement)
		// assert line & column
		assert.Equal(t, 6, stmt1.Line())
		assert.Equal(t, 1, stmt1.Column())
		assert.Equal(t, stmt1.Left, "LOCAL_PARAM")
		assert.Equal(t, stmt1.Right, "0")
	}
}

func TestCursorDeclaration(t *testing.T) {
	text := `create or replace procedure test is
	Cursor c_AllAws Is
		Select *
		From Asc_Work_Status
		Where Aws_Interfacer='ABB'
		Order By Aws_Asc_Id;
	Rec_AllAws c_AllAws%Rowtype;
	begin
		select 1 from dual;
	end;`

	node := getRoot(t, text)

	{ // assert that the statement is a CreateProcedureStatement
		assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
		stmt := node.Statements[0].(*semantic.CreateProcedureStatement)
		assert.Equal(t, 1, stmt.Line())
		assert.Equal(t, 1, stmt.Column())
		assert.True(t, stmt.IsReplace)
		assert.Equal(t, stmt.Name, "test")

		// assert declaration
		assert.NotNil(t, stmt.Declarations)
		assert.Equal(t, len(stmt.Declarations), 2)
		// assert the declaration is a CursorDeclaration
		assert.IsType(t, &semantic.CursorDeclaration{}, stmt.Declarations[0])
		decl := stmt.Declarations[0].(*semantic.CursorDeclaration)
		assert.Equal(t, decl.Name, "c_AllAws")
		// assert the parameters is nil
		assert.Nil(t, decl.Parameters)
		// assert the declaration is a VariableDeclaration
		assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[1])
		vardecl := stmt.Declarations[1].(*semantic.VariableDeclaration)
		assert.Equal(t, vardecl.Name, "Rec_AllAws")
		assert.Equal(t, vardecl.DataType, "c_AllAws%Rowtype")
		// assert the statement is a SelectStatement
		assert.NotNil(t, decl.Stmt)
		assert.IsType(t, &semantic.SelectStatement{}, decl.Stmt)
	}

	text = `create or replace procedure test is
	Cursor c_Aws(m_Areano Varchar2) Is
		Select *
		From Asc_Work_Status
		Where Aws_Interfacer='ABB'
		  And Aws_Curarea=m_Areano;
		Rec_Aws c_Aws%Rowtype;
	begin
		select 1 from dual;
	end;`

	node = getRoot(t, text)

	{ // assert that the statement is a CreateProcedureStatement
		assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
		stmt := node.Statements[0].(*semantic.CreateProcedureStatement)
		assert.Equal(t, 1, stmt.Line())
		assert.Equal(t, 1, stmt.Column())
		assert.True(t, stmt.IsReplace)
		assert.Equal(t, stmt.Name, "test")

		// assert declaration
		assert.NotNil(t, stmt.Declarations)
		assert.Equal(t, len(stmt.Declarations), 2)
		// assert the declaration is a CursorDeclaration
		assert.IsType(t, &semantic.CursorDeclaration{}, stmt.Declarations[0])
		decl := stmt.Declarations[0].(*semantic.CursorDeclaration)
		assert.Equal(t, decl.Name, "c_Aws")
		// assert the parameters is not nil
		assert.NotNil(t, decl.Parameters)
		assert.Equal(t, len(decl.Parameters), 1)
		// assert the parameter is a Parameter
		param := decl.Parameters[0]
		assert.Equal(t, param.Name, "m_Areano")
		assert.Equal(t, param.DataType, "Varchar2")
		// assert the declaration is a VariableDeclaration
		assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[1])
		vardecl := stmt.Declarations[1].(*semantic.VariableDeclaration)
		assert.Equal(t, vardecl.Name, "Rec_Aws")
		assert.Equal(t, vardecl.DataType, "c_Aws%Rowtype")
		// assert the statement is a SelectStatement
		assert.NotNil(t, decl.Stmt)
		assert.IsType(t, &semantic.SelectStatement{}, decl.Stmt)
	}
}

func getRoot(t *testing.T, text string) *semantic.Script {
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
	return node
}
