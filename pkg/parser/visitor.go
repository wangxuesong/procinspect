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
		case *plsql.Create_procedure_bodyContext:
			c := child.(*plsql.Create_procedure_bodyContext)
			nodes = append(nodes, v.VisitCreate_procedure_body(c))
		case *plsql.Variable_declarationContext:
			c := child.(*plsql.Variable_declarationContext)
			nodes = append(nodes, v.VisitVariable_declaration(c))
		case *plsql.Exception_declarationContext:
			c := child.(*plsql.Exception_declarationContext)
			nodes = append(nodes, v.VisitException_declaration(c))
		case *plsql.Cursor_declarationContext:
			c := child.(*plsql.Cursor_declarationContext)
			nodes = append(nodes, v.VisitCursor_declaration(c))
		case *plsql.Assignment_statementContext:
			c := child.(*plsql.Assignment_statementContext)
			nodes = append(nodes, v.VisitAssignment_statement(c))
		case *plsql.If_statementContext:
			c := child.(*plsql.If_statementContext)
			nodes = append(nodes, v.VisitIf_statement(c))
		case *plsql.Open_statementContext:
			c := child.(*plsql.Open_statementContext)
			nodes = append(nodes, v.VisitOpen_statement(c))
		case *plsql.Close_statementContext:
			c := child.(*plsql.Close_statementContext)
			nodes = append(nodes, v.VisitClose_statement(c))
		case *plsql.Fetch_statementContext:
			c := child.(*plsql.Fetch_statementContext)
			nodes = append(nodes, v.VisitFetch_statement(c))
		case *plsql.Exit_statementContext:
			c := child.(*plsql.Exit_statementContext)
			nodes = append(nodes, v.VisitExit_statement(c))
		case *plsql.Loop_statementContext:
			c := child.(*plsql.Loop_statementContext)
			nodes = append(nodes, v.VisitLoop_statement(c))
		case *plsql.Function_callContext:
			c := child.(*plsql.Function_callContext)
			nodes = append(nodes, v.VisitFunction_call(c))
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
	if ctx.Where_clause() != nil {
		if ctx.Where_clause().Expression() != nil {
			visitor := &exprVisitor{}
			stmt.Where = visitor.VisitExpression(ctx.Where_clause().Expression().(*plsql.ExpressionContext)).(semantic.Expr)
		}
	}
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

func (v *plsqlVisitor) VisitCreate_procedure_body(ctx *plsql.Create_procedure_bodyContext) interface{} {
	stmt := &semantic.CreateProcedureStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Procedure_name().GetText()
	stmt.IsReplace = ctx.REPLACE() != nil
	for _, p := range ctx.AllParameter() {
		stmt.Parameters = append(stmt.Parameters, v.VisitParameter(p.(*plsql.ParameterContext)).(*semantic.Parameter))
	}
	if ctx.Seq_of_declare_specs() != nil {
		stmt.Declarations = v.VisitSeq_of_declare_specs(ctx.Seq_of_declare_specs().(*plsql.Seq_of_declare_specsContext)).([]semantic.Declaration)
	}
	stmt.Body = v.VisitBody(ctx.Body().(*plsql.BodyContext)).(*semantic.Body)
	return stmt
}

func (v *plsqlVisitor) VisitParameter(ctx *plsql.ParameterContext) interface{} {
	param := &semantic.Parameter{}
	param.SetLine(ctx.GetStart().GetLine())
	param.SetColumn(ctx.GetStart().GetColumn())
	param.Name = ctx.Parameter_name().GetText()
	param.DataType = ctx.Type_spec().GetText()
	return param
}

func (v *plsqlVisitor) VisitSeq_of_declare_specs(ctx *plsql.Seq_of_declare_specsContext) interface{} {
	decls := make([]semantic.Declaration, 0, len(ctx.AllDeclare_spec()))
	for _, d := range ctx.AllDeclare_spec() {
		decls = append(decls, d.Accept(v).(semantic.Declaration))
	}
	return decls
}

func (v *plsqlVisitor) VisitVariable_declaration(ctx *plsql.Variable_declarationContext) interface{} {
	varDecl := &semantic.VariableDeclaration{}
	varDecl.SetLine(ctx.GetStart().GetLine())
	varDecl.SetColumn(ctx.GetStart().GetColumn())
	varDecl.Name = ctx.Identifier().GetText()
	varDecl.DataType = ctx.Type_spec().GetText()
	return varDecl
}

