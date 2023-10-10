package parser

import (
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
)

type (
	plsqlVisitor struct {
		plsql.BasePlSqlParserVisitor
		StartLine int

		errors []ParseError
	}
)

func newAstNode[T any](ctx antlr.ParserRuleContext) *T {
	node := new(T)
	stmt := interface{}(node).(semantic.SetPosition)
	setAstSpan(ctx, stmt)
	return node
}

func setAstSpan(ctx antlr.ParserRuleContext, stmt semantic.SetPosition) {
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.SetSpan(semantic.Span{
		Start: ctx.GetStart().GetStart(),
		End:   ctx.GetStop().GetStop(),
	})
}

func newPlSqlVisitor(startLineNo ...int) *plsqlVisitor {
	v := &plsqlVisitor{}
	v.BasePlSqlParserVisitor.ParseTreeVisitor = v
	if len(startLineNo) == 1 {
		v.StartLine = startLineNo[0]
	}
	return v
}

func (v *plsqlVisitor) Error() (err error) {
	if len(v.errors) == 0 {
		return nil
	}
	for _, e := range v.errors {
		err = errors.Join(err, e)
	}
	return
}

func (v *plsqlVisitor) Errors() []ParseError {
	return v.errors
}

func (v *plsqlVisitor) ReportError(msg string, line, column int) {
	v.errors = append(v.errors, ParseError{
		Kind:   ErrSemantic,
		Line:   line + v.StartLine,
		Column: column,
		Msg:    msg,
	})
}

func (v *plsqlVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *plsqlVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	// TODO implement me
	panic("implement me")
}

func (v *plsqlVisitor) VisitErrorNode(node antlr.ErrorNode) interface{} {
	// TODO implement me
	panic("implement me")
}

func GeneralScript(root plsql.ISql_scriptContext) (script *semantic.Script, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	visitor := newPlSqlVisitor()
	script = visitor.VisitSql_script(root.(*plsql.Sql_scriptContext)).(*semantic.Script)
	return script, visitor.Error()
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
		case *plsql.Create_packageContext:
			c := child.(*plsql.Create_packageContext)
			nodes = append(nodes, v.VisitCreate_package(c))
		case *plsql.Create_package_bodyContext:
			c := child.(*plsql.Create_package_bodyContext)
			nodes = append(nodes, v.VisitCreate_package_body(c))
		case *plsql.Create_procedure_bodyContext:
			c := child.(*plsql.Create_procedure_bodyContext)
			nodes = append(nodes, v.VisitCreate_procedure_body(c))
		case *plsql.Create_function_bodyContext:
			c := child.(*plsql.Create_function_bodyContext)
			nodes = append(nodes, v.VisitCreate_function_body(c))
		case *plsql.Create_typeContext:
			c := child.(*plsql.Create_typeContext)
			nodes = append(nodes, v.VisitCreate_type(c))
		case *plsql.Type_definitionContext:
			c := child.(*plsql.Type_definitionContext)
			nodes = append(nodes, v.VisitType_definition(c))
		case *plsql.Nested_table_type_defContext:
			c := child.(*plsql.Nested_table_type_defContext)
			nodes = append(nodes, v.VisitNested_table_type_def(c))
		case *plsql.Type_declarationContext:
			c := child.(*plsql.Type_declarationContext)
			nodes = append(nodes, v.VisitType_declaration(c))
		case *plsql.Procedure_specContext:
			c := child.(*plsql.Procedure_specContext)
			nodes = append(nodes, v.VisitProcedure_spec(c))
		case *plsql.Procedure_bodyContext:
			c := child.(*plsql.Procedure_bodyContext)
			nodes = append(nodes, v.VisitProcedure_body(c))
		case *plsql.Function_specContext:
			c := child.(*plsql.Function_specContext)
			nodes = append(nodes, v.VisitFunction_spec(c))
		case *plsql.Function_bodyContext:
			c := child.(*plsql.Function_bodyContext)
			nodes = append(nodes, v.VisitFunction_body(c))
		case *plsql.Variable_declarationContext:
			c := child.(*plsql.Variable_declarationContext)
			nodes = append(nodes, v.VisitVariable_declaration(c))
		case *plsql.Exception_declarationContext:
			c := child.(*plsql.Exception_declarationContext)
			nodes = append(nodes, v.VisitException_declaration(c))
		case *plsql.Cursor_declarationContext:
			c := child.(*plsql.Cursor_declarationContext)
			nodes = append(nodes, v.VisitCursor_declaration(c))
		case *plsql.Anonymous_blockContext:
			c := child.(*plsql.Anonymous_blockContext)
			nodes = append(nodes, v.VisitAnonymous_block(c))
		case *plsql.Assignment_statementContext:
			c := child.(*plsql.Assignment_statementContext)
			nodes = append(nodes, v.VisitAssignment_statement(c))
		case *plsql.If_statementContext:
			c := child.(*plsql.If_statementContext)
			nodes = append(nodes, v.VisitIf_statement(c))
		case *plsql.Open_statementContext:
			c := child.(*plsql.Open_statementContext)
			nodes = append(nodes, v.VisitOpen_statement(c))
		case *plsql.Open_for_statementContext:
			c := child.(*plsql.Open_for_statementContext)
			nodes = append(nodes, v.VisitOpen_for_statement(c))
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
		case *plsql.Return_statementContext:
			c := child.(*plsql.Return_statementContext)
			nodes = append(nodes, v.VisitReturn_statement(c))
		case *plsql.Null_statementContext:
			c := child.(*plsql.Null_statementContext)
			nodes = append(nodes, v.VisitNull_statement(c))
		case *plsql.Call_statementContext:
			c := child.(*plsql.Call_statementContext)
			nodes = append(nodes, v.VisitCall_statement(c))
		case *plsql.Simple_case_statementContext:
			c := child.(*plsql.Simple_case_statementContext)
			nodes = append(nodes, v.VisitSimple_case_statement(c))
		case *plsql.Searched_case_statementContext:
			c := child.(*plsql.Searched_case_statementContext)
			nodes = append(nodes, v.VisitSearched_case_statement(c))
		case *plsql.Commit_statementContext:
			c := child.(*plsql.Commit_statementContext)
			nodes = append(nodes, v.VisitCommit_statement(c))
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
	script := newAstNode[semantic.Script](ctx)
	for _, stmt := range ctx.AllUnit_statement() {
		o := stmt.Accept(v)
		switch o.(type) {
		case semantic.Statement:
			break
		default:
			v.ReportError(fmt.Sprintf("unprocessed syntax %T",
				stmt.GetChild(0)),
				stmt.GetStart().GetLine(), stmt.GetStart().GetColumn())
			continue
		}
		script.Statements = append(script.Statements, o.(semantic.Statement))
	}

	return script
}

