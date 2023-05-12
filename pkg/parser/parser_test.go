package parser

import (
	"fmt"
	"testing"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"

	"github.com/stretchr/testify/assert"
)

type (
	testCase struct {
		name string
		text string
		Func testCaseFunc
	}

	testCaseFunc func(*testing.T, *semantic.Script)

	testSuite []testCase
)

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

func runTest(t *testing.T, input string, testFunc func(t *testing.T, node *semantic.Script)) {
	node := getRoot(t, input)
	testFunc(t, node)
}

func runTestSuite(t *testing.T, tests testSuite) {
	for _, test := range tests {
		fmt.Println(test.name)
		runTest(t, test.text, test.Func)
	}
}

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
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "simple",
		text: `select * from dual, test;`,
		Func: func(t *testing.T, node *semantic.Script) {
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
		Func: func(t *testing.T, node *semantic.Script) {
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
BEGIN
LOCAL_PARAM:=1;
END;`,
		Func: func(t *testing.T, node *semantic.Script) {
			// assert that the statement is a CreateProcedureStatement
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
			assert.IsType(t, &semantic.NumericLiteral{}, stmt1.Right)
			right := stmt1.Right.(*semantic.NumericLiteral)
			assert.Equal(t, right.Value, int64(1))
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
		Func: func(t *testing.T, node *semantic.Script) {
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
		Func: func(t *testing.T, node *semantic.Script) {
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
		Func: func(t *testing.T, node *semantic.Script) {
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
						assert.IsType(t, &semantic.NameExpression{}, stmt.Right)
						nameExp = stmt.Right.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Rec_AllAws.Aws_Asc_Id")
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
						stmt = ifStmt.ElseBlock[0].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Asc_Ids")
						assert.IsType(t, &semantic.BinaryExpression{}, stmt.Right)
						binaryExp := stmt.Right.(*semantic.BinaryExpression)
						assert.Equal(t, binaryExp.Operator, "||")
						assert.IsType(t, &semantic.BinaryExpression{}, binaryExp.Left)
						assert.IsType(t, &semantic.NameExpression{}, binaryExp.Right)
						nameExp = binaryExp.Right.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Rec_AllAws.Aws_Asc_Id")
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
					assert.Equal(t, funcCallExp.Name, "Instr")
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
						assert.Equal(t, len(ifStmt.ThenBlock), 2)
						assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ThenBlock[0])
						stmt := ifStmt.ThenBlock[0].(*semantic.AssignmentStatement)
						assert.Equal(t, stmt.Left, "v_Areano")
						assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Right)
						funcCallExp := stmt.Right.(*semantic.FunctionCallExpression)
						assert.Equal(t, funcCallExp.Name, "SUBSTR")
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
						assert.Equal(t, funcCallExp.Name, "SUBSTR")
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
							assert.IsType(t, &semantic.NameExpression{}, stmt.Right)
							relExp := stmt.Right.(*semantic.NameExpression)
							assert.Equal(t, relExp.Name, "Rec_Aws.Aws_Asc_Id")
							assert.NotNil(t, ifStmt.ElseBlock)
							assert.Equal(t, len(ifStmt.ElseBlock), 1)
							assert.IsType(t, &semantic.AssignmentStatement{}, ifStmt.ElseBlock[0])
							stmt = ifStmt.ElseBlock[0].(*semantic.AssignmentStatement)
							assert.Equal(t, stmt.Left, "v_Asc_Ids")
							assert.IsType(t, &semantic.BinaryExpression{}, stmt.Right)
							binaryExp := stmt.Right.(*semantic.BinaryExpression)
							assert.Equal(t, binaryExp.Operator, "||")
							assert.IsType(t, &semantic.BinaryExpression{}, binaryExp.Left)
							assert.IsType(t, &semantic.NameExpression{}, binaryExp.Right)
							nameExp = binaryExp.Right.(*semantic.NameExpression)
							assert.Equal(t, nameExp.Name, "Rec_Aws.Aws_Asc_Id")
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
						assert.Equal(t, procCall.Name, "Mon_Abb_Pak.gen_Anchor_Job_P")
						assert.Equal(t, len(procCall.Arguments), 1)
						assert.IsType(t, &semantic.Argument{}, procCall.Arguments[0])
						assert.Equal(t, procCall.Arguments[0].Name, "v_Asc_Ids")
					}
					assert.Nil(t, ifStmt.ElseBlock)
				}
			}
		},
	})

	runTestSuite(t, tests)
}
