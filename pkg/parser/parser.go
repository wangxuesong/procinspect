package parser

import (
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

func (l *sqlListener) EnterSql_script(ctx *plsql.Sql_scriptContext) {
	l.script = &semantic.Script{
		Statements: make([]semantic.Statement, 0),
	}
	l.nodeStack.Push(l.script)
}

func (l *sqlListener) EnterSelect_statement(ctx *plsql.Select_statementContext) {
	stmt := &semantic.SelectStatement{}
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitSelect_statement(ctx *plsql.Select_statementContext) {
	stmt := l.nodeStack.Pop().(*semantic.SelectStatement)
	l.script.Statements = append(l.script.Statements, stmt)
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