func (v *plsqlVisitor) VisitSelect_statement(ctx *plsql.Select_statementContext) interface{} {
	object := ctx.Select_only_statement().Accept(v)
	switch object.(type) {
	case *semantic.SelectStatement:
		stmt, ok := object.(*semantic.SelectStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Select_only_statement()),
				ctx.GetStart().GetLine(),
				ctx.GetStart().GetColumn())
		}
		if ctx.For_update_clause(0) != nil {
			expr, ok := ctx.For_update_clause(0).Accept(v).(*semantic.ForUpdateClause)
			if !ok {
				v.ReportError(fmt.Sprintf("unprocessed expression %T", ctx.For_update_clause(0)),
					ctx.For_update_clause(0).GetStart().GetLine(),
					ctx.For_update_clause(0).GetStart().GetColumn())
				return stmt
			}
			stmt.ForUpdate = expr
		}
		return stmt
	case *semantic.SetOperationStatement:
		stmt := object.(*semantic.SetOperationStatement)
		return stmt
	default:
		v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Select_only_statement()),
			ctx.GetStart().GetLine(),
			ctx.GetStart().GetColumn())
		return nil
	}
}

func (v *plsqlVisitor) VisitSelect_only_statement(ctx *plsql.Select_only_statementContext) interface{} {
	var with *semantic.WithClause
	if ctx.Subquery_factoring_clause() != nil {
		var ok bool
		with, ok = ctx.Subquery_factoring_clause().Accept(v).(*semantic.WithClause)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Subquery_factoring_clause()),
				ctx.Subquery_factoring_clause().GetStart().GetLine(),
				ctx.Subquery_factoring_clause().GetStart().GetColumn())
		}
	}
	object := ctx.Subquery().Accept(v)
	switch object.(type) {
	case *semantic.SelectStatement:
		stmt := object.(*semantic.SelectStatement)
		stmt.With = with
		return object
	case *semantic.SetOperationStatement:
		return object
	default:
		v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Subquery()),
			ctx.Subquery().GetStart().GetLine(),
			ctx.Subquery().GetStart().GetColumn())
		return nil
	}
}

func (v *plsqlVisitor) VisitSubquery_factoring_clause(ctx *plsql.Subquery_factoring_clauseContext) interface{} {
	clause := newAstNode[semantic.WithClause](ctx)
	for _, item := range ctx.AllFactoring_element() {
		visitor := newExprVisitor(v)
		expr := item.Accept(visitor).(*semantic.CommonTableExpression)
		clause.CTEs = append(clause.CTEs, expr)
	}
	return clause
}

func (v *plsqlVisitor) VisitFor_update_clause(ctx *plsql.For_update_clauseContext) interface{} {
	clause := newAstNode[semantic.ForUpdateClause](ctx)

	if ctx.For_update_of_part() != nil {
		v.ReportError(fmt.Sprintf("unsupported %T", ctx.For_update_of_part()),
			ctx.For_update_of_part().GetStart().GetLine(),
			ctx.For_update_of_part().GetStart().GetColumn())
	}

	if ctx.For_update_options() != nil {
		visitor := newExprVisitor(v)
		var ok bool
		clause.Options, ok = ctx.For_update_options().Accept(visitor).(semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed expression %T", ctx.For_update_options()),
				ctx.For_update_options().GetStart().GetLine(),
				ctx.For_update_options().GetStart().GetColumn())
		}
	}

	return clause
}