func (v *plsqlVisitor) VisitException_declaration(ctx *plsql.Exception_declarationContext) interface{} {
	exception := &semantic.ExceptionDeclaration{}
	exception.SetLine(ctx.GetStart().GetLine())
	exception.SetColumn(ctx.GetStart().GetColumn())
	exception.Name = ctx.Identifier().GetText()
	return exception
}

func (v *plsqlVisitor) VisitCursor_declaration(ctx *plsql.Cursor_declarationContext) interface{} {
	cursor := &semantic.CursorDeclaration{}
	cursor.SetLine(ctx.GetStart().GetLine())
	cursor.SetColumn(ctx.GetStart().GetColumn())
	cursor.Name = ctx.Identifier().GetText()
	for _, p := range ctx.AllParameter_spec() {
		cursor.Parameters = append(cursor.Parameters, v.VisitParameter_spec(p.(*plsql.Parameter_specContext)).(*semantic.Parameter))
	}
	cursor.Stmt = ctx.Select_statement().Accept(v).(*semantic.SelectStatement)
	return cursor
}

func (v *plsqlVisitor) VisitParameter_spec(ctx *plsql.Parameter_specContext) interface{} {
	param := &semantic.Parameter{}
	param.SetLine(ctx.GetStart().GetLine())
	param.SetColumn(ctx.GetStart().GetColumn())
	param.Name = ctx.Parameter_name().GetText()
	param.DataType = ctx.Type_spec().GetText()
	return param
}

func (v *plsqlVisitor) VisitBody(ctx *plsql.BodyContext) interface{} {
	stmt := &semantic.Body{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Statements = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	return stmt
}

func (v *plsqlVisitor) VisitSeq_of_statements(ctx *plsql.Seq_of_statementsContext) interface{} {
	stmts := make([]semantic.Statement, 0, len(ctx.AllStatement()))
	for _, stmt := range ctx.AllStatement() {
		stmts = append(stmts, stmt.Accept(v).(semantic.Statement))
	}
	return stmts
}

func (v *plsqlVisitor) VisitAssignment_statement(ctx *plsql.Assignment_statementContext) interface{} {
	stmt := &semantic.AssignmentStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Left = ctx.General_element().GetText()
	visitor := &exprVisitor{}
	stmt.Right = visitor.VisitExpression(ctx.Expression().(*plsql.ExpressionContext)).(semantic.Expr)

	return stmt
}

func (v *plsqlVisitor) VisitIf_statement(ctx *plsql.If_statementContext) interface{} {
	stmt := &semantic.IfStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	if ctx.Condition() != nil {
		vistior := exprVisitor{}
		stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
	}
	stmt.ThenBlock = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	if ctx.Else_part() != nil {
		stmt.ElseBlock = v.VisitSeq_of_statements(ctx.Else_part().Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	}
	return stmt
}

func (v *plsqlVisitor) VisitOpen_statement(ctx *plsql.Open_statementContext) interface{} {
	stmt := &semantic.OpenStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	return stmt
}

func (v *plsqlVisitor) VisitClose_statement(ctx *plsql.Close_statementContext) interface{} {
	stmt := &semantic.CloseStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	return stmt
}

func (v *plsqlVisitor) VisitFetch_statement(ctx *plsql.Fetch_statementContext) interface{} {
	stmt := &semantic.FetchStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Cursor = ctx.Cursor_name().GetText()
	stmt.Into = ctx.Variable_name(0).GetText()
	return stmt
}

func (v *plsqlVisitor) VisitExit_statement(ctx *plsql.Exit_statementContext) interface{} {
	stmt := &semantic.ExitStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	if ctx.Condition() != nil {
		vistior := exprVisitor{}
		stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
	}
	return stmt
}

func (v *plsqlVisitor) VisitLoop_statement(ctx *plsql.Loop_statementContext) interface{} {
	stmt := &semantic.LoopStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Statements = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	return stmt
}

func (v *plsqlVisitor) VisitFunction_call(ctx *plsql.Function_callContext) interface{} {
	stmt := &semantic.ProcedureCall{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Routine_name().GetText()
	if ctx.Function_argument() != nil {
		stmt.Arguments = v.VisitFunction_argument(ctx.Function_argument().(*plsql.Function_argumentContext)).([]*semantic.Argument)
	}
	return stmt
}

func (v *plsqlVisitor) VisitFunction_argument(ctx *plsql.Function_argumentContext) interface{} {
	args := make([]*semantic.Argument, 0, len(ctx.AllArgument()))
	for _, c := range ctx.AllArgument() {
		arg := &semantic.Argument{}
		arg.SetLine(c.GetStart().GetLine())
		arg.SetColumn(c.GetStart().GetColumn())

		arg.Name = c.GetText()
		args = append(args, arg)
	}
	return args
}