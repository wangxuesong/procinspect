package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
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
			assert.Equal(t, semantic.Span{Start: 0, End: 23}, stmt.Span())
			assert.Equal(t, len(stmt.Fields.Fields), 1)
			assert.Equal(t, stmt.Fields.Fields[0].WildCard.Table, "*")
			assert.Equal(t, len(stmt.From.TableRefs), 2)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.Equal(t, stmt.From.TableRefs[1].Table, "test")

		},
	})

	tests = append(tests, testCase{
		name: "simple projection",
		text: `select t.* from dual, test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, len(stmt.Fields.Fields), 1)
			assert.Equal(t, stmt.Fields.Fields[0].WildCard.Table, "t")
			assert.Equal(t, len(stmt.From.TableRefs), 2)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.Equal(t, stmt.From.TableRefs[1].Table, "test")

		},
	})

	tests = append(tests, testCase{
		name: "in expression",
		text: `select * from dual where a in (1, 2, 3);`,
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
			assert.Equal(t, len(stmt.From.TableRefs), 1)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.NotNil(t, stmt.Where)
			assert.IsType(t, &semantic.InExpression{}, stmt.Where)
			expr := stmt.Where.(*semantic.InExpression)
			assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
			name := expr.Expr.(*semantic.NameExpression)
			assert.Equal(t, "a", name.Name)
			assert.Equal(t, 3, len(expr.Elems))
			assert.IsType(t, &semantic.NumericLiteral{}, expr.Elems[0])
			elem := expr.Elems[0].(*semantic.NumericLiteral)
			assert.Equal(t, int64(1), elem.Value)
			assert.IsType(t, &semantic.NumericLiteral{}, expr.Elems[1])
			elem = expr.Elems[1].(*semantic.NumericLiteral)
			assert.Equal(t, int64(2), elem.Value)
			assert.IsType(t, &semantic.NumericLiteral{}, expr.Elems[2])
			elem = expr.Elems[2].(*semantic.NumericLiteral)
			assert.Equal(t, int64(3), elem.Value)
		},
	})

	tests = append(tests, testCase{
		name: "between expression",
		text: `select * from dual where a between 1 and 2;`,
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
			assert.Equal(t, len(stmt.From.TableRefs), 1)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.NotNil(t, stmt.Where)
			assert.IsType(t, &semantic.BetweenExpression{}, stmt.Where)
			expr := stmt.Where.(*semantic.BetweenExpression)
			assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
			name := expr.Expr.(*semantic.NameExpression)
			assert.Equal(t, "a", name.Name)
			assert.Equal(t, 2, len(expr.Elems))
			assert.IsType(t, &semantic.NumericLiteral{}, expr.Elems[0])
			elem := expr.Elems[0].(*semantic.NumericLiteral)
			assert.Equal(t, int64(1), elem.Value)
			assert.IsType(t, &semantic.NumericLiteral{}, expr.Elems[1])
			elem = expr.Elems[1].(*semantic.NumericLiteral)
			assert.Equal(t, int64(2), elem.Value)
		},
	})

	tests = append(tests, testCase{
		name: "like expression",
		text: `select * from dual where a like '%1%';`,
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
			assert.Equal(t, len(stmt.From.TableRefs), 1)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.NotNil(t, stmt.Where)
			assert.IsType(t, &semantic.LikeExpression{}, stmt.Where)
			expr := stmt.Where.(*semantic.LikeExpression)
			assert.IsType(t, &semantic.NameExpression{}, expr.Expr)
			name := expr.Expr.(*semantic.NameExpression)
			assert.Equal(t, "a", name.Name)
		},
	})

	tests = append(tests, testCase{
		name: "exists expression",
		text: `select * from dual where exists (select * from dual);`,
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
			assert.Equal(t, len(stmt.From.TableRefs), 1)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
			assert.NotNil(t, stmt.Where)
			assert.IsType(t, &semantic.ExistsExpression{}, stmt.Where)
			expr := stmt.Where.(*semantic.ExistsExpression)
			assert.IsType(t, &semantic.QueryExpression{}, expr.Expr)
			query := expr.Expr.(*semantic.QueryExpression).Query
			assert.Equal(t, query.Fields.Fields[0].WildCard.Table, "*")
			assert.Equal(t, stmt.From.TableRefs[0].Table, "dual")
		},
	})

	tests = append(tests, testCase{
		name: "out join expression",
		text: `select * from t1,t2 where t1.id=t2.id(+);`,
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
			assert.Equal(t, stmt.From.TableRefs[0].Table, "t1")
			assert.NotNil(t, stmt.Where)
			assert.IsType(t, &semantic.RelationalExpression{}, stmt.Where)
			expr := stmt.Where.(*semantic.RelationalExpression)
			assert.IsType(t, &semantic.DotExpression{}, expr.Left)
			dot := expr.Left.(*semantic.DotExpression)
			assert.IsType(t, &semantic.NameExpression{}, dot.Name)
			name := dot.Name.(*semantic.NameExpression)
			assert.Equal(t, "id", name.Name)
			assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
			dot = dot.Parent.(*semantic.DotExpression)
			assert.IsType(t, &semantic.NameExpression{}, dot.Name)
			name = dot.Name.(*semantic.NameExpression)
			assert.Equal(t, "t1", name.Name)
			assert.IsType(t, &semantic.OuterJoinExpression{}, expr.Right)
			join := expr.Right.(*semantic.OuterJoinExpression)
			assert.IsType(t, &semantic.DotExpression{}, join.Expr)
			dot = join.Expr.(*semantic.DotExpression)
			assert.IsType(t, &semantic.NameExpression{}, dot.Name)
			name = dot.Name.(*semantic.NameExpression)
			assert.Equal(t, "id", name.Name)
			assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
			dot = dot.Parent.(*semantic.DotExpression)
			assert.IsType(t, &semantic.NameExpression{}, dot.Name)
			name = dot.Name.(*semantic.NameExpression)
			assert.Equal(t, "t2", name.Name)
			assert.Nil(t, dot.Parent)
		},
	})

	tests = append(tests, testCase{
		name: "commit & rollback",
		text: `commit;
rollback;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			assert.IsType(t, &semantic.CommitStatement{}, node.Statements[0])
			assert.IsType(t, &semantic.RollbackStatement{}, node.Statements[1])
		},
	})

	tests = append(tests, testCase{
		name: "delete",
		text: `delete from t1 where t1.id =1;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.DeleteStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.DeleteStatement)
			assert.NotNil(t, stmt.Table)
			assert.NotNil(t, stmt.Where)
		},
	})

	tests = append(tests, testCase{
		name: "update",
		text: `update t1 set id = 2 where t1.id =1;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.UpdateStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.UpdateStatement)
			assert.NotNil(t, stmt.Table)
			assert.NotNil(t, stmt.Where)
			assert.NotNil(t, stmt.SetExprs)
			assert.Equal(t, 1, len(stmt.SetExprs))
			assert.IsType(t, &semantic.BinaryExpression{}, stmt.SetExprs[0])
			expr := stmt.SetExprs[0].(*semantic.BinaryExpression)
			assert.Equal(t, "id", expr.Left.(*semantic.NameExpression).Name)
			assert.Equal(t, int64(2), expr.Right.(*semantic.NumericLiteral).Value)
		},
	})

	tests = append(tests, testCase{
		name: "drop function",
		text: `drop function test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.DropFunctionStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.DropFunctionStatement)
			assert.Equal(t, "test", stmt.Name)
		},
	})

	tests = append(tests, testCase{
		name: "drop procedure",
		text: `drop procedure test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.DropProcedureStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.DropProcedureStatement)
			assert.Equal(t, "test", stmt.Name)
		},
	})

	tests = append(tests, testCase{
		name: "drop package",
		text: `drop package body a.test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.DropPackageStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.DropPackageStatement)
			assert.Equal(t, "test", stmt.Name)
			assert.Equal(t, "a", stmt.Schema)
			assert.True(t, stmt.IsBody)
		},
	})
	tests = append(tests, testCase{
		name: "drop trigger",
		text: `drop trigger test;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.DropTriggerStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.DropTriggerStatement)
			assert.Equal(t, "test", stmt.Name)
		},
	})

	runTestSuite(t, tests)
}