func (v *plsqlVisitor) VisitQuery_block(ctx *plsql.Query_blockContext) interface{} {
	stmt := newAstNode[semantic.SelectStatement](ctx)
	stmt.Fields = v.VisitSelected_list(ctx.Selected_list().(*plsql.Selected_listContext)).(*semantic.FieldList)
	stmt.From = ctx.From_clause().Accept(v).(*semantic.FromClause)
	if ctx.Where_clause() != nil {
		if ctx.Where_clause().Expression() != nil {
			visitor := newExprVisitor(v)
			stmt.Where = visitor.VisitExpression(ctx.Where_clause().Expression().(*plsql.ExpressionContext)).(semantic.Expr)
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitSelected_list(ctx *plsql.Selected_listContext) interface{} {
	fields := newAstNode[semantic.FieldList](ctx)
	if ctx.ASTERISK() != nil {
		fields.Fields = append(
			fields.Fields,
			&semantic.SelectField{WildCard: &semantic.WildCardField{Table: "*", Schema: "*"}},
		)
	} else {
		for _, elem := range ctx.AllSelect_list_elements() {
			ev := newExprVisitor(v)
			expt := elem.Accept(ev)
			var filed *semantic.SelectField
			switch expt.(type) {
			case *semantic.SelectField:
				filed = expt.(*semantic.SelectField)
				fields.Fields = append(fields.Fields, filed)
			case semantic.Expr:
				expr := expt.(semantic.Expr)
				filed = &semantic.SelectField{Expr: expr}
				fields.Fields = append(fields.Fields, filed)
			default:
				v.ReportError(fmt.Sprintf("unprocessed syntax %T", elem),
					elem.GetStart().GetLine(), elem.GetStart().GetColumn())
			}
		}
	}

	return fields
}

func (v *plsqlVisitor) VisitTable_ref_list(ctx *plsql.Table_ref_listContext) interface{} {
	from := newAstNode[semantic.FromClause](ctx)
	// tables := make([]*semantic.TableRef, 0, len(ctx.AllTable_ref()))
	for _, t := range ctx.AllTable_ref() {
		from.TableRefs = append(from.TableRefs, &semantic.TableRef{
			Table: t.GetText(),
		})
	}
	return from
}

func (v *plsqlVisitor) VisitCreate_procedure_body(ctx *plsql.Create_procedure_bodyContext) interface{} {
	stmt := newAstNode[semantic.CreateProcedureStatement](ctx)
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

func (v *plsqlVisitor) VisitCreate_function_body(ctx *plsql.Create_function_bodyContext) interface{} {
	stmt := newAstNode[semantic.CreateFunctionStatement](ctx)
	stmt.Name = ctx.Function_name().GetText()
	stmt.IsReplace = ctx.REPLACE() != nil
	for _, p := range ctx.AllParameter() {
		stmt.Parameters = append(stmt.Parameters, v.VisitParameter(p.(*plsql.ParameterContext)).(*semantic.Parameter))
	}
	if ctx.Seq_of_declare_specs() != nil {
		stmt.Declarations = v.VisitSeq_of_declare_specs(ctx.Seq_of_declare_specs().(*plsql.Seq_of_declare_specsContext)).([]semantic.Declaration)
	}
	if ctx.RETURN() != nil {
		stmt.Return = ctx.Type_spec().GetText()
	}
	if ctx.Body() != nil {
		stmt.Body = v.VisitBody(ctx.Body().(*plsql.BodyContext)).(*semantic.Body)
	}
	return stmt
}

func (v *plsqlVisitor) VisitParameter(ctx *plsql.ParameterContext) interface{} {
	param := newAstNode[semantic.Parameter](ctx)
	param.Name = ctx.Parameter_name().GetText()
	param.DataType = ctx.Type_spec().GetText()
	return param
}

func (v *plsqlVisitor) VisitSeq_of_declare_specs(ctx *plsql.Seq_of_declare_specsContext) interface{} {
	decls := make([]semantic.Declaration, 0, len(ctx.AllDeclare_spec()))
	for _, d := range ctx.AllDeclare_spec() {
		node := d.Accept(v)
		decl, ok := node.(semantic.Declaration)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T",
				d.GetChild(0)),
				d.GetStart().GetLine(), d.GetStart().GetColumn())
			continue
		}
		decls = append(decls, decl)
	}
	return decls
}

func (v *plsqlVisitor) VisitVariable_declaration(ctx *plsql.Variable_declarationContext) interface{} {
	varDecl := newAstNode[semantic.VariableDeclaration](ctx)
	varDecl.Name = ctx.Identifier().GetText()
	varDecl.DataType = ctx.Type_spec().GetText()
	if ctx.Default_value_part() != nil {
		visitor := newExprVisitor(v)
		varDecl.Initialization = visitor.VisitExpression(ctx.Default_value_part().Expression().(*plsql.ExpressionContext)).(semantic.Expr)
	}
	return varDecl
}

func (v *plsqlVisitor) VisitException_declaration(ctx *plsql.Exception_declarationContext) interface{} {
	exception := newAstNode[semantic.ExceptionDeclaration](ctx)
	exception.Name = ctx.Identifier().GetText()
	return exception
}

func (v *plsqlVisitor) VisitCursor_declaration(ctx *plsql.Cursor_declarationContext) interface{} {
	cursor := newAstNode[semantic.CursorDeclaration](ctx)
	cursor.Name = ctx.Identifier().GetText()
	for _, p := range ctx.AllParameter_spec() {
		cursor.Parameters = append(cursor.Parameters, v.VisitParameter_spec(p.(*plsql.Parameter_specContext)).(*semantic.Parameter))
	}
	var ok bool
	cursor.Stmt, ok = ctx.Select_statement().Accept(v).(*semantic.SelectStatement)
	if !ok {
		v.ReportError(fmt.Sprintf("unprocessed syntax %T",
			ctx.Select_statement().GetChild(0)),
			ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
	}
	return cursor
}

func (v *plsqlVisitor) VisitParameter_spec(ctx *plsql.Parameter_specContext) interface{} {
	param := newAstNode[semantic.Parameter](ctx)
	param.Name = ctx.Parameter_name().GetText()
	param.DataType = ctx.Type_spec().GetText()
	return param
}

func (v *plsqlVisitor) VisitBody(ctx *plsql.BodyContext) interface{} {
	stmt := newAstNode[semantic.Body](ctx)
	stmt.Statements = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	return stmt
}

func (v *plsqlVisitor) VisitBlock(ctx *plsql.BlockContext) interface{} {
	stmt := newAstNode[semantic.BlockStatement](ctx)
	for _, p := range ctx.AllDeclare_spec() {
		decl, ok := p.Accept(v).(semantic.Declaration)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T",
				p.GetChild(0)),
				p.GetStart().GetLine(), p.GetStart().GetColumn())
			continue
		}
		stmt.Declarations = append(stmt.Declarations, decl)
	}
	stmt.Body = v.VisitBody(ctx.Body().(*plsql.BodyContext)).(*semantic.Body)
	return stmt
}

func (v *plsqlVisitor) VisitAnonymous_block(ctx *plsql.Anonymous_blockContext) interface{} {
	stmt := newAstNode[semantic.BlockStatement](ctx)
	if ctx.Seq_of_declare_specs() != nil {
		for _, p := range ctx.Seq_of_declare_specs().AllDeclare_spec() {
			stmt.Declarations = append(stmt.Declarations, p.Accept(v).(semantic.Declaration))
		}
	}
	stmts := v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)

	stmt.Body = &semantic.Body{Statements: stmts}
	return stmt
}

func (v *plsqlVisitor) VisitSeq_of_statements(ctx *plsql.Seq_of_statementsContext) interface{} {
	stmts := make([]semantic.Statement, 0)
	for _, node := range ctx.GetChildren() {
		stmt, ok := node.(antlr.ParserRuleContext)
		if !ok {
			continue
		}
		s, ok := stmt.Accept(v).(semantic.Statement)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T", stmt.GetChild(0)),
				stmt.GetStart().GetLine(),
				stmt.GetStart().GetColumn())
			return stmts
		}
		stmts = append(stmts, s)
	}
	return stmts
}

func (v *plsqlVisitor) VisitAssignment_statement(ctx *plsql.Assignment_statementContext) interface{} {
	stmt := newAstNode[semantic.AssignmentStatement](ctx)
	if ctx.General_element() != nil {
		stmt.Left = ctx.General_element().GetText()
	} else if ctx.Bind_variable() != nil {
		stmt.Left = ctx.Bind_variable().GetText()
	}
	visitor := newExprVisitor(v)
	stmt.Right = visitor.VisitExpression(ctx.Expression().(*plsql.ExpressionContext)).(semantic.Expr)

	return stmt
}

