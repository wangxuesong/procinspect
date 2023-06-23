package parser

import (
	"testing"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	testCase struct {
		name string
		text string
		root rootFunc
		Func testCaseFunc
	}

	testCaseFunc func(*testing.T, any)

	testSuite []testCase

	rootFunc func(*testing.T, string) any
)

func getRoot(t *testing.T, text string) any {
	p := plsql.NewParser(text)
	root := p.Sql_script()
	assert.Nil(t, p.Error(), p.Error())
	assert.NotNil(t, root)
	assert.IsType(t, &plsql.Sql_scriptContext{}, root)
	_, ok := root.(*plsql.Sql_scriptContext)
	assert.True(t, ok)

	node, err := GeneralScript(root)
	require.Nil(t, err, err)
	assert.NotNil(t, node)
	assert.Greater(t, len(node.Statements), 0)
	return node
}

func runTest(t *testing.T, input string, testFunc func(t *testing.T, node any), rootFunc ...rootFunc) {
	if len(rootFunc) > 0 {
		node := rootFunc[0](t, input)
		testFunc(t, node)
		return
	}
	node := getRoot(t, input)
	testFunc(t, node)
}

func runTestSuite(t *testing.T, tests testSuite) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.root != nil {
				runTest(t, test.text, test.Func, test.root)
			} else {
				runTest(t, test.text, test.Func)
			}
		})
	}
}

func TestParseCreatePackage(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "create package",
		text: `create or replace package test is
	procedure swth(a number);
end;
create or replace package body test is
	procedure swth(a number) is
	begin
		a := 1;
	end swth;
end;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 2)
			{
				stmt, ok := node.Statements[0].(*semantic.CreatePackageStatement)
				assert.True(t, ok)
				assert.NotNil(t, stmt)
				assert.Equal(t, "test", stmt.Name)
				assert.Equal(t, len(stmt.Procedures), 1)
				assert.Equal(t, "swth", stmt.Procedures[0].Name)
				assert.Equal(t, len(stmt.Procedures[0].Parameters), 1)
				assert.Equal(t, "a", stmt.Procedures[0].Parameters[0].Name)
				assert.Equal(t, "number", stmt.Procedures[0].Parameters[0].DataType)
			}
			{
				stmt, ok := node.Statements[1].(*semantic.CreatePackageBodyStatement)
				assert.True(t, ok)
				assert.NotNil(t, stmt)
				assert.Equal(t, "test", stmt.Name)
				assert.Equal(t, len(stmt.Procedures), 1)
				assert.Equal(t, "swth", stmt.Procedures[0].Name)
				assert.Equal(t, len(stmt.Procedures[0].Parameters), 1)
				assert.Equal(t, "a", stmt.Procedures[0].Parameters[0].Name)
				assert.Equal(t, "number", stmt.Procedures[0].Parameters[0].DataType)
				assert.NotNil(t, stmt.Procedures[0].Body)
			}
		},
	})

	tests = append(tests, testCase{
		name: "create package with type",
		text: `
create or replace package zznode.pkg_task_info is
  type type_date_tab is table of date index by binary_integer;
end;
`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			{
				stmt, ok := node.Statements[0].(*semantic.CreatePackageStatement)
				assert.True(t, ok)
				assert.NotNil(t, stmt)
				assert.Equal(t, "pkg_task_info", stmt.Name)
				assert.Equal(t, len(stmt.Types), 1)
				assert.IsType(t, &semantic.NestTableTypeDeclaration{}, stmt.Types[0])
				typeStmt := stmt.Types[0].(*semantic.NestTableTypeDeclaration)
				assert.Equal(t, "type_date_tab", typeStmt.Name)
				assert.Equal(t, 3, typeStmt.Line())
				assert.Equal(t, 3, typeStmt.Column())
			}
		},
	})

	tests = append(tests, testCase{
		name: "create package with function declaration",
		text: `
create or replace package zznode.pkg_task_info is
  function fun_rid_holiday(in_begin in date, in_end in date) return number;
end;
`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			{
				stmt, ok := node.Statements[0].(*semantic.CreatePackageStatement)
				assert.True(t, ok)
				assert.NotNil(t, stmt)
				assert.Equal(t, "pkg_task_info", stmt.Name)
				assert.Equal(t, len(stmt.Types), 1)
				assert.IsType(t, &semantic.FunctionDeclaration{}, stmt.Types[0])
				typeStmt := stmt.Types[0].(*semantic.FunctionDeclaration)
				assert.Equal(t, "fun_rid_holiday", typeStmt.Name)
				assert.Equal(t, 3, typeStmt.Line())
				assert.Equal(t, 3, typeStmt.Column())
			}
		},
	})

	runTestSuite(t, tests)
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
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "simple",
		text: `select * from dual, test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
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

		},
	})

	runTestSuite(t, tests)
}