func TestParseFunctionCall(t *testing.T) {
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "functions",
		text: `select to_char(id) a, to_char(id, 'DD-MM-YYYY'), to_char(id, 'DD-MM-YYYY', 'HH24:MI:SS'),
decode(m.move_kind||m.order_type,'LOADDELIVER',wi1.wi_dest_loc,'') from t1;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, len(stmt.Fields.Fields), 4)
			{ // to_char(id) a
				assert.NotNil(t, stmt.Fields.Fields[0].Expr)
				assert.IsType(t, &semantic.AliasExpression{}, stmt.Fields.Fields[0].Expr)
				alias := stmt.Fields.Fields[0].Expr.(*semantic.AliasExpression)
				assert.Equal(t, "a", alias.Alias)
				assert.IsType(t, &semantic.FunctionCallExpression{}, alias.Expr)
				expr := alias.Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "TO_CHAR", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
			{ // to_char(id, 'DD-MM-YYYY')
				assert.NotNil(t, stmt.Fields.Fields[1].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[1].Expr)
				expr := stmt.Fields.Fields[1].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "TO_CHAR", name.Name)
				assert.Equal(t, len(expr.Args), 2)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[1])
				str := expr.Args[1].(*semantic.StringLiteral)
				assert.Equal(t, "'DD-MM-YYYY'", str.Value)
			}
			{ // to_char(id, 'DD-MM-YYYY', 'HH24:MI:SS')
				assert.NotNil(t, stmt.Fields.Fields[2].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[2].Expr)
				expr := stmt.Fields.Fields[2].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "TO_CHAR", name.Name)
				assert.Equal(t, len(expr.Args), 3)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[1])
				str := expr.Args[1].(*semantic.StringLiteral)
				assert.Equal(t, "'DD-MM-YYYY'", str.Value)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[2])
				str = expr.Args[2].(*semantic.StringLiteral)
				assert.Equal(t, "'HH24:MI:SS'", str.Value)
			}
			{ // decode(m.move_kind||m.order_type,'LOADDELIVER',wi1.wi_dest_loc,'')
				i := 3
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "DECODE", name.Name)
				assert.Equal(t, len(expr.Args), 4)
				assert.IsType(t, &semantic.BinaryExpression{}, expr.Args[0])
				binary := expr.Args[0].(*semantic.BinaryExpression)
				assert.IsType(t, &semantic.DotExpression{}, binary.Left)
				assert.IsType(t, &semantic.DotExpression{}, binary.Right)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[1])
				str := expr.Args[1].(*semantic.StringLiteral)
				assert.Equal(t, "'LOADDELIVER'", str.Value)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[3])
				str = expr.Args[3].(*semantic.StringLiteral)
				assert.Equal(t, "''", str.Value)
			}
		},
	})

	tests = append(tests, testCase{
		name: "count",
		text: `select count(*) a, count(1), count(id), count(t.id)
-- , count(distinct t.d),count(unique t.id), count(all t.id)
-- ,count(id) over partition by t.id
from t;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, len(stmt.Fields.Fields), 4)
			{ // count(*) a
				assert.NotNil(t, stmt.Fields.Fields[0].Expr)
				assert.IsType(t, &semantic.AliasExpression{}, stmt.Fields.Fields[0].Expr)
				alias := stmt.Fields.Fields[0].Expr.(*semantic.AliasExpression)
				assert.Equal(t, "a", alias.Alias)
				assert.IsType(t, &semantic.FunctionCallExpression{}, alias.Expr)
				expr := alias.Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "COUNT", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.StringLiteral{}, expr.Args[0])
				str := expr.Args[0].(*semantic.StringLiteral)
				assert.Equal(t, "*", str.Value)
			}
			{ // count(1)
				assert.NotNil(t, stmt.Fields.Fields[1].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[1].Expr)
				expr := stmt.Fields.Fields[1].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "COUNT", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.NumericLiteral{}, expr.Args[0])
				num := expr.Args[0].(*semantic.NumericLiteral)
				assert.Equal(t, int64(1), num.Value)
			}
			{ // count(id)
				assert.NotNil(t, stmt.Fields.Fields[2].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[2].Expr)
				expr := stmt.Fields.Fields[2].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "COUNT", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
			{ // count(t.id)
				i := 3
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "COUNT", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.DotExpression{}, expr.Args[0])
				dot := expr.Args[0].(*semantic.DotExpression)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
		},
	})

	tests = append(tests, testCase{
		name: "min",
		text: `select min(id), min(t.id)
from t;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, len(stmt.Fields.Fields), 2)
			{ // min(id)
				i := 0
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "min", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
			{ // min(t.id)
				i := 1
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "min", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.DotExpression{}, expr.Args[0])
				dot := expr.Args[0].(*semantic.DotExpression)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
		},
	})

	tests = append(tests, testCase{
		name: "trim",
		text: `select trim(id), avg(t.id)
from t;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 1, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, 2, len(stmt.Fields.Fields))
			{ // trim(id)
				i := 0
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "TRIM", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.NameExpression{}, expr.Args[0])
				name = expr.Args[0].(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
			{ // avg(t.id)
				i := 1
				assert.NotNil(t, stmt.Fields.Fields[i].Expr)
				assert.IsType(t, &semantic.FunctionCallExpression{}, stmt.Fields.Fields[i].Expr)
				expr := stmt.Fields.Fields[i].Expr.(*semantic.FunctionCallExpression)
				assert.IsType(t, &semantic.NameExpression{}, expr.Name)
				name := expr.Name.(*semantic.NameExpression)
				assert.Equal(t, "AVG", name.Name)
				assert.Equal(t, len(expr.Args), 1)
				assert.IsType(t, &semantic.DotExpression{}, expr.Args[0])
				dot := expr.Args[0].(*semantic.DotExpression)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "id", name.Name)
			}
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
		return(a);
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
			assert.IsType(t, &semantic.NameExpression{}, r.Name)
			name := r.Name.(*semantic.NameExpression)
			assert.Equal(t, "a", name.Name)
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
return null;
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
			assert.IsType(t, &semantic.NullExpression{}, r.Name)
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
			assert.Equal(t, semantic.Span{Start: 0, End: 1711}, stmt.Span())
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
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp = dotExp.Name.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Aws_Asc_Id")
						assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
						dotExp = dotExp.Parent.(*semantic.DotExpression)
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp = dotExp.Name.(*semantic.NameExpression)
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
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp = dotExp.Name.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "Aws_Asc_Id")
						assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
						dotExp = dotExp.Parent.(*semantic.DotExpression)
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp = dotExp.Name.(*semantic.NameExpression)
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
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
							nameExp := dotExp.Name.(*semantic.NameExpression)
							assert.Equal(t, nameExp.Name, "Aws_Asc_Id")
							assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
							dotExp = dotExp.Parent.(*semantic.DotExpression)
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
							relExp := dotExp.Name.(*semantic.NameExpression)
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
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
							nameExp = dotExp.Name.(*semantic.NameExpression)
							assert.Equal(t, nameExp.Name, "Aws_Asc_Id")
							assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
							dotExp = dotExp.Parent.(*semantic.DotExpression)
							assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
							nameExp = dotExp.Name.(*semantic.NameExpression)
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
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp := dotExp.Name.(*semantic.NameExpression)
						assert.Equal(t, nameExp.Name, "gen_Anchor_Job_P")
						assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
						dotExp = dotExp.Parent.(*semantic.DotExpression)
						assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
						nameExp = dotExp.Name.(*semantic.NameExpression)
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
	tests = append(tests, testCase{
		name: "bind variable",
		root: getBlock,
		text: `
DECLARE
	a NUMBER := 1;
BEGIN
	a:=:New.id.c;
	a:=:Old.id;
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
			assert.Equal(t, len(node.Body.Statements), 2)
			{ // a:=:New.id.c;
				i := 0
				assert.IsType(t, &semantic.AssignmentStatement{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.AssignmentStatement)
				assert.Equal(t, stmt.Left, "a")
				assert.IsType(t, &semantic.DotExpression{}, stmt.Right)
				dotExp := stmt.Right.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
				name := dotExp.Name.(*semantic.NameExpression)
				assert.Equal(t, name.Name, "c")
				assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
				dotExp = dotExp.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
				name = dotExp.Name.(*semantic.NameExpression)
				assert.Equal(t, name.Name, "id")
				assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
				dotExp = dotExp.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.BindNameExpression{}, dotExp.Name)
				bindExp := dotExp.Name.(*semantic.BindNameExpression)
				assert.IsType(t, &semantic.NameExpression{}, bindExp.Name)
				nameExp := bindExp.Name.(*semantic.NameExpression)
				assert.Equal(t, nameExp.Name, ":New")
			}
			{ // a:=:Old.id;
				i := 1
				assert.IsType(t, &semantic.AssignmentStatement{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.AssignmentStatement)
				assert.Equal(t, stmt.Left, "a")
				assert.IsType(t, &semantic.DotExpression{}, stmt.Right)
				dotExp := stmt.Right.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
				name := dotExp.Name.(*semantic.NameExpression)
				assert.Equal(t, name.Name, "id")
				assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
				dotExp = dotExp.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.BindNameExpression{}, dotExp.Name)
				bindExp := dotExp.Name.(*semantic.BindNameExpression)
				assert.IsType(t, &semantic.NameExpression{}, bindExp.Name)
				nameExp := bindExp.Name.(*semantic.NameExpression)
				assert.Equal(t, nameExp.Name, ":Old")
			}
		},
	})
	tests = append(tests, testCase{
		name: "execute_immediate",
		root: getBlock,
		text: `
DECLARE
	a NUMBER := 1;
BEGIN
	execute immediate 'select * from t';
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
			{ // execute immediate 'select * from t';
				i := 0
				assert.IsType(t, &semantic.ExecuteImmediateStatement{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.ExecuteImmediateStatement)
				assert.Equal(t, "'select * from t'", stmt.Sql)
			}
		},
	})
	tests = append(tests, testCase{
		name: "raise exception",
		root: getBlock,
		text: `
BEGIN
	raise test;
END`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.BlockStatement)
			assert.NotNil(t, node.Body)
			assert.Equal(t, len(node.Body.Statements), 1)
			{
				i := 0
				assert.IsType(t, &semantic.RaiseStatement{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.RaiseStatement)
				assert.Equal(t, "test", stmt.Name)
			}
		},
	})
	tests = append(tests, testCase{
		name: "goto & label",
		root: getBlock,
		text: `
BEGIN
	<<test>>
    goto test;
END`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.BlockStatement)
			assert.NotNil(t, node.Body)
			assert.Equal(t, len(node.Body.Statements), 2)
			{
				i := 0
				assert.IsType(t, &semantic.LabelDeclaration{}, node.Body.Statements[i])
				i++
				assert.IsType(t, &semantic.GotoStatement{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.GotoStatement)
				assert.Equal(t, "test", stmt.Label)
			}
		},
	})
	tests = append(tests, testCase{
		name: "named argument",
		root: getBlock,
		text: `
BEGIN
	func1(a=>1);
	func2(a=>a,b=>2);
	func3(a=>t.id);
END`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.BlockStatement)
			assert.NotNil(t, node.Body)
			assert.Equal(t, 3, len(node.Body.Statements))
			{
				i := 0
				assert.IsType(t, &semantic.ProcedureCall{}, node.Body.Statements[i])
				stmt := node.Body.Statements[i].(*semantic.ProcedureCall)
				assert.IsType(t, &semantic.NameExpression{}, stmt.Name)
				nameExp := stmt.Name.(*semantic.NameExpression)
				assert.Equal(t, "func1", nameExp.Name)
				assert.Equal(t, 1, len(stmt.Arguments))
				assert.IsType(t, &semantic.NamedArgumentExpression{}, stmt.Arguments[0])
				arg := stmt.Arguments[0].(*semantic.NamedArgumentExpression)
				assert.IsType(t, &semantic.NameExpression{}, arg.Name)
				nameExp = arg.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", nameExp.Name)
				assert.IsType(t, &semantic.NumericLiteral{}, arg.Value)
				assert.Equal(t, int64(1), arg.Value.(*semantic.NumericLiteral).Value)
				i++
				assert.IsType(t, &semantic.ProcedureCall{}, node.Body.Statements[i])
				stmt = node.Body.Statements[i].(*semantic.ProcedureCall)
				assert.IsType(t, &semantic.NameExpression{}, stmt.Name)
				nameExp = stmt.Name.(*semantic.NameExpression)
				assert.Equal(t, "func2", nameExp.Name)
				assert.Equal(t, 2, len(stmt.Arguments))
				assert.IsType(t, &semantic.NamedArgumentExpression{}, stmt.Arguments[0])
				arg = stmt.Arguments[0].(*semantic.NamedArgumentExpression)
				assert.IsType(t, &semantic.NameExpression{}, arg.Name)
				nameExp = arg.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", nameExp.Name)
				assert.IsType(t, &semantic.NameExpression{}, arg.Value)
				nameExp = arg.Value.(*semantic.NameExpression)
				assert.Equal(t, "a", nameExp.Name)
				assert.IsType(t, &semantic.NamedArgumentExpression{}, stmt.Arguments[1])
				arg = stmt.Arguments[1].(*semantic.NamedArgumentExpression)
				assert.IsType(t, &semantic.NameExpression{}, arg.Name)
				nameExp = arg.Name.(*semantic.NameExpression)
				assert.Equal(t, "b", nameExp.Name)
				assert.IsType(t, &semantic.NumericLiteral{}, arg.Value)
				assert.Equal(t, int64(2), arg.Value.(*semantic.NumericLiteral).Value)
				i++
				assert.IsType(t, &semantic.ProcedureCall{}, node.Body.Statements[i])
				stmt = node.Body.Statements[i].(*semantic.ProcedureCall)
				assert.IsType(t, &semantic.NameExpression{}, stmt.Name)
				nameExp = stmt.Name.(*semantic.NameExpression)
				assert.Equal(t, "func3", nameExp.Name)
				assert.Equal(t, 1, len(stmt.Arguments))
				assert.IsType(t, &semantic.NamedArgumentExpression{}, stmt.Arguments[0])
				arg = stmt.Arguments[0].(*semantic.NamedArgumentExpression)
				assert.IsType(t, &semantic.NameExpression{}, arg.Name)
				nameExp = arg.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", nameExp.Name)
				assert.IsType(t, &semantic.DotExpression{}, arg.Value)
				dotExp := arg.Value.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
				nameExp = dotExp.Name.(*semantic.NameExpression)
				assert.Equal(t, "id", nameExp.Name)
				assert.IsType(t, &semantic.DotExpression{}, dotExp.Parent)
				dotExp = dotExp.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dotExp.Name)
				nameExp = dotExp.Name.(*semantic.NameExpression)
				assert.Equal(t, "t", nameExp.Name)
			}
		},
	})

	runTestSuite(t, tests)
}