func (v *plsqlVisitor) VisitIf_statement(ctx *plsql.If_statementContext) interface{} {
	stmt := newAstNode[semantic.IfStatement](ctx)
	if ctx.Condition() != nil {
		vistior := newExprVisitor(v)
		stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
	}
	stmt.ThenBlock = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	if ctx.Else_part() != nil {
		stmt.ElseBlock = v.VisitSeq_of_statements(ctx.Else_part().Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	}
	return stmt
}

func (v *plsqlVisitor) VisitOpen_statement(ctx *plsql.Open_statementContext) interface{} {
	stmt := newAstNode[semantic.OpenStatement](ctx)
	return stmt
}

func (v *plsqlVisitor) VisitOpen_for_statement(ctx *plsql.Open_for_statementContext) interface{} {
	stmt := newAstNode[semantic.OpenForStatement](ctx)
	visitor := newExprVisitor(v)
	expr, ok := ctx.Variable_name().Accept(visitor).(semantic.Expr)
	if !ok {
		v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Variable_name()),
			ctx.GetStart().GetLine(),
			ctx.GetStart().GetColumn())
		return stmt
	}
	stmt.Name = expr

	if ctx.Select_statement() != nil {
		s, ok := ctx.Select_statement().Accept(v).(semantic.Statement)
		if !ok {
			v.ReportError(fmt.Sprintf("unprocessed syntax %T", ctx.Select_statement()),
				ctx.GetStart().GetLine(),
				ctx.GetStart().GetColumn())
			return stmt
		}
		stmtExpr := newAstNode[semantic.StatementExpression](ctx.Select_statement())
		stmtExpr.Stmt = s
		stmt.For = stmtExpr
	}

	return stmt
}

func (v *plsqlVisitor) VisitClose_statement(ctx *plsql.Close_statementContext) interface{} {
	stmt := newAstNode[semantic.CloseStatement](ctx)
	return stmt
}

func (v *plsqlVisitor) VisitFetch_statement(ctx *plsql.Fetch_statementContext) interface{} {
	stmt := newAstNode[semantic.FetchStatement](ctx)
	stmt.Cursor = ctx.Cursor_name().GetText()
	stmt.Into = ctx.General_element(0).GetText()
	return stmt
}

func (v *plsqlVisitor) VisitExit_statement(ctx *plsql.Exit_statementContext) interface{} {
	stmt := newAstNode[semantic.ExitStatement](ctx)
	if ctx.Condition() != nil {
		vistior := newExprVisitor(v)
		stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
	}
	return stmt
}