func TestCreateProcedure(t *testing.T) {
	var tests = testSuite{}

	tests = append(tests, testCase{
		name: "create procedure",
		text: `create or replace procedure test is
	begin
		select 1 from dual;
		select 2 from t;
	end;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
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
		},
	})

	tests = append(tests, testCase{
		name: "create procedure with parameters",
		text: `CREATE OR REPLACE PROCEDURE PROC(PARAM NUMBER)
IS
LOCAL_PARAM NUMBER;
USER_EXCEPTION EXCEPTION;
  type type_date_tab is table of date index by binary_integer;
BEGIN
LOCAL_PARAM:=1;
END;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
			assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.CreateProcedureStatement)

			// assert line & column
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())

			assert.True(t, stmt.IsReplace)

			assert.Equal(t, stmt.Name, "PROC")

			assert.Equal(t, len(stmt.Parameters), 1)
			assert.Equal(t, len(stmt.Body.Statements), 1)
			assert.Equal(t, stmt.Parameters[0].Name, "PARAM")
			assert.Equal(t, stmt.Parameters[0].DataType, "NUMBER")

			assert.NotNil(t, stmt.Declarations)
			assert.Equal(t, len(stmt.Declarations), 3)

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
			// assert nest table type declaration
			assert.IsType(t, &semantic.NestTableTypeDeclaration{}, stmt.Declarations[2])
			decl3 := stmt.Declarations[2].(*semantic.NestTableTypeDeclaration)
			assert.Equal(t, decl3.Name, "type_date_tab")

			assert.NotNil(t, stmt.Body)
			assert.Equal(t, len(stmt.Body.Statements), 1)

			assert.IsType(t, &semantic.AssignmentStatement{}, stmt.Body.Statements[0])
			stmt1 := stmt.Body.Statements[0].(*semantic.AssignmentStatement)
			// assert line & column
			assert.Equal(t, 7, stmt1.Line())
			assert.Equal(t, 1, stmt1.Column())
			assert.Equal(t, stmt1.Left, "LOCAL_PARAM")
			assert.IsType(t, &semantic.NumericLiteral{}, stmt1.Right)
			right := stmt1.Right.(*semantic.NumericLiteral)
			assert.Equal(t, right.Value, int64(1))
		},
	})

	runTestSuite(t, tests)
}

func TestCreateFunction(t *testing.T) {
	var tests = testSuite{}

	tests = append(tests, testCase{
		name: "create function",
		text: `create or replace function test return number is
	begin
		select 1 from dual;
		return 1;
	end;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
			assert.IsType(t, &semantic.CreateFunctionStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.CreateFunctionStatement)

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

			assert.IsType(t, &semantic.ReturnStatement{}, stmt.Body.Statements[1])
			r := stmt.Body.Statements[1].(*semantic.ReturnStatement)
			// assert line & column
			assert.Equal(t, 4, r.Line())
			assert.Equal(t, 3, r.Column())
		},
	})

	tests = append(tests, testCase{
		name: "create function with parameters",
		text: `CREATE OR REPLACE FUNCTION PROC(PARAM NUMBER) RETURN NUMBER
IS
LOCAL_PARAM NUMBER;
USER_EXCEPTION EXCEPTION;
  type type_date_tab is table of date index by binary_integer;
BEGIN
LOCAL_PARAM:=1;
return 1;
END;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
			assert.IsType(t, &semantic.CreateFunctionStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.CreateFunctionStatement)

			// assert line & column
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())

			assert.True(t, stmt.IsReplace)

			assert.Equal(t, stmt.Name, "PROC")

			assert.Equal(t, len(stmt.Parameters), 1)
			assert.Equal(t, len(stmt.Body.Statements), 2)
			assert.Equal(t, stmt.Parameters[0].Name, "PARAM")
			assert.Equal(t, stmt.Parameters[0].DataType, "NUMBER")

			assert.NotNil(t, stmt.Declarations)
			assert.Equal(t, len(stmt.Declarations), 3)

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
			// assert nest table type declaration
			assert.IsType(t, &semantic.NestTableTypeDeclaration{}, stmt.Declarations[2])
			decl3 := stmt.Declarations[2].(*semantic.NestTableTypeDeclaration)
			assert.Equal(t, decl3.Name, "type_date_tab")

			assert.NotNil(t, stmt.Body)
			assert.Equal(t, len(stmt.Body.Statements), 2)

			assert.IsType(t, &semantic.AssignmentStatement{}, stmt.Body.Statements[0])
			stmt1 := stmt.Body.Statements[0].(*semantic.AssignmentStatement)
			// assert line & column
			assert.Equal(t, 7, stmt1.Line())
			assert.Equal(t, 1, stmt1.Column())
			assert.Equal(t, stmt1.Left, "LOCAL_PARAM")
			assert.IsType(t, &semantic.NumericLiteral{}, stmt1.Right)
			right := stmt1.Right.(*semantic.NumericLiteral)
			assert.Equal(t, right.Value, int64(1))

			assert.IsType(t, &semantic.ReturnStatement{}, stmt.Body.Statements[1])
			r := stmt.Body.Statements[1].(*semantic.ReturnStatement)
			// assert line & column
			assert.Equal(t, 8, r.Line())
			assert.Equal(t, 1, r.Column())
		},
	})

	runTestSuite(t, tests)
}