func TestCaseWhenStatement(t *testing.T) {
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "searched case statement",
		text: `
select case when a=b then 1 else 2 end from dual;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Greater(t, len(node.Statements), 0)
			stmt, ok := node.Statements[0].(*semantic.SelectStatement)
			assert.True(t, ok)
			assert.NotNil(t, stmt)
			assert.Equal(t, 2, stmt.Line())
			assert.Equal(t, 1, stmt.Column())
			assert.Equal(t, len(stmt.Fields.Fields), 1)
			assert.IsType(t, &semantic.StatementExpression{}, stmt.Fields.Fields[0].Expr)
			stmtExpr := stmt.Fields.Fields[0].Expr.(*semantic.StatementExpression)
			assert.IsType(t, &semantic.CaseWhenStatement{}, stmtExpr.Stmt)
			caseStmt := stmtExpr.Stmt.(*semantic.CaseWhenStatement)
			assert.Nil(t, caseStmt.Expr)
			assert.Equal(t, 1, len(caseStmt.WhenClauses))
			assert.IsType(t, &semantic.CaseWhenBlock{}, caseStmt.WhenClauses[0])
			whenBlock := caseStmt.WhenClauses[0]
			assert.NotNil(t, whenBlock.Expr)
			assert.Equal(t, 0, len(whenBlock.Stmts))
			assert.NotNil(t, caseStmt.ElseClause.Expr)
			assert.Equal(t, 0, len(caseStmt.ElseClause.Stmts))
		},
	})

	runTestSuite(t, tests)
}

func TestParseSelectStatement(t *testing.T) {
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "select for update",
		text: `select * from test for update nowait;`,
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
			assert.Equal(t, len(stmt.From.TableRefs), 1)
			assert.Equal(t, stmt.From.TableRefs[0].Table, "test")
			assert.NotNil(t, stmt.ForUpdate)
			assert.NotNil(t, stmt.ForUpdate.Options)
		},
	})

	runTestSuite(t, tests)
}

func TestParseInsertStatement(t *testing.T) {
	tests := testSuite{}

	tests = append(tests, testCase{
		name: "simple insert",
		text: `insert into t1 values (1);`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.InsertStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.InsertStatement)
			assert.NotNil(t, stmt.AllInto)
			assert.Equal(t, 1, len(stmt.AllInto))
			into := stmt.AllInto[0]
			assert.NotNil(t, into.Table)
			assert.Equal(t, "t1", into.Table.Table)
			assert.NotNil(t, into.Values)
			assert.Equal(t, 1, len(into.Values))
			assert.IsType(t, &semantic.NumericLiteral{}, into.Values[0])
			num := into.Values[0].(*semantic.NumericLiteral)
			assert.Equal(t, int64(1), num.Value)
		},
	})

	tests = append(tests, testCase{
		name: " insert table with columns",
		text: `insert into t1 (a,b,c) values (1,2,3);`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.InsertStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.InsertStatement)
			assert.NotNil(t, stmt.AllInto)
			assert.Equal(t, 1, len(stmt.AllInto))
			into := stmt.AllInto[0]
			assert.NotNil(t, into.Table)
			assert.Equal(t, "t1", into.Table.Table)
			assert.Equal(t, 3, len(into.Columns))
			assert.IsType(t, &semantic.NameExpression{}, into.Columns[0])
			name := into.Columns[0].(*semantic.NameExpression)
			assert.Equal(t, name.Name, "a")
			assert.IsType(t, &semantic.NameExpression{}, into.Columns[1])
			name = into.Columns[1].(*semantic.NameExpression)
			assert.Equal(t, name.Name, "b")
			assert.IsType(t, &semantic.NameExpression{}, into.Columns[2])
			name = into.Columns[2].(*semantic.NameExpression)
			assert.Equal(t, name.Name, "c")
			assert.NotNil(t, into.Values)
			assert.Equal(t, 3, len(into.Values))
			assert.IsType(t, &semantic.NumericLiteral{}, into.Values[0])
			num := into.Values[0].(*semantic.NumericLiteral)
			assert.Equal(t, int64(1), num.Value)
			assert.IsType(t, &semantic.NumericLiteral{}, into.Values[1])
			num = into.Values[1].(*semantic.NumericLiteral)
			assert.Equal(t, int64(2), num.Value)
		},
	})

	tests = append(tests, testCase{
		name: " insert table with select",
		text: `insert into t1 (a,b,c) select a,b,c from t;`,
		Func: func(t *testing.T, root any) {
			node := root.(*semantic.Script)
			assert.Equal(t, len(node.Statements), 1)
			assert.IsType(t, &semantic.InsertStatement{}, node.Statements[0])
			stmt := node.Statements[0].(*semantic.InsertStatement)
			assert.NotNil(t, stmt.AllInto)
			assert.Equal(t, 1, len(stmt.AllInto))
			into := stmt.AllInto[0]
			assert.NotNil(t, into.Table)
			assert.Equal(t, "t1", into.Table.Table)
			assert.Nil(t, into.Values)
			assert.NotNil(t, stmt.Select)
			assert.Equal(t, 3, len(stmt.Select.Fields.Fields))
		},
	})

	runTestSuite(t, tests)
}

func TestParseMergeStatement(t *testing.T) {
	//tests := testSuite{}

	tests := []testCase{
		{
			name: "simple merge",
			text: `
merge into t1
using t2
on (t1.a=t2.a)
when matched then
	update set t1.b=t2.b;`,
			Func: func(t *testing.T, root any) {
				node := root.(*semantic.Script)
				assert.Equal(t, len(node.Statements), 1)
				assert.IsType(t, &semantic.MergeStatement{}, node.Statements[0])
				stmt := node.Statements[0].(*semantic.MergeStatement)
				assert.Equal(t, "t1", stmt.Table.Table)
				assert.IsType(t, &semantic.NameExpression{}, stmt.Using)
				name := stmt.Using.(*semantic.NameExpression)
				assert.Equal(t, name.Name, "t2")
				assert.IsType(t, &semantic.RelationalExpression{}, stmt.OnCondition)
				rel := stmt.OnCondition.(*semantic.RelationalExpression)
				assert.Equal(t, "=", rel.Operator)
				assert.IsType(t, &semantic.DotExpression{}, rel.Left)
				dot := rel.Left.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", name.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t1", name.Name)
				assert.IsType(t, &semantic.DotExpression{}, rel.Right)
				dot = rel.Right.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", name.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				name = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t2", name.Name)
				assert.Nil(t, stmt.MergeInsert)
				assert.NotNil(t, stmt.MergeUpdate)
			},
		},
		{
			name: "merge using select",
			text: `
merge into t1
using (select * from t2)
on (t1.a=t2.a)
when matched then
	update set t1.b=t2.b;`,
			Func: func(t *testing.T, root any) {
				node := root.(*semantic.Script)
				assert.Equal(t, len(node.Statements), 1)
				assert.IsType(t, &semantic.MergeStatement{}, node.Statements[0])
				stmt := node.Statements[0].(*semantic.MergeStatement)
				assert.Equal(t, "t1", stmt.Table.Table)
				assert.IsType(t, &semantic.StatementExpression{}, stmt.Using)
				exprStmt := stmt.Using.(*semantic.StatementExpression)
				assert.IsType(t, &semantic.SelectStatement{}, exprStmt.Stmt)
				assert.IsType(t, &semantic.RelationalExpression{}, stmt.OnCondition)
				rel := stmt.OnCondition.(*semantic.RelationalExpression)
				assert.Equal(t, "=", rel.Operator)
				assert.IsType(t, &semantic.DotExpression{}, rel.Left)
				dot := rel.Left.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr := dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t1", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, rel.Right)
				dot = rel.Right.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "a", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t2", expr.Name)
				assert.Nil(t, stmt.MergeInsert)
				assert.NotNil(t, stmt.MergeUpdate)
				update := stmt.MergeUpdate
				assert.NotNil(t, update)
				assert.Equal(t, 1, len(update.SetElems))
				assert.IsType(t, &semantic.BinaryExpression{}, update.SetElems[0])
				binary := update.SetElems[0].(*semantic.BinaryExpression)
				assert.Equal(t, "=", binary.Operator)
				assert.IsType(t, &semantic.DotExpression{}, binary.Left)
				dot = binary.Left.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "b", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t1", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, binary.Right)
				dot = binary.Right.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "b", expr.Name)
				assert.IsType(t, &semantic.DotExpression{}, dot.Parent)
				dot = dot.Parent.(*semantic.DotExpression)
				assert.IsType(t, &semantic.NameExpression{}, dot.Name)
				expr = dot.Name.(*semantic.NameExpression)
				assert.Equal(t, "t2", expr.Name)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.root != nil {
				runTest(t, test.text, test.Func, test.root)
			} else {
				runTest(t, test.text, test.Func)
			}
		})
	}

	//runTestSuite(t, tests)
}