func (v *plsqlVisitor) VisitLoop_statement(ctx *plsql.Loop_statementContext) interface{} {
	stmt := newAstNode[semantic.LoopStatement](ctx)
	stmt.Statements = v.VisitSeq_of_statements(ctx.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
	return stmt
}

func (v *plsqlVisitor) VisitReturn_statement(ctx *plsql.Return_statementContext) interface{} {
	stmt := newAstNode[semantic.ReturnStatement](ctx)
	if ctx.Expression() != nil {
		visitor := newExprVisitor(v)
		stmt.Name = visitor.VisitExpression(ctx.Expression().(*plsql.ExpressionContext)).(semantic.Expr)
	}
	return stmt
}

func (v *plsqlVisitor) VisitNull_statement(ctx *plsql.Null_statementContext) interface{} {
	stmt := newAstNode[semantic.NullStatement](ctx)
	return stmt
}

func (v *plsqlVisitor) VisitCall_statement(ctx *plsql.Call_statementContext) interface{} {
	stmt := newAstNode[semantic.ProcedureCall](ctx)
	visitor := newExprVisitor(v)
	stmt.Name = visitor.VisitRoutine_name(ctx.Routine_name().(*plsql.Routine_nameContext)).(semantic.Expr)
	if ctx.Function_argument() != nil {
		visitor := newExprVisitor(v)
		stmt.Arguments = ctx.Function_argument().Accept(visitor).([]semantic.Expr)
	}
	return stmt
}

func (v *plsqlVisitor) VisitFunction_argument(ctx *plsql.Function_argumentContext) interface{} {
	exprs := make([]semantic.Expr, 0, len(ctx.AllArgument()))
	for _, c := range ctx.AllArgument() {
		visitor := newExprVisitor(v)
		expr := visitor.VisitExpression(c.Expression().(*plsql.ExpressionContext)).(semantic.Expr)
		exprs = append(exprs, expr)
	}
	return exprs
}

func (v *plsqlVisitor) VisitCreate_package(ctx *plsql.Create_packageContext) interface{} {
	stmt := newAstNode[semantic.CreatePackageStatement](ctx)

	stmt.Name = ctx.Package_name(0).GetText()
	for _, p := range ctx.AllPackage_obj_spec() {
		spec := p.Accept(v)
		switch spec.(type) {
		case *semantic.CreateProcedureStatement:
			stmt.Procedures = append(stmt.Procedures, spec.(*semantic.CreateProcedureStatement))
		case semantic.Declaration:
			stmt.Types = append(stmt.Types, spec.(semantic.Declaration))
		case *semantic.VariableDeclaration:
			stmt.Variables = append(stmt.Variables, spec.(semantic.Declaration))
		default:
			v.ReportError(fmt.Sprintf("unprocessed syntax %T", p.GetChild(0)),
				p.GetStart().GetLine(),
				p.GetStart().GetColumn())
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitCreate_package_body(ctx *plsql.Create_package_bodyContext) interface{} {
	stmt := newAstNode[semantic.CreatePackageBodyStatement](ctx)

	stmt.Name = ctx.Package_name(0).GetText()
	for _, p := range ctx.AllPackage_obj_body() {
		s := p.Accept(v)
		switch s.(type) {
		case *semantic.CreateProcedureStatement:
			stmt.Procedures = append(stmt.Procedures, s.(*semantic.CreateProcedureStatement))
		case *semantic.CreateFunctionStatement:
			stmt.Functions = append(stmt.Functions, s.(*semantic.CreateFunctionStatement))
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitProcedure_spec(ctx *plsql.Procedure_specContext) interface{} {
	stmt := newAstNode[semantic.CreateProcedureStatement](ctx)

	stmt.Name = ctx.Identifier().GetText()
	for _, p := range ctx.AllParameter() {
		stmt.Parameters = append(stmt.Parameters, v.VisitParameter(p.(*plsql.ParameterContext)).(*semantic.Parameter))
	}

	return stmt
}

func (v *plsqlVisitor) VisitProcedure_body(ctx *plsql.Procedure_bodyContext) interface{} {
	stmt := newAstNode[semantic.CreateProcedureStatement](ctx)

	stmt.Name = ctx.Identifier().GetText()
	for _, p := range ctx.AllParameter() {
		stmt.Parameters = append(stmt.Parameters, v.VisitParameter(p.(*plsql.ParameterContext)).(*semantic.Parameter))
	}
	if ctx.Seq_of_declare_specs() != nil {
		stmt.Declarations = v.VisitSeq_of_declare_specs(ctx.Seq_of_declare_specs().(*plsql.Seq_of_declare_specsContext)).([]semantic.Declaration)
	}
	stmt.Body = v.VisitBody(ctx.Body().(*plsql.BodyContext)).(*semantic.Body)
	return stmt
}

func (v *plsqlVisitor) VisitCreate_type(ctx *plsql.Create_typeContext) interface{} {
	if ctx.Type_definition() != nil {
		return v.VisitType_definition(ctx.Type_definition().(*plsql.Type_definitionContext))
	}

	return v.VisitChildren(ctx)
}

func (v *plsqlVisitor) VisitType_definition(ctx *plsql.Type_definitionContext) interface{} {
	value := ctx.Object_type_def().Accept(v)
	switch value.(type) {
	case *semantic.CreateNestTableStatement:
		stmt := value.(*semantic.CreateNestTableStatement)
		stmt.Name = ctx.Type_name().GetText()
		return stmt
	}

	return value
}

func (v *plsqlVisitor) VisitNested_table_type_def(ctx *plsql.Nested_table_type_defContext) interface{} {
	stmt := newAstNode[semantic.CreateNestTableStatement](ctx)

	return stmt
}

func (v *plsqlVisitor) VisitType_declaration(ctx *plsql.Type_declarationContext) interface{} {
	if ctx.Table_type_def() != nil {
		decl := v.VisitTable_type_def(ctx.Table_type_def().(*plsql.Table_type_defContext)).(*semantic.NestTableTypeDeclaration)
		decl.Name = ctx.Identifier().GetText()
		setAstSpan(ctx, decl)
		return decl
	}

	if ctx.Ref_cursor_type_def() != nil {
		decl := v.VisitRef_cursor_type_def(ctx.Ref_cursor_type_def().(*plsql.Ref_cursor_type_defContext)).(*semantic.CursorDeclaration)
		decl.Name = ctx.Identifier().GetText()
		return decl
	}
	return nil
}

func (v *plsqlVisitor) VisitTable_type_def(ctx *plsql.Table_type_defContext) interface{} {
	stmt := newAstNode[semantic.NestTableTypeDeclaration](ctx)
	return stmt
}

func (v *plsqlVisitor) VisitRef_cursor_type_def(ctx *plsql.Ref_cursor_type_defContext) interface{} {
	decl := newAstNode[semantic.CursorDeclaration](ctx)
	decl.IsReference = true
	if ctx.Type_spec() != nil {
		decl.Return = ctx.Type_spec().GetText()
	}
	return decl
}

func (v *plsqlVisitor) VisitFunction_spec(ctx *plsql.Function_specContext) interface{} {
	decl := newAstNode[semantic.FunctionDeclaration](ctx)

	decl.Name = ctx.Identifier().GetText()

	return decl
}

func (v *plsqlVisitor) VisitPragma_declaration(ctx *plsql.Pragma_declarationContext) interface{} {
	if ctx.AUTONOMOUS_TRANSACTION() != nil {
		decl := newAstNode[semantic.AutonomousTransactionDeclaration](ctx)
		return decl
	}

	return nil
}

func (v *plsqlVisitor) VisitFunction_body(ctx *plsql.Function_bodyContext) interface{} {
	stmt := newAstNode[semantic.CreateFunctionStatement](ctx)

	stmt.Name = ctx.Identifier().GetText()
	for _, p := range ctx.AllParameter() {
		stmt.Parameters = append(stmt.Parameters, v.VisitParameter(p.(*plsql.ParameterContext)).(*semantic.Parameter))
	}
	if ctx.Seq_of_declare_specs() != nil {
		stmt.Declarations = v.VisitSeq_of_declare_specs(ctx.Seq_of_declare_specs().(*plsql.Seq_of_declare_specsContext)).([]semantic.Declaration)
	}
	stmt.Body = v.VisitBody(ctx.Body().(*plsql.BodyContext)).(*semantic.Body)
	return stmt
}

func (v *plsqlVisitor) VisitSimple_case_statement(ctx *plsql.Simple_case_statementContext) interface{} {
	stmt := newAstNode[semantic.CaseWhenStatement](ctx)
	if ctx.Expression() != nil {
		visitor := newExprVisitor(v)
		stmt.Expr = visitor.VisitExpression(ctx.Expression().(*plsql.ExpressionContext)).(semantic.Expr)
	}
	for _, item := range ctx.AllSimple_case_when_part() {
		block := newAstNode[semantic.CaseWhenBlock](ctx)
		if item.Expression(0) != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression(0).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression(0).GetStart().GetLine(),
					item.Expression(0).GetStart().GetColumn())
			} else {
				block.Condition = expr
			}
		}
		if item.Seq_of_statements() != nil {
			block.Stmts = v.VisitSeq_of_statements(item.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
		}
		if item.Expression(1) != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression(1).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression(1).GetStart().GetLine(),
					item.Expression(1).GetStart().GetColumn())
			} else {
				block.Expr = expr
			}
		}
		stmt.WhenClauses = append(stmt.WhenClauses, block)
	}
	if ctx.Case_else_part() != nil {
		item := ctx.Case_else_part()
		block := newAstNode[semantic.CaseWhenBlock](ctx)
		if item.Seq_of_statements() != nil {
			block.Stmts = v.VisitSeq_of_statements(item.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
		}
		if item.Expression() != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression().(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression().GetStart().GetLine(),
					item.Expression().GetStart().GetColumn())
			} else {
				block.Expr = expr
			}
		}
		stmt.ElseClause = block
	}
	return stmt
}

func (v *plsqlVisitor) VisitSearched_case_statement(ctx *plsql.Searched_case_statementContext) interface{} {
	stmt := newAstNode[semantic.CaseWhenStatement](ctx)
	for _, item := range ctx.AllSearched_case_when_part() {
		block := newAstNode[semantic.CaseWhenBlock](ctx)
		if item.Expression(0) != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression(0).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression(0).GetStart().GetLine(),
					item.Expression(0).GetStart().GetColumn())
			} else {
				block.Condition = expr
			}
		}
		if item.Seq_of_statements() != nil {
			block.Stmts = v.VisitSeq_of_statements(item.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
		}
		if item.Expression(1) != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression(1).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression(1).GetStart().GetLine(),
					item.Expression(1).GetStart().GetColumn())
			} else {
				block.Expr = expr
			}
		}
		stmt.WhenClauses = append(stmt.WhenClauses, block)
	}
	if ctx.Case_else_part() != nil {
		item := ctx.Case_else_part()
		block := newAstNode[semantic.CaseWhenBlock](ctx)
		if item.Seq_of_statements() != nil {
			block.Stmts = v.VisitSeq_of_statements(item.Seq_of_statements().(*plsql.Seq_of_statementsContext)).([]semantic.Statement)
		}
		if item.Expression() != nil {
			visitor := newExprVisitor(v)
			expr, ok := visitor.VisitExpression(item.Expression().(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					item.Expression().GetStart().GetLine(),
					item.Expression().GetStart().GetColumn())
			} else {
				block.Expr = expr
			}
		}
		stmt.ElseClause = block
	}
	return stmt
}