func TestCursorDeclaration(t *testing.T) {
	var tests testSuite

	tests = append(tests, testCase{
		name: "无参数游标声明",
		text: `create or replace procedure test is
	Cursor c_AllAws Is
		Select *
		From Asc_Work_Status
		Where Aws_Interfacer='ABB'
		Order By Aws_Asc_Id;
	Rec_AllAws c_AllAws%Rowtype;
	begin
		select 1 from dual;
	end;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
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
		},
	})

	tests = append(tests, testCase{
		name: "有参数游标声明",
		text: `create or replace procedure test is
	Cursor c_Aws(m_Areano Varchar2) Is
		Select *
		From Asc_Work_Status
		Where Aws_Interfacer='ABB'
		  And Aws_Curarea=m_Areano;
		Rec_Aws c_Aws%Rowtype;
	begin
		select 1 from dual;
	end;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
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
		},
	})

	runTestSuite(t, tests)
}

func TestCursor(t *testing.T) {
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "游标声明及使用",
		text: `CREATE OR REPLACE Procedure Abb_GenAnchorJob_P(i_Areanos Varchar2) Is
  v_Asc_Ids Varchar2(4000);

  Cursor c_AllAws Is
    Select *
    From Asc_Work_Status
    Where Aws_Interfacer='ABB'
    Order By Aws_Asc_Id;
  Rec_AllAws c_AllAws%Rowtype;

  Cursor c_Aws(m_Areano Varchar2) Is
    Select *
    From Asc_Work_Status
    Where Aws_Interfacer='ABB'
      And Aws_Curarea=m_Areano;
  Rec_Aws c_Aws%Rowtype;
  v_Areano Varchar2(100);
  v_Areanos Varchar2(4000);
  v_Index Integer;
Begin
  If i_Areanos Is Null Then
    If c_AllAws%Isopen Then
      Close c_AllAws;
    End If;
    Open c_AllAws;
    Loop
      Fetch c_AllAws Into Rec_AllAws;
      Exit When c_AllAws%NotFound;
      If v_Asc_Ids Is Null Then
        v_Asc_Ids := Rec_AllAws.Aws_Asc_Id;
      Else
        v_Asc_Ids := v_Asc_Ids || ','||Rec_AllAws.Aws_Asc_Id;
      End If;
    End Loop;
  Else
    v_Areanos := i_Areanos;
    Loop
      v_Index := Instr(v_Areanos,',');
      If v_Index>0 Then
        v_Areano := Substr(v_Areanos,1,v_Index-1);
        v_Areanos := Substr(v_Areanos,v_Index+1);
        v_Areanos := Nvl(v_Areanos,v_Index+1);
      Else
        v_Areano := v_Areanos;
        v_Areanos := '';
      End If;
      Exit When v_Areano Is Null;

      If c_Aws%Isopen Then
        Close c_Aws;
      End If;
      Open c_Aws(v_Areano);
      Loop
        Fetch c_Aws Into Rec_Aws;
        Exit When c_Aws%NotFound;
        If v_Asc_Ids Is Null Then
          v_Asc_Ids := Rec_Aws.Aws_Asc_Id;
        Else
          v_Asc_Ids := v_Asc_Ids||','||Rec_Aws.Aws_Asc_Id;
        End If;
      End Loop;
      Close c_Aws;
    End Loop;
  End If;

  If v_Asc_Ids Is Not Null Then
    Mon_Abb_Pak.gen_Anchor_Job_P(v_Asc_Ids);
  End If;
End;
/`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			// assert that the statement is a CreateProcedureStatement
			assert.IsType(t, &semantic.CreateProcedureStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.CreateProcedureStatement)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.True(t, stmt.IsReplace)
			assert.Equal(t, stmt.Name, "Abb_GenAnchorJob_P")

			// assert declaration
			{
				assert.NotNil(t, stmt.Declarations)
				assert.Equal(t, len(stmt.Declarations), 8)
				// assert the declaration is a variableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[0])
				varDecl := stmt.Declarations[0].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "v_Asc_Ids")
				assert.Equal(t, varDecl.DataType, "Varchar2(4000)")
				// assert the declaration is a CursorDeclaration
				assert.IsType(t, &semantic.CursorDeclaration{}, stmt.Declarations[1])
				cursorDecl := stmt.Declarations[1].(*semantic.CursorDeclaration)
				assert.Equal(t, cursorDecl.Name, "c_AllAws")
				assert.IsType(t, &semantic.SelectStatement{}, cursorDecl.Stmt)
				selectStmt := cursorDecl.Stmt.(*semantic.SelectStatement)
				assert.Equal(t, len(selectStmt.Fields.Fields), 1)
				assert.NotNil(t, selectStmt.From)
				assert.Equal(t, len(selectStmt.From.TableRefs), 1)
				assert.Equal(t, selectStmt.From.TableRefs[0].Table, "Asc_Work_Status")
				assert.IsType(t, &semantic.RelationalExpression{}, selectStmt.Where)
				expr := selectStmt.Where.(*semantic.RelationalExpression)
				assert.Equal(t, expr.Operator, "=")
				assert.IsType(t, &semantic.NameExpression{}, expr.Left)
				nameExpr := expr.Left.(*semantic.NameExpression)
				assert.Equal(t, nameExpr.Name, "Aws_Interfacer")
				assert.IsType(t, &semantic.StringLiteral{}, expr.Right)
				strLit := expr.Right.(*semantic.StringLiteral)
				assert.Equal(t, strLit.Value, "'ABB'")
				// assert the declaration is a VariableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[2])
				varDecl = stmt.Declarations[2].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "Rec_AllAws")
				assert.Equal(t, varDecl.DataType, "c_AllAws%Rowtype")
				// assert the declaration is a CursorDeclaration
				assert.IsType(t, &semantic.CursorDeclaration{}, stmt.Declarations[3])
				cursorDecl = stmt.Declarations[3].(*semantic.CursorDeclaration)
				assert.Equal(t, cursorDecl.Name, "c_Aws")
				// assert the declaration is a VariableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[4])
				varDecl = stmt.Declarations[4].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "Rec_Aws")
				assert.Equal(t, varDecl.DataType, "c_Aws%Rowtype")
				// assert the declaration is a VariableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[5])
				varDecl = stmt.Declarations[5].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "v_Areano")
				assert.Equal(t, varDecl.DataType, "Varchar2(100)")
				// assert the declaration is a VariableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[6])
				varDecl = stmt.Declarations[6].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "v_Areanos")
				assert.Equal(t, varDecl.DataType, "Varchar2(4000)")
				// assert the declaration is a VariableDeclaration
				assert.IsType(t, &semantic.VariableDeclaration{}, stmt.Declarations[7])
				varDecl = stmt.Declarations[7].(*semantic.VariableDeclaration)
				assert.Equal(t, varDecl.Name, "v_Index")
				assert.Equal(t, varDecl.DataType, "Integer")
			}

			// assert body
			{
				assert.NotNil(t, stmt.Body)
				assert.Equal(t, len(stmt.Body.Statements), 2)
				assert.IsType(t, &semantic.IfStatement{}, stmt.Body.Statements[0])
				ifStmt := stmt.Body.Statements[0].(*semantic.IfStatement)
				// assert the condition of the if statement
				assert.NotNil(t, ifStmt.Condition)
				assert.IsType(t, &semantic.UnaryLogicalExpression{}, ifStmt.Condition)
				unaryExp := ifStmt.Condition.(*semantic.UnaryLogicalExpression)
				assert.Equal(t, unaryExp.Operator, "NULL")
				assert.IsType(t, &semantic.NameExpression{}, unaryExp.Expr)
				nameExp := unaryExp.Expr.(*semantic.NameExpression)
				assert.Equal(t, nameExp.Name, "i_Areanos")
				//assert.Equal(t, ifStmt.Condition, "i_AreanosIsNull")
				// assert the then_block of the if statement
				assert.NotNil(t, ifStmt.ThenBlock)
				assert.Equal(t, len(ifStmt.ThenBlock), 3)
				assert.IsType(t, &semantic.IfStatement{}, ifStmt.ThenBlock[0])
				// assert the first if statement
				{
					ifStmt := ifStmt.ThenBlock[0].(*semantic.IfStatement)
					assert.NotNil(t, ifStmt.Condition)
					assert.IsType(t, &semantic.CursorAttribute{}, ifStmt.Condition)
					cursorAttr := ifStmt.Condition.(*semantic.CursorAttribute)
					assert.Equal(t, cursorAttr.Cursor, "c_AllAws")
					assert.Equal(t, cursorAttr.Attr, "ISOPEN")
					//assert.Equal(t, ifStmt.Condition, "c_AllAws%Isopen")
					assert.NotNil(t, ifStmt.ThenBlock)
					assert.IsType(t, &semantic.CloseStatement{}, ifStmt.ThenBlock[0])
					assert.Nil(t, ifStmt.ElseBlock)
				}
				assert.IsType(t, &semantic.OpenStatement{}, ifStmt.ThenBlock[1])
				assert.IsType(t, &semantic.LoopStatement{}, ifStmt.ThenBlock[2])
				// assert the loop statement
				{
					loopStmt := ifStmt.ThenBlock[2].(*semantic.LoopStatement)
					assert.NotNil(t, loopStmt.Statements)
					assert.Equal(t, len(loopStmt.Statements), 3)
					assert.IsType(t, &semantic.FetchStatement{}, loopStmt.Statements[0])
					// assert the fetch statement
					{
						fetchStmt := loopStmt.Statements[0].(*semantic.FetchStatement)
						assert.NotNil(t, fetchStmt.Cursor)
						assert.Equal(t, fetchStmt.Cursor, "c_AllAws")
						assert.NotNil(t, fetchStmt.Into)
						assert.Equal(t, fetchStmt.Into, "Rec_AllAws")
					}
					assert.IsType(t, &semantic.ExitStatement{}, loopStmt.Statements[1])
					// assert the exit statement
					{
						exitStmt := loopStmt.Statements[1].(*semantic.ExitStatement)
						assert.NotNil(t, exitStmt.Condition)
						assert.IsType(t, &semantic.CursorAttribute{}, exitStmt.Condition)
						attr := exitStmt.Condition.(*semantic.CursorAttribute)
						assert.Equal(t, attr.Cursor, "c_AllAws")
						assert.Equal(t, attr.Attr, "NOTFOUND")
					}
					assert.IsType(t, &semantic.IfStatement{}, loopStmt.Statements[2])
					// assert the if statement
					{
						ifStmt := loopStmt.Statements[2].(*semantic.IfStatement)
						assert.NotNil(t, ifStmt.Condition)
						assert.IsType(t, &semantic.UnaryLogicalExpression{}, ifStmt.Condition)
						unaryExp := ifStmt.Condition.(*semantic.UnaryLogicalExpression)
						assert.Equal(t, unaryExp.Operator, "NULL")
						assert.IsType(t, &semantic.NameExpression{}, unaryExp.Expr)
						nameExp := unaryExp.Expr.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Asc_Ids")
						assert.NotNil(t, ifStmt.ThenBlock)
						assert.NotNil(t, ifStmt.ElseBlock)
						assert.Equal(t, len(ifStmt.ThenBlock), 1)
						assert.Equal(t, len(ifStmt.ElseBlock), 1)
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[0])
						stmt := ifStmt.ThenBlock[0].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Asc_Ids")
						assert.IsType(t, &semantic.DotExpression{}, stmt.Right)
						dotExp := stmt.Right.(*semantic.DotExpression)
						assert.Equal(t, dotExp.Name, "Aws_Asc_Id")
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Parent)
						nameExp = dotExp.Parent.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Rec_AllAws")
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
						stmt = ifStmt.ElseBlock[0].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Asc_Ids")
						assert.IsType(t, &semantic.BinaryExpression{}, stmt.Right)
						binaryExp := stmt.Right.(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "||")
						assert.IsType(t, &semantic.BinaryExpression{}, binaryExp.Left)
						assert.IsType(t, &semantic.DotExpression{}, binaryExp.Right)
						dotExp = binaryExp.Right.(*semantic.DotExpression)
						assert.Equal(t, dotExp.Name, "Aws_Asc_Id")
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Parent)
						nameExp = dotExp.Parent.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Rec_AllAws")
						binaryExp = binaryExp.Left.(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "||")
						assert.IsType(t, &semantic.NameExpression{}, binaryExp.Left)
						nameExp = binaryExp.Left.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Asc_Ids")
						assert.IsType(t, &semantic.StringLiteral{}, binaryExp.Right)
						stringLit := binaryExp.Right.(*semantic.StringLiteral)
						assert.Equal(t, stringLit.Value, "','")
					}
				}
				// assert the else_block of the if statement
				assert.NotNil(t, ifStmt.ElseBlock)
				assert.Equal(t, len(ifStmt.ElseBlock), 2)
				assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
				assert.IsType(t, &semantic.LoopStatement{}, ifStmt.ElseBlock[1])
				// assert the loop statement
				{
					loopStmt := ifStmt.ElseBlock[1].(*semantic.LoopStatement)
					assert.NotNil(t, loopStmt.Statements)
					assert.Equal(t, len(loopStmt.Statements), 7)
					assert.IsType(t, &semantic.AssignmentStatement{}, loopStmt.Statements[0])
					stmt := loopStmt.Statements[0].(*semantic.AssignmentStatement)
					assert.Equal(t, stmt.Left, "v_Index")
					assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Right)
					funcCallExp := stmt.Right.(*semantic.FunctionCallExpression)
					assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Name)
					expr := funcCallExp.Name.(*semantic.NameExpression)
					assert.Equal(t, expr.Name, "Instr")
					assert.Equal(t, len(funcCallExp.Args), 2)
					assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Args[0])
					nameExp := funcCallExp.Args[0].(*semantic.NameExpression)
					assert.Equal(t, nameExp.Name, "v_Areanos")
					assert.IsType(t, &semantic.StringLiteral{}, funcCallExp.Args[1])
					stringLit := funcCallExp.Args[1].(*semantic.StringLiteral)
					assert.Equal(t, stringLit.Value, "','")
					assert.IsType(t, &semantic.IfStatement{}, loopStmt.Statements[1])
					// assert the else_block of the if statement
					{
						ifStmt := loopStmt.Statements[1].(*semantic.IfStatement)
						assert.NotNil(t, ifStmt.Condition)
						assert.IsType(t, &semantic.RelationalExpression{}, ifStmt.Condition)
						relExp := ifStmt.Condition.(*semantic.RelationalExpression)
						assert.Equal(t, relExp.Operator, ">")
						assert.IsType(t, &semantic.NameExpression{}, relExp.Left)
						nameExp := relExp.Left.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Index")
						assert.IsType(t, &semantic.NumericLiteral{}, relExp.Right)
						numericLit := relExp.Right.(*semantic.NumericLiteral)
						assert.Equal(t, numericLit.Value, int64(0))
						assert.NotNil(t, ifStmt.ThenBlock)
						assert.Equal(t, len(ifStmt.ThenBlock), 3)
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[0])
						stmt := ifStmt.ThenBlock[0].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Areano")
						assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Right)
						funcCallExp := stmt.Right.(*semantic.FunctionCallExpression)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Name)
						expr := funcCallExp.Name.(*semantic.NameExpression)
						assert.Equal(t, expr.Name, "SUBSTR")
						assert.Equal(t, len(funcCallExp.Args), 3)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Args[0])
						nameExp = funcCallExp.Args[0].(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Areanos")
						assert.IsType(t, &semantic.NumericLiteral{}, funcCallExp.Args[1])
						numericLit = funcCallExp.Args[1].(*semantic.NumericLiteral)
						assert.Equal(t, numericLit.Value, int64(1))
						assert.IsType(t, &semantic.BinaryExpression{}, funcCallExp.Args[2])
						binaryExp := funcCallExp.Args[2].(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "-")
						assert.IsType(t, &semantic.NameExpression{}, binaryExp.Left)
						nameExp = binaryExp.Left.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Index")
						assert.IsType(t, &semantic.NumericLiteral{}, binaryExp.Right)
						numericLit = binaryExp.Right.(*semantic.NumericLiteral)
						assert.Equal(t, numericLit.Value, int64(1))
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[1])
						stmt = ifStmt.ThenBlock[1].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Areanos")
						assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Right)
						funcCallExp = stmt.Right.(*semantic.FunctionCallExpression)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Name)
						nameExp = funcCallExp.Name.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "SUBSTR")
						assert.Equal(t, len(funcCallExp.Args), 2)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Args[0])
						nameExp = funcCallExp.Args[0].(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Areanos")
						assert.IsType(t, &semantic.BinaryExpression{}, funcCallExp.Args[1])
						binaryExp = funcCallExp.Args[1].(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "+")
						assert.IsType(t, &semantic.NameExpression{}, binaryExp.Left)
						nameExp = binaryExp.Left.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Index")
						assert.IsType(t, &semantic.NumericLiteral{}, binaryExp.Right)
						numericLit = binaryExp.Right.(*semantic.NumericLiteral)
						assert.Equal(t, numericLit.Value, int64(1))
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[2])
						stmt = ifStmt.ThenBlock[2].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Areanos")
						assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Right)
						funcCallExp = stmt.Right.(*semantic.FunctionCallExpression)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Name)
						nameExp = funcCallExp.Name.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "NVL")
						assert.Equal(t, len(funcCallExp.Args), 2)
						assert.IsType(t, &semantic.NameExpression{}, funcCallExp.Args[0])
						nameExp = funcCallExp.Args[0].(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Areanos")
						assert.IsType(t, &semantic.BinaryExpression{}, funcCallExp.Args[1])
						binaryExp = funcCallExp.Args[1].(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "+")
						assert.IsType(t, &semantic.NameExpression{}, binaryExp.Left)
						nameExp = binaryExp.Left.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Index")
						assert.IsType(t, &semantic.NumericLiteral{}, binaryExp.Right)
						numericLit = binaryExp.Right.(*semantic.NumericLiteral)
						assert.Equal(t, numericLit.Value, int64(1))
						assert.NotNil(t, ifStmt.ElseBlock)
						assert.Equal(t, len(ifStmt.ElseBlock), 2)
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[1])
					}
					assert.IsType(t, &semantic.ExitStatement{}, loopStmt.Statements[2])
					// assert the exit statement
					{
						exitStmt := loopStmt.Statements[2].(*semantic.ExitStatement)
						assert.NotNil(t, exitStmt.Condition)
						assert.IsType(t, &semantic.UnaryLogicalExpression{}, exitStmt.Condition)
						expr := exitStmt.Condition.(*semantic.UnaryLogicalExpression)
						assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
						name := expr.Expr.(*semantic.NameExpression)
						assert.Equal(t, name.Name, "v_Areano")
						assert.Equal(t, expr.Operator, "NULL")
					}
					assert.IsType(t, &semantic.IfStatement{}, loopStmt.Statements[3])
					{
						ifStmt := loopStmt.Statements[3].(*semantic.IfStatement)
						assert.NotNil(t, ifStmt.Condition)
						assert.IsType(t, &semantic.CursorAttribute{}, ifStmt.Condition)
						attr := ifStmt.Condition.(*semantic.CursorAttribute)
						assert.Equal(t, attr.Cursor, "c_Aws")
						assert.Equal(t, attr.Attr, "ISOPEN")
						assert.NotNil(t, ifStmt.ThenBlock)
						assert.Equal(t, len(ifStmt.ThenBlock), 1)
						assert.IsType(t, &semantic.CloseStatement{}, ifStmt.ThenBlock[0])
					}
					assert.IsType(t, &semantic.OpenStatement{}, loopStmt.Statements[4])
					assert.IsType(t, &semantic.LoopStatement{}, loopStmt.Statements[5])
					// assert the loop statement
					{
						loopStmt := loopStmt.Statements[5].(*semantic.LoopStatement)
						assert.NotNil(t, loopStmt.Statements)
						assert.Equal(t, len(loopStmt.Statements), 3)
						assert.IsType(t, &semantic.FetchStatement{}, loopStmt.Statements[0])
						// assert the fetch statement
						{
							fetchStmt := loopStmt.Statements[0].(*semantic.FetchStatement)
							assert.NotNil(t, fetchStmt.Cursor)
							assert.Equal(t, fetchStmt.Cursor, "c_Aws")
							assert.NotNil(t, fetchStmt.Into)
							assert.Equal(t, fetchStmt.Into, "Rec_Aws")
						}
						assert.IsType(t, &semantic.ExitStatement{}, loopStmt.Statements[1])
						// assert the exit statement
						{
							exitStmt := loopStmt.Statements[1].(*semantic.ExitStatement)
							assert.NotNil(t, exitStmt.Condition)
							assert.IsType(t, &semantic.CursorAttribute{}, exitStmt.Condition)
							attr := exitStmt.Condition.(*semantic.CursorAttribute)
							assert.Equal(t, attr.Cursor, "c_Aws")
							assert.Equal(t, attr.Attr, "NOTFOUND")
						}
						assert.IsType(t, &semantic.IfStatement{}, loopStmt.Statements[2])
						// assert the if statement
						{
							ifStmt := loopStmt.Statements[2].(*semantic.IfStatement)
							assert.NotNil(t, ifStmt.Condition)
							assert.IsType(t, &semantic.UnaryLogicalExpression{}, ifStmt.Condition)
							expr := ifStmt.Condition.(*semantic.UnaryLogicalExpression)
							assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
							name := expr.Expr.(*semantic.NameExpression)
							assert.Equal(t, name.Name, "v_Asc_Ids")
							assert.Equal(t, expr.Operator, "NULL")
							assert.NotNil(t, ifStmt.ThenBlock)
							assert.Equal(t, len(ifStmt.ThenBlock), 1)
							assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[0])
							stmt := ifStmt.ThenBlock[0].(*semantic.AssignmentStatement)
							assert.Equal(t, stmt.Left, "v_Asc_Ids")
							assert.IsType(t, &semantic.DotExpression{}, stmt.Right)
							dotExp := stmt.Right.(*semantic.DotExpression)
							assert.Equal(t, dotExp.Name, "Aws_Asc_Id")
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Parent)
							relExp := dotExp.Parent.(*semantic.NameExpression)
							assert.Equal(t, relExp.Name, "Rec_Aws")
							assert.NotNil(t, ifStmt.ElseBlock)
							assert.Equal(t, len(ifStmt.ElseBlock), 1)
							assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
							stmt = ifStmt.ElseBlock[0].(*semantic.AssignmentStatement)
							assert.Equal(t, stmt.Left, "v_Asc_Ids")
							assert.IsType(t, &semantic.BinaryExpression{}, stmt.Right)
							binaryExp := stmt.Right.(*semantic.BinaryExpression)
							assert.Equal(t, binaryExp.Operator, "||")
							assert.IsType(t, &semantic.BinaryExpression{}, binaryExp.Left)
							assert.IsType(t, &semantic.DotExpression{}, binaryExp.Right)
							dotExp = binaryExp.Right.(*semantic.DotExpression)
							assert.Equal(t, dotExp.Name, "Aws_Asc_Id")
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Parent)
							nameExp = dotExp.Parent.(*semantic.NameExpression)
							assert.Equal(t, nameExp.Name, "Rec_Aws")
							binaryExp = binaryExp.Left.(*semantic.BinaryExpression)
							assert.Equal(t, binaryExp.Operator, "||")
							assert.IsType(t, &semantic.NameExpression{}, binaryExp.Left)
							nameExp = binaryExp.Left.(*semantic.NameExpression)
							assert.Equal(t, nameExp.Name, "v_Asc_Ids")
							assert.IsType(t, &semantic.StringLiteral{}, binaryExp.Right)
							stringLit := binaryExp.Right.(*semantic.StringLiteral)
							assert.Equal(t, stringLit.Value, "','")
						}
					}
					assert.IsType(t, &semantic.CloseStatement{}, loopStmt.Statements[6])
				}
				assert.IsType(t, &semantic.IfStatement{}, stmt.Body.Statements[1])
				// assert the if statement
				{
					ifStmt := stmt.Body.Statements[1].(*semantic.IfStatement)
					assert.NotNil(t, ifStmt.Condition)
					assert.IsType(t, &semantic.UnaryLogicalExpression{}, ifStmt.Condition)
					expr := ifStmt.Condition.(*semantic.UnaryLogicalExpression)
					assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
					name := expr.Expr.(*semantic.NameExpression)
					assert.Equal(t, name.Name, "v_Asc_Ids")
					assert.Equal(t, expr.Operator, "NULL")
					assert.True(t, expr.Not)
					assert.NotNil(t, ifStmt.ThenBlock)
					assert.Equal(t, len(ifStmt.ThenBlock), 1)
					assert.IsType(t, &semantic.ProcedureCall{}, ifStmt.ThenBlock[0])
					// assert the procedure call
					{
						procCall := ifStmt.ThenBlock[0].(*semantic.ProcedureCall)
						assert.IsType(t, &semantic.DotExpression{}, procCall.Name)
						dotExp := procCall.Name.(*semantic.DotExpression)
						assert.Equal(t, dotExp.Name, "gen_Anchor_Job_P")
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Parent)
						nameExp := dotExp.Parent.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Mon_Abb_Pak")
						assert.Equal(t, len(procCall.Arguments), 1)
						assert.IsType(t, &semantic.NameExpression{}, procCall.Arguments[0])
						nameExp = procCall.Arguments[0].(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "v_Asc_Ids")
					}
					assert.Nil(t, ifStmt.ElseBlock)
				}
			}
		},
	})

	runTestSuite(t, tests)
}

