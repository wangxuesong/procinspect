package parser

import (
	"fmt"
	"strconv"
	"sync/atomic"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type (
	sqlListener struct {
		*plsql.BasePlSqlParserListener

		script    *semantic.Script
		nodeStack semantic.Stack[semantic.Node]
	}
)

var (
	stmtDepth int64 = 0
)

func Parse(text string) error {
	p := plsql.NewParser(text)
	_ = p.Sql_script()
	return p.Error()
}

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

func peekNodeDepth[T semantic.Node](l *sqlListener, depth int64) (t T, err error) {
	for i := len(l.nodeStack) - 1; i >= 0; i-- {
		node := l.nodeStack[i]
		if _, ok := node.(T); ok {
			if d, ok := node.(semantic.StatementDepth); ok {
				if d.Get() == depth {
					return node.(T), nil
				}
			}
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

func (l *sqlListener) ExitSelected_list(ctx *plsql.Selected_listContext) {
	//stmt := l.nodeStack.Top().(*semantic.SelectStatement)
	stmt, err := peekNode[*semantic.SelectStatement](l)
	if err != nil {
		panic(err)
	}

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
	} else {
		for _, _ = range ctx.AllSelect_list_elements() {
			node := l.nodeStack.Top()
			if _, ok := node.(semantic.Expr); ok {
				selectField := &semantic.SelectField{}
				stmt.Fields.Fields = append(stmt.Fields.Fields, selectField)
				l.nodeStack.Pop()
			}
		}
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
	// set right
	node := l.nodeStack.Top()
	if _, ok := node.(semantic.Expr); ok {
		stmt.Right = node.(semantic.Expr)
		l.nodeStack.Pop()
	} else {
		stmt.Right = nil
	}
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

func (l *sqlListener) ExitCursor_declaration(ctx *plsql.Cursor_declarationContext) {
	decl := &semantic.CursorDeclaration{}
	// set line & column
	decl.SetLine(ctx.GetStart().GetLine())
	decl.SetColumn(ctx.GetStart().GetColumn())
	decl.Name = ctx.Identifier().GetText()
	if ctx.Select_statement() != nil {
		stmt := l.nodeStack.Pop().(*semantic.SelectStatement)
		decl.Stmt = stmt
	}
	for range ctx.AllParameter_spec() {
		decl.Parameters = append(decl.Parameters, l.nodeStack.Pop().(*semantic.Parameter))
	}
	l.nodeStack.Push(decl)
}

func (l *sqlListener) ExitParameter_spec(ctx *plsql.Parameter_specContext) {
	para := &semantic.Parameter{}
	para.Name = ctx.Parameter_name().GetText()
	para.DataType = ctx.Type_spec().GetText()
	l.nodeStack.Push(para)
}

func (l *sqlListener) ExitParameter(ctx *plsql.ParameterContext) {
	stmt := &semantic.Parameter{}
	stmt.Name = ctx.Parameter_name().GetText()
	stmt.DataType = ctx.Type_spec().GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) EnterIf_statement(ctx *plsql.If_statementContext) {
	stmt := &semantic.IfStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Set(atomic.LoadInt64(&stmtDepth))
	atomic.AddInt64(&stmtDepth, 1)
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitIf_statement(ctx *plsql.If_statementContext) {
	stmt, err := peekNodeDepth[*semantic.IfStatement](l, stmtDepth-1)
	if err != nil {
		panic(err)
	}
	for {
		node := l.nodeStack.Top()
		if s, ok := node.(*semantic.IfStatement); ok {
			if s == stmt {
				//stmt.Condition = ctx.Condition().GetText()
				if ctx.Condition() != nil {
					vistior := experVisitor{}
					stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
				}

				break
			}
		}
		switch node.(type) {
		case semantic.Statement:
			stmt.ThenBlock = append([]semantic.Statement{node.(semantic.Statement)}, stmt.ThenBlock...)
		case *semantic.ElseBlock:
			stmt.ElseBlock = node.(*semantic.ElseBlock).Statements
		}

		l.nodeStack.Pop()
	}
	atomic.AddInt64(&stmtDepth, -1)
}

func (l *sqlListener) EnterElse_part(ctx *plsql.Else_partContext) {
	stmt := &semantic.ElseBlock{}
	stmt.Set(atomic.LoadInt64(&stmtDepth))
	atomic.AddInt64(&stmtDepth, 1)
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitElse_part(ctx *plsql.Else_partContext) {
	stmt, err := peekNodeDepth[*semantic.ElseBlock](l, stmtDepth-1)
	if err != nil {
		panic(err)
	}
	for {
		node := l.nodeStack.Top()
		if s, ok := node.(*semantic.ElseBlock); ok {
			if s == stmt {
				break
			}
		}
		switch node.(type) {
		case semantic.Statement:
			stmt.Statements = append([]semantic.Statement{node.(semantic.Statement)}, stmt.Statements...)
		}
		l.nodeStack.Pop()
	}
	atomic.AddInt64(&stmtDepth, -1)
}

func (l *sqlListener) ExitOpen_statement(ctx *plsql.Open_statementContext) {
	stmt := &semantic.OpenStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Cursor_name().GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitClose_statement(ctx *plsql.Close_statementContext) {
	stmt := &semantic.CloseStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Cursor_name().GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) EnterLoop_statement(ctx *plsql.Loop_statementContext) {
	stmt := &semantic.LoopStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Set(atomic.LoadInt64(&stmtDepth))
	atomic.AddInt64(&stmtDepth, 1)
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitLoop_statement(ctx *plsql.Loop_statementContext) {
	stmt, err := peekNodeDepth[*semantic.LoopStatement](l, stmtDepth-1)
	if err != nil {
		panic(err)
	}
	for {
		node := l.nodeStack.Top()
		if s, ok := node.(*semantic.LoopStatement); ok {
			if s == stmt {
				break
			}
		}
		switch node.(type) {
		case semantic.Statement:
			stmt.Statements = append([]semantic.Statement{node.(semantic.Statement)}, stmt.Statements...)
		}
		l.nodeStack.Pop()
	}
	atomic.AddInt64(&stmtDepth, -1)
}

func (l *sqlListener) ExitFetch_statement(ctx *plsql.Fetch_statementContext) {
	stmt := &semantic.FetchStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Cursor = ctx.Cursor_name().GetText()
	stmt.Into = ctx.Variable_name(0).GetText()
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitExit_statement(ctx *plsql.Exit_statementContext) {
	stmt := &semantic.ExitStatement{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	if ctx.Condition() != nil {
		vistior := experVisitor{}
		stmt.Condition = vistior.VisitCondition(ctx.Condition().(*plsql.ConditionContext)).(semantic.Expr)
	}
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitFunction_call(ctx *plsql.Function_callContext) {
	stmt := &semantic.ProcedureCall{}
	stmt.SetLine(ctx.GetStart().GetLine())
	stmt.SetColumn(ctx.GetStart().GetColumn())
	stmt.Name = ctx.Routine_name().GetText()
	for range ctx.Function_argument().AllArgument() {
		node := l.nodeStack.Top()
		switch node.(type) {
		case *semantic.Argument:
			stmt.Arguments = append([]*semantic.Argument{node.(*semantic.Argument)}, stmt.Arguments...)
			l.nodeStack.Pop()
		}
	}
	l.nodeStack.Push(stmt)
}

func (l *sqlListener) ExitFunction_argument(ctx *plsql.Function_argumentContext) {
	for _, arg := range ctx.AllArgument() {
		stmt := &semantic.Argument{}
		stmt.SetLine(arg.GetStart().GetLine())
		stmt.SetColumn(arg.GetStart().GetColumn())

		stmt.Name = arg.GetText()
		l.nodeStack.Push(stmt)
	}
}

func (l *sqlListener) ExitNumeric(ctx *plsql.NumericContext) {
	number := &semantic.NumericLiteral{}
	number.SetLine(ctx.GetStart().GetLine())
	number.SetColumn(ctx.GetStart().GetColumn())
	if ctx.UNSIGNED_INTEGER() != nil {
		if v, err := strconv.ParseInt(ctx.GetText(), 10, 64); err == nil {
			number.Value = v
		}
	}
	l.nodeStack.Push(number)
}
