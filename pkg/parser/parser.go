package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
)

type (
	sqlListener struct {
		*plsql.BasePlSqlParserListener

		script    *semantic.Script
		nodeStack semantic.Stack[semantic.Node]
	}
)

func GeneralScript(root plsql.ISql_scriptContext) *semantic.Script {
	listener := &sqlListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, root)
	return listener.script
}

func peekNode[T semantic.Node](l *sqlListener) (t T, err error) {
	for i := len(l.nodeStack) - 1; i >= 0; i-- {
		node := l.nodeStack[i]
		if _, ok := node.(T); ok {
			return node.(T), nil
		}
	}
	return t, fmt.Errorf("node is not of type %T", t)
}

func (l *sqlListener) EnterSql_script(ctx *plsql.Sql_scriptContext) {
	l.script = &semantic.Script{
		Statements: make([]semantic.Statement, 0),
	}
	l.script.SetLine(ctx.GetStart().GetLine())
	l.script.SetColumn(ctx.GetStart().GetColumn())
	l.nodeStack.Push(l.script)
}

func (l *sqlListener) ExitSql_script(ctx *plsql.Sql_scriptContext) {
	l.nodeStack.Pop()
}

func (l *sqlListener) ExitUnit_statement(ctx *plsql.Unit_statementContext) {
	for {
		node := l.nodeStack.Top()
		if _, ok := node.(*semantic.Script); ok {
			break
		}
		l.script.Statements = append(l.script.Statements, node.(semantic.Statement))
		l.nodeStack.Pop()
	}
}

func (l *sqlListener) EnterSelect_statement(ctx *plsql.Select_statementContext) {
	stmt := &semantic.SelectStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitSelect_statement(ctx *plsql.Select_statementContext) {
}

func (l *sqlListener) EnterSelected_list(ctx *plsql.Selected_listContext) {
	stmt := l.nodeStack.Top().(*semantic.SelectStatement)
	stmt.Fields = &semantic.FieldList{
		Fields: make([]*semantic.SelectField, 0),
	}
	if ctx.ASTERISK() != nil {
		field := &semantic.WildCardField{}
		field.Table = ctx.ASTERISK().GetText()
		selectField := &semantic.SelectField{
			WildCard: field,
		}
		stmt.Fields.Fields = append(stmt.Fields.Fields, selectField)
	}
}

func (l *sqlListener) EnterTable_ref_list(ctx *plsql.Table_ref_listContext) {
	from := &semantic.FromClause{
		TableRefs: make([]*semantic.TableRef, 0),
	}
	l.nodeStack.Push(from)
}

func (l *sqlListener) ExitTable_ref_list(ctx *plsql.Table_ref_listContext) {
	from := l.nodeStack.Pop().(*semantic.FromClause)
	stmt := l.nodeStack.Top().(*semantic.SelectStatement)

	for i := 0; i < ctx.GetChildCount(); i++ {
		if ctx.Table_ref(i) != nil {
			from.TableRefs = append(from.TableRefs, &semantic.TableRef{
				Table: ctx.Table_ref(i).GetText(),
			})
		}
		stmt.From = from
	}
}

func (l *sqlListener) EnterCreate_procedure_body(ctx *plsql.Create_procedure_bodyContext) {
	stmt := &semantic.CreateProcedureStatement{}
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitCreate_procedure_body(ctx *plsql.Create_procedure_bodyContext) {
	stmt, err := peekNode[*semantic.CreateProcedureStatement](l)
	if err != nil {
		panic(err)
	}

	// set line & column
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())

	// set name
	stmt.Name = ctx.Procedure_name().GetText()

	// set is replace
	stmt.IsReplace = ctx.REPLACE() != nil

	// Add statements
	for {
		node := l.nodeStack.Top()
		if _, ok := node.(*semantic.CreateProcedureStatement); ok {
			break
		}
		switch node.(type) {
		case *semantic.Body:
			stmt.Body = node.(*semantic.Body)
		case semantic.Declaration:
			stmt.Declarations = append([]semantic.Declaration{node.(semantic.Declaration)}, stmt.Declarations...)
		case *semantic.Parameter:
			stmt.Parameters = append([]*semantic.Parameter{node.(*semantic.Parameter)}, stmt.Parameters...)
		}
		l.nodeStack.Pop()
	}
}

func (l *sqlListener) EnterBody(ctx *plsql.BodyContext) {
	body := &semantic.Body{}
	l.nodeStack.Push(body)
}

func (l *sqlListener) ExitBody(ctx *plsql.BodyContext) {
	body, err := peekNode[*semantic.Body](l)
	if err != nil {
		panic(err)
	}
	for {
		node := l.nodeStack.Top()
		if _, ok := node.(*semantic.Body); ok {
			break
		}
		switch node.(type) {
		case semantic.Statement:
			body.Statements = append([]semantic.Statement{node.(semantic.Statement)}, body.Statements...)
		}
		l.nodeStack.Pop()
	}
}

func (l *sqlListener) EnterAssignment_statement(ctx *plsql.Assignment_statementContext) {
	stmt := &semantic.AssignmentStatement{}
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitAssignment_statement(ctx *plsql.Assignment_statementContext) {
	stmt, err := peekNode[*semantic.AssignmentStatement](l)
	if err != nil {
		panic(err)
	}

	// set line & column
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())

	// set left
	stmt.Left = ctx.General_element().GetText()
	stmt.Right = ctx.Expression().GetText()
}

func (l *sqlListener) ExitVariable_declaration(ctx *plsql.Variable_declarationContext) {
	stmt := &semantic.VariableDeclaration{}
	// set line & column
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Identifier().GetText()
	stmt.DataType = ctx.Type_spec().GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitException_declaration(ctx *plsql.Exception_declarationContext) {
	stmt := &semantic.ExceptionDeclaration{}
	// set line & column
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Identifier().GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitParameter(ctx *plsql.ParameterContext) {
	stmt := &semantic.Parameter{}
	stmt.Name = ctx.Parameter_name().GetText()
	stmt.DataType = ctx.Type_spec().GetText()
	l.nodeStack.Push(stmt)
}