func TestBlock(t *testing.T) {
	getBlock := func(t *testing.T, text string) any {
		p := plsql.NewParser(text)
		root := p.Block()
		assert.Nil(t, p.Error())

		visitor := newPlSqlVisitor()
		node := visitor.VisitBlock(root.(*plsql.BlockContext)).(*semantic.BlockStatement)

		return node
	}

	tests := testSuite{}
	tests = append(tests, testCase{
		name: "block",
		root: getBlock,
		text: `
DECLARE
	a NUMBER := 1;
BEGIN
	a:=1;
END`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.BlockStatement)
			assert.Equal(t, len(node.Declarations), 1)
			assert.IsType(t, &semantic.VariableDeclaration{}, node.Declarations[0])
			decl := node.Declarations[0].(*semantic.VariableDeclaration)
			assert.Equal(t, decl.Name, "a")
			assert.Equal(t, decl.DataType, "NUMBER")
			assert.NotNil(t, decl.Initialization)
			assert.IsType(t, &semantic.NumericLiteral{}, decl.Initialization)
			assert.NotNil(t, node.Body)
			assert.Equal(t, len(node.Body.Statements), 1)
			assert.IsType(t, &semantic.AssignmentStatement{}, node.Body.Statements[0])
			stmt := node.Body.Statements[0].(*semantic.AssignmentStatement)
			assert.Equal(t, stmt.Left, "a")
			assert.IsType(t, &semantic.NumericLiteral{}, stmt.Right)
			lit := stmt.Right.(*semantic.NumericLiteral)
			assert.Equal(t, lit.Value, int64(1))
		},
	})

	runTestSuite(t, tests)
}