func (v *plsqlVisitor) VisitCommit_statement(ctx *plsql.Commit_statementContext) interface{} {
	stmt := newAstNode[semantic.CommitStatement](ctx)
	if ctx.COMMENT() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.COMMENT()),
			ctx.COMMENT().GetSymbol().GetLine(),
			ctx.COMMENT().GetSymbol().GetColumn())
	}
	if ctx.FORCE() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.FORCE()),
			ctx.FORCE().GetSymbol().GetLine(),
			ctx.FORCE().GetSymbol().GetColumn())
	}

	return stmt
}

func (v *plsqlVisitor) VisitRollback_statement(ctx *plsql.Rollback_statementContext) interface{} {
	stmt := newAstNode[semantic.RollbackStatement](ctx)
	if ctx.TO() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.TO()),
			ctx.TO().GetSymbol().GetLine(),
			ctx.TO().GetSymbol().GetColumn())
	}
	if ctx.FORCE() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.FORCE()),
			ctx.FORCE().GetSymbol().GetLine(),
			ctx.FORCE().GetSymbol().GetColumn())
	}
	return stmt
}

func (v *plsqlVisitor) VisitDelete_statement(ctx *plsql.Delete_statementContext) interface{} {
	stmt := newAstNode[semantic.DeleteStatement](ctx)
	if ctx.General_table_ref() != nil {
		visitor := newExprVisitor(v)
		object := ctx.General_table_ref().Accept(visitor)
		expr, ok := object.(semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.General_table_ref().GetChild(0)),
				ctx.General_table_ref().GetStart().GetLine(),
				ctx.General_table_ref().GetStart().GetColumn())
		}
		stmt.Table = expr
	}

	if ctx.Where_clause() != nil {
		if ctx.Where_clause().Expression() != nil {
			visitor := newExprVisitor(v)
			stmt.Where = visitor.VisitExpression(ctx.Where_clause().Expression().(*plsql.ExpressionContext)).(semantic.Expr)
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitUpdate_statement(ctx *plsql.Update_statementContext) interface{} {
	stmt := newAstNode[semantic.UpdateStatement](ctx)
	if ctx.General_table_ref() != nil {
		visitor := newExprVisitor(v)
		object := ctx.General_table_ref().Accept(visitor)
		expr, ok := object.(semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.General_table_ref().GetChild(0)),
				ctx.General_table_ref().GetStart().GetLine(),
				ctx.General_table_ref().GetStart().GetColumn())
		}
		stmt.Table = expr
	}

	if ctx.Update_set_clause() != nil {
		if ctx.Update_set_clause().VALUE() != nil {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Update_set_clause().VALUE()),
				ctx.Update_set_clause().VALUE().GetSymbol().GetLine(),
				ctx.Update_set_clause().VALUE().GetSymbol().GetColumn())
		} else {
			exprs := make([]semantic.Expr, 0)
			for _, item := range ctx.Update_set_clause().AllColumn_based_update_set_clause() {
				visitor := newExprVisitor(v)
				object := item.Accept(visitor)
				expr, ok := object.(semantic.Expr)
				if !ok {
					v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Update_set_clause().GetChild(1)),
						ctx.Update_set_clause().GetStart().GetLine(),
						ctx.Update_set_clause().GetStart().GetColumn())
				}
				exprs = append(exprs, expr)
			}
			stmt.SetExprs = exprs
		}
	}

	if ctx.Where_clause() != nil {
		if ctx.Where_clause().Expression() != nil {
			visitor := newExprVisitor(v)
			stmt.Where = visitor.VisitExpression(ctx.Where_clause().Expression().(*plsql.ExpressionContext)).(semantic.Expr)
		}
	}

	if ctx.Static_returning_clause() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Static_returning_clause()),
			ctx.Static_returning_clause().GetStart().GetLine(),
			ctx.Static_returning_clause().GetStart().GetColumn())
	}

	if ctx.Error_logging_clause() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Error_logging_clause()),
			ctx.Error_logging_clause().GetStart().GetLine(),
			ctx.Error_logging_clause().GetStart().GetColumn())
	}
	return stmt
}

func (v *plsqlVisitor) VisitContinue_statement(ctx *plsql.Continue_statementContext) interface{} {
	stmt := newAstNode[semantic.ContinueStatement](ctx)
	if ctx.Label_name() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Label_name()),
			ctx.Label_name().GetStart().GetLine(),
			ctx.Label_name().GetStart().GetColumn())
	}
	if ctx.Condition() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Condition()),
			ctx.Condition().GetStart().GetLine(),
			ctx.Condition().GetStart().GetColumn())
	}
	return stmt
}

func (v *plsqlVisitor) VisitExecute_immediate(ctx *plsql.Execute_immediateContext) interface{} {
	stmt := newAstNode[semantic.ExecuteImmediateStatement](ctx)
	stmt.Sql = ctx.Expression().GetText()

	if ctx.Into_clause() != nil {
		into, ok := ctx.Into_clause().Accept(v).(*semantic.IntoClause)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Into_clause()),
				ctx.Into_clause().GetStart().GetLine(),
				ctx.Into_clause().GetStart().GetColumn())
		}
		stmt.Into = into
	}
	if ctx.Using_clause() != nil {
		visitor := newExprVisitor(v)
		using, ok := ctx.Using_clause().Accept(visitor).(*semantic.UsingClause)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Using_clause()),
				ctx.Using_clause().GetStart().GetLine(),
				ctx.Using_clause().GetStart().GetColumn())
		}
		stmt.Using = using
	}
	if ctx.Dynamic_returning_clause() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Dynamic_returning_clause()),
			ctx.Dynamic_returning_clause().GetStart().GetLine(),
			ctx.Dynamic_returning_clause().GetStart().GetColumn())
	}
	return stmt
}

