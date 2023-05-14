//go:build visitor

package parser

import (
	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type (
	plsqlVisitor struct {
		plsql.BasePlSqlParserVisitor
	}
)

func GeneralScript(root plsql.ISql_scriptContext) *semantic.Script {
	visitor := &plsqlVisitor{}
	script := visitor.VisitSql_script(root.(*plsql.Sql_scriptContext)).(*semantic.Script)
	return script
}

func (v *plsqlVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	children := node.GetChildren()
	nodes := make([]interface{}, 0, len(children))
	for _, child := range children {
		switch child.(type) {
		case *plsql.Query_blockContext:
			c := child.(*plsql.Query_blockContext)
			nodes = append(nodes, v.VisitQuery_block(c))
		case *plsql.Selected_listContext:
			c := child.(*plsql.Selected_listContext)
			nodes = append(nodes, v.VisitSelected_list(c))
		case *plsql.Table_ref_listContext:
			c := child.(*plsql.Table_ref_listContext)
			nodes = append(nodes, v.VisitTable_ref_list(c))
		case antlr.TerminalNode:
			break
		default:
			tree := child.(antlr.ParseTree)
			c := tree.Accept(v)
			nodes = append(nodes, c)
		}
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return nodes
}

func (v *plsqlVisitor) VisitSql_script(ctx *plsql.Sql_scriptContext) interface{} {
	script := &semantic.Script{}
	for _, stmt := range ctx.AllUnit_statement() {
		script.Statements = append(script.Statements, stmt.Accept(v).(semantic.Statement))
	}

	return script
}

func (v *plsqlVisitor) VisitQuery_block(ctx *plsql.Query_blockContext) interface{} {
	stmt := &semantic.SelectStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Fields = v.VisitSelected_list(ctx.Selected_list().(*plsql.Selected_listContext)).(*semantic.FieldList)
	stmt.From = ctx.From_clause().Accept(v).(*semantic.FromClause)
	return stmt
}

func (v *plsqlVisitor) VisitSelected_list(ctx *plsql.Selected_listContext) interface{} {
	fields := &semantic.FieldList{}
	fields.SetLine(ctx.GetStart().GetLine())
	fields.SetColumn(ctx.GetStart().GetColumn())
	if ctx.ASTERISK() != nil {
		fields.Fields = append(
			fields.Fields,
			&semantic.SelectField{WildCard: &semantic.WildCardField{Table: "*", Schema: "*"}},
		)
	}

	return fields
}

func (v *plsqlVisitor) VisitTable_ref_list(ctx *plsql.Table_ref_listContext) interface{} {
	from := &semantic.FromClause{}
	from.SetLine(ctx.GetStart().GetLine())
	from.SetColumn(ctx.GetStart().GetColumn())
	//tables := make([]*semantic.TableRef, 0, len(ctx.AllTable_ref()))
	for _, t := range ctx.AllTable_ref() {
		from.TableRefs = append(from.TableRefs, &semantic.TableRef{
			Table: t.GetText(),
		})
	}
	return from
}