func (v *plsqlVisitor) VisitInto_clause(ctx *plsql.Into_clauseContext) interface{} {
	stmt := newAstNode[semantic.IntoClause](ctx)
	if ctx.BULK() != nil {
		stmt.IsBulk = true
	}

	for _, child := range ctx.GetChildren() {
		switch child.(type) {
		case antlr.TerminalNode:
			continue
		default:
			visitor := newExprVisitor(v)
			object, ok := child.(antlr.ParseTree).Accept(visitor).(semantic.Expr)
			if !ok {
				v.ReportError(fmt.Sprintf("unsupported syntax %T", child),
					child.(antlr.ParserRuleContext).GetStart().GetLine(),
					child.(antlr.ParserRuleContext).GetStart().GetColumn())
			}
			stmt.Vars = append(stmt.Vars, object)
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitCreate_synonym(ctx *plsql.Create_synonymContext) interface{} {
	stmt := newAstNode[semantic.CreateSynonymStatement](ctx)

	visitor := newExprVisitor(v)
	object := ctx.Synonym_name().Accept(visitor)
	expr, ok := object.(semantic.Expr)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Synonym_name().GetChild(0)),
			ctx.Synonym_name().GetStart().GetLine(),
			ctx.Synonym_name().GetStart().GetColumn())
	}
	stmt.Synonym = expr
	object = ctx.Schema_object_name().Accept(visitor)
	expr, ok = object.(semantic.Expr)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Schema_object_name().GetChild(0)),
			ctx.Schema_object_name().GetStart().GetLine(),
			ctx.Schema_object_name().GetStart().GetColumn())
	}
	stmt.Original = expr
	return stmt
}

func (v *plsqlVisitor) VisitInsert_statement(ctx *plsql.Insert_statementContext) interface{} {
	if ctx.Multi_table_insert() != nil {
		stmt, ok := ctx.Multi_table_insert().Accept(v).(*semantic.InsertStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Multi_table_insert()),
				ctx.Multi_table_insert().GetStart().GetLine(),
				ctx.Multi_table_insert().GetStart().GetColumn())
		}
		return stmt
	} else { // single table insert
		stmt, ok := ctx.Single_table_insert().Accept(v).(*semantic.InsertStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Single_table_insert().GetChild(0)),
				ctx.Single_table_insert().GetStart().GetLine(),
				ctx.Single_table_insert().GetStart().GetColumn())
		}
		return stmt
	}
}

func (v *plsqlVisitor) VisitSingle_table_insert(ctx *plsql.Single_table_insertContext) interface{} {
	stmt := newAstNode[semantic.InsertStatement](ctx)

	clause, ok := ctx.Insert_into_clause().Accept(v).(*semantic.InsertIntoClause)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Insert_into_clause().GetChild(0)),
			ctx.Insert_into_clause().GetStart().GetLine(),
			ctx.Insert_into_clause().GetStart().GetColumn())
	} else {
		stmt.AllInto = append(stmt.AllInto, clause)
	}

	if ctx.Values_clause() != nil {
		values, ok := ctx.Values_clause().Accept(v).([]semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Values_clause().GetChild(0)),
				ctx.Values_clause().GetStart().GetLine(),
				ctx.Values_clause().GetStart().GetColumn())
		} else {
			clause.Values = values
		}
	} else if ctx.Select_statement() != nil {
		var ok bool
		stmt.Select, ok = ctx.Select_statement().Accept(v).(*semantic.SelectStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Select_statement()),
				ctx.Select_statement().GetStart().GetLine(),
				ctx.Select_statement().GetStart().GetColumn())
		}
	}

	return stmt
}

func (v *plsqlVisitor) VisitInsert_into_clause(ctx *plsql.Insert_into_clauseContext) interface{} {
	stmt := newAstNode[semantic.InsertIntoClause](ctx)

	stmt.Table = &semantic.TableRef{
		Table: ctx.General_table_ref().GetText(),
	}
	if ctx.Paren_column_list() != nil {
		visitor := newExprVisitor(v)
		objects, ok := ctx.Paren_column_list().Accept(visitor).([]semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Paren_column_list().GetChild(0)),
				ctx.Paren_column_list().GetStart().GetLine(),
				ctx.Paren_column_list().GetStart().GetColumn())
		}
		stmt.Columns = objects
	}
	return stmt
}

func (v *plsqlVisitor) VisitValues_clause(ctx *plsql.Values_clauseContext) interface{} {
	if ctx.REGULAR_ID() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.REGULAR_ID()),
			ctx.REGULAR_ID().GetSymbol().GetLine(),
			ctx.REGULAR_ID().GetSymbol().GetColumn())
		return nil
	} else if ctx.Expressions() != nil {
		visitor := newExprVisitor(v)
		exprs, ok := ctx.Expressions().Accept(visitor).([]semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Expressions().GetChild(0)),
				ctx.Expressions().GetStart().GetLine(),
				ctx.Expressions().GetStart().GetColumn())
		}
		return exprs
	}
	return nil
}

func (v *plsqlVisitor) VisitMulti_table_insert(ctx *plsql.Multi_table_insertContext) interface{} {
	stmt := newAstNode[semantic.InsertStatement](ctx)
	if ctx.ALL() != nil {
		for _, item := range ctx.AllMulti_table_element() {
			into, ok := item.Accept(v).(*semantic.InsertIntoClause)
			if !ok {
				v.ReportError(fmt.Sprintf("unsupported syntax %T", item),
					item.GetStart().GetLine(),
					item.GetStart().GetColumn())
			} else {
				stmt.AllInto = append(stmt.AllInto, into)
			}
		}
	}
	if ctx.Conditional_insert_clause() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Conditional_insert_clause()),
			ctx.Conditional_insert_clause().GetStart().GetLine(),
			ctx.Conditional_insert_clause().GetStart().GetColumn())
	}
	selStmt, ok := ctx.Select_statement().Accept(v).(*semantic.SelectStatement)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Select_statement()),
			ctx.Select_statement().GetStart().GetLine(),
			ctx.Select_statement().GetStart().GetColumn())
	} else {
		stmt.Select = selStmt
	}
	return stmt
}

func (v *plsqlVisitor) VisitMulti_table_element(ctx *plsql.Multi_table_elementContext) interface{} {
	into, ok := ctx.Insert_into_clause().Accept(v).(*semantic.InsertIntoClause)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Insert_into_clause()),
			ctx.Insert_into_clause().GetStart().GetLine(),
			ctx.Insert_into_clause().GetStart().GetColumn())
	}
	if ctx.Values_clause() != nil {
		values, ok := ctx.Values_clause().Accept(v).([]semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Values_clause()),
				ctx.Values_clause().GetStart().GetLine(),
				ctx.Values_clause().GetStart().GetColumn())
		} else {
			into.Values = values
		}
	}
	return into
}

func (v *plsqlVisitor) VisitSubquery(ctx *plsql.SubqueryContext) interface{} {
	selStmt, ok := ctx.Subquery_basic_elements().Accept(v).(*semantic.SelectStatement)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Subquery_basic_elements().GetChild(0)),
			ctx.Subquery_basic_elements().GetStart().GetLine(),
			ctx.Subquery_basic_elements().GetStart().GetColumn())
		return nil
	}
	if len(ctx.AllSubquery_operation_part()) == 0 {
		return selStmt
	}

	stmt := newAstNode[semantic.SetOperationStatement](ctx)
	stmt.SelectList = append(stmt.SelectList, selStmt)

	for _, child := range ctx.AllSubquery_operation_part() {
		selStmt, ok = child.Accept(v).(*semantic.SelectStatement)
		stmt.SelectList = append(stmt.SelectList, selStmt)
	}

	return stmt
}

func (v *plsqlVisitor) VisitSubquery_operation_part(ctx *plsql.Subquery_operation_partContext) interface{} {
	stmt, ok := ctx.Subquery_basic_elements().Accept(v).(*semantic.SelectStatement)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Subquery_basic_elements().GetChild(0)),
			ctx.Subquery_basic_elements().GetStart().GetLine(),
			ctx.Subquery_basic_elements().GetStart().GetColumn())
		return nil
	}

	if ctx.UNION() != nil {
		var op semantic.SetOperator = semantic.Union
		if ctx.ALL() != nil {
			op = semantic.UnionAll
		}
		stmt.SetOperator = &op
	}
	if ctx.INTERSECT() != nil {
		var op semantic.SetOperator = semantic.Intersect
		stmt.SetOperator = &op
	}
	if ctx.MINUS() != nil {
		var op semantic.SetOperator = semantic.Minus
		stmt.SetOperator = &op
	}
	return stmt
}

func (v *plsqlVisitor) VisitDrop_function(ctx *plsql.Drop_functionContext) interface{} {
	stmt := newAstNode[semantic.DropFunctionStatement](ctx)
	stmt.Name = ctx.Function_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitDrop_procedure(ctx *plsql.Drop_procedureContext) interface{} {
	stmt := newAstNode[semantic.DropProcedureStatement](ctx)
	stmt.Name = ctx.Procedure_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitDrop_package(ctx *plsql.Drop_packageContext) interface{} {
	stmt := newAstNode[semantic.DropPackageStatement](ctx)
	stmt.Name = ctx.Package_name().GetText()
	if ctx.BODY() != nil {
		stmt.IsBody = true
	}
	if ctx.Schema_object_name() != nil {
		stmt.Schema = ctx.Schema_object_name().GetText()
	}
	return stmt
}

func (v *plsqlVisitor) VisitDrop_trigger(ctx *plsql.Drop_triggerContext) interface{} {
	stmt := newAstNode[semantic.DropTriggerStatement](ctx)
	stmt.Name = ctx.Trigger_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitRaise_statement(ctx *plsql.Raise_statementContext) interface{} {
	stmt := newAstNode[semantic.RaiseStatement](ctx)
	stmt.Name = ctx.Exception_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitGoto_statement(ctx *plsql.Goto_statementContext) interface{} {
	stmt := newAstNode[semantic.GotoStatement](ctx)
	stmt.Label = ctx.Label_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitLabel_declaration(ctx *plsql.Label_declarationContext) interface{} {
	stmt := newAstNode[semantic.LabelDeclaration](ctx)
	stmt.Label = ctx.Label_name().GetText()
	return stmt
}

func (v *plsqlVisitor) VisitMerge_statement(ctx *plsql.Merge_statementContext) interface{} {
	stmt := newAstNode[semantic.MergeStatement](ctx)
	stmt.Table = &semantic.TableRef{Table: ctx.Tableview_name().GetText()}

	// Using
	visitor := newExprVisitor(v)
	expr, ok := ctx.Selected_tableview().Accept(visitor).(semantic.Expr)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Selected_tableview().GetChild(0)),
			ctx.Selected_tableview().GetStart().GetLine(),
			ctx.Selected_tableview().GetStart().GetColumn())
	}
	stmt.Using = expr

	// On condition
	expr, ok = ctx.Condition().Accept(visitor).(semantic.Expr)
	if !ok {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Condition().GetChild(0)),
			ctx.Condition().GetStart().GetLine(),
			ctx.Condition().GetStart().GetColumn())
	}
	stmt.OnCondition = expr

	// Merge matched
	if ctx.Merge_update_clause() != nil {
		ok := false
		stmt.MergeUpdate, ok = ctx.Merge_update_clause().Accept(v).(*semantic.MergeUpdateStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Merge_update_clause()),
				ctx.Merge_update_clause().GetStart().GetLine(),
				ctx.Merge_update_clause().GetStart().GetColumn())
		}
	}

	// Merge not matched
	if ctx.Merge_insert_clause() != nil {
		ok := false
		stmt.MergeInsert, ok = ctx.Merge_insert_clause().Accept(v).(*semantic.MergeInsertStatement)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Merge_insert_clause()),
				ctx.Merge_insert_clause().GetStart().GetLine(),
				ctx.Merge_insert_clause().GetStart().GetColumn())
		}
	}
	return stmt
}

func (v *plsqlVisitor) VisitMerge_update_clause(ctx *plsql.Merge_update_clauseContext) interface{} {
	stmt := newAstNode[semantic.MergeUpdateStatement](ctx)
	for _, ele := range ctx.AllMerge_element() {
		visitor := newExprVisitor(v)
		expr, ok := ele.Accept(visitor).(semantic.Expr)
		if !ok {
			v.ReportError(fmt.Sprintf("unsupported syntax %T", ele),
				ele.GetStart().GetLine(),
				ele.GetStart().GetColumn())
		} else {
			stmt.SetElems = append(stmt.SetElems, expr)
		}
	}
	// merge_update_delete_part
	if ctx.Merge_update_delete_part() != nil {
		v.ReportError(fmt.Sprintf("unsupported syntax %T", ctx.Merge_update_delete_part()),
			ctx.Merge_update_delete_part().GetStart().GetLine(),
			ctx.Merge_update_delete_part().GetStart().GetColumn())
		return stmt
	}
	return stmt
}

func (v *plsqlVisitor) VisitMerge_insert_clause(ctx *plsql.Merge_insert_clauseContext) interface{} {
	stmt := newAstNode[semantic.MergeInsertStatement](ctx)
	return stmt
}
