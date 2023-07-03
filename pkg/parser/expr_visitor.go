package parser

import (
	"procinspect/pkg/log"
	"strconv"
	"strings"

	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"

	"github.com/antlr4-go/antlr/v4"
)

type (
	exprVisitor struct {
		plsql.BasePlSqlParserVisitor
		stmtVisitor *plsqlVisitor
	}
)

func newExprVisitor(visitor *plsqlVisitor) *exprVisitor {
	v := &exprVisitor{stmtVisitor: visitor}
	v.BasePlSqlParserVisitor.ParseTreeVisitor = v
	return v
}

func (v *exprVisitor) ReportError(msg string, line, column int) {
	defer log.Sync()
	log.Warn(msg, log.Int("line", line), log.Int("column", column))
}

func (v *exprVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *exprVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	return node.GetText()
}

func (v *exprVisitor) VisitErrorNode(_ antlr.ErrorNode) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *exprVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	children := node.GetChildren()
	nodes := make([]interface{}, 0, len(children))
	for _, child := range children {
		switch child.(type) {
		case *plsql.Other_functionContext:
			c := child.(*plsql.Other_functionContext)
			nodes = append(nodes, v.VisitOther_function(c))
		case *plsql.Logical_expressionContext:
			c := child.(*plsql.Logical_expressionContext)
			nodes = append(nodes, v.VisitLogical_expression(c))
		case *plsql.Unary_logical_expressionContext:
			c := child.(*plsql.Unary_logical_expressionContext)
			nodes = append(nodes, v.VisitUnary_logical_expression(c))
		case *plsql.Relational_expressionContext:
			c := child.(*plsql.Relational_expressionContext)
			nodes = append(nodes, v.VisitRelational_expression(c))
		case *plsql.Compound_expressionContext:
			c := child.(*plsql.Compound_expressionContext)
			nodes = append(nodes, v.VisitCompound_expression(c))
		case *plsql.In_elementsContext:
			c := child.(*plsql.In_elementsContext)
			nodes = append(nodes, v.VisitIn_elements(c))
		case *plsql.ConcatenationContext:
			c := child.(*plsql.ConcatenationContext)
			nodes = append(nodes, v.VisitConcatenation(c))
		case *plsql.Unary_expressionContext:
			c := child.(*plsql.Unary_expressionContext)
			nodes = append(nodes, v.VisitUnary_expression(c))
		case *plsql.Quantified_expressionContext:
			c := child.(*plsql.Quantified_expressionContext)
			nodes = append(nodes, v.VisitQuantified_expression(c))
		case *plsql.AtomContext:
			c := child.(*plsql.AtomContext)
			nodes = append(nodes, v.VisitAtom(c))
		case *plsql.Table_elementContext:
			c := child.(*plsql.Table_elementContext)
			nodes = append(nodes, v.VisitTable_element(c))
		case *plsql.String_functionContext:
			c := child.(*plsql.String_functionContext)
			nodes = append(nodes, v.VisitString_function(c))
		case *plsql.Numeric_functionContext:
			c := child.(*plsql.Numeric_functionContext)
			nodes = append(nodes, v.VisitNumeric_function(c))
		case *plsql.General_element_partContext:
			c := child.(*plsql.General_element_partContext)
			nodes = append(nodes, v.VisitGeneral_element_part(c))
		case *plsql.ConstantContext:
			c := child.(*plsql.ConstantContext)
			nodes = append(nodes, v.VisitConstant(c))
		case *plsql.Quoted_stringContext:
			c := child.(*plsql.Quoted_stringContext)
			nodes = append(nodes, v.VisitQuoted_string(c))
		case *plsql.Variable_nameContext:
			c := child.(*plsql.Variable_nameContext)
			nodes = append(nodes, v.VisitVariable_name(c))
		case *plsql.NumericContext:
			c := child.(*plsql.NumericContext)
			nodes = append(nodes, v.VisitNumeric(c))
		case *plsql.Routine_nameContext:
			c := child.(*plsql.Routine_nameContext)
			nodes = append(nodes, v.VisitRoutine_name(c))
		default:
			tree := child.(antlr.ParseTree)
			c := tree.Accept(v)
			switch c.(type) {
			case []interface{}:
				cc := c.([]interface{})
				nodes = append(nodes, cc...)
			default:
				nodes = append(nodes, c)
			}
		}
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return nodes
}

func (v *exprVisitor) VisitCondition(ctx *plsql.ConditionContext) interface{} {
	if ctx.Expression() != nil {
		expr := v.VisitExpression(ctx.Expression().(*plsql.ExpressionContext))
		return expr
	} else {
		return nil
	}
}

func (v *exprVisitor) VisitExpressions(ctx *plsql.ExpressionsContext) interface{} {
	exprs := make([]semantic.Expr, 0, len(ctx.AllExpression()))
	for _, e := range ctx.AllExpression() {
		var ok bool
		expr, ok := v.VisitExpression(e.(*plsql.ExpressionContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				e.GetStart().GetLine(),
				e.GetStart().GetColumn())
			return expr
		}

		exprs = append(exprs, expr)
	}
	//if len(exprs) == 1 {
	//	return exprs[0]
	//}
	return exprs
}

func (v *exprVisitor) VisitOther_function(ctx *plsql.Other_functionContext) interface{} {
	if ctx.Cursor_name() != nil {
		node := ctx.Cursor_name().(*plsql.Cursor_nameContext)
		ca := &semantic.CursorAttribute{Cursor: node.GetText()}
		ca.SetLine(node.GetStart().GetLine())
		ca.SetColumn(node.GetStart().GetColumn())
		if ctx.PERCENT_NOTFOUND() != nil {
			ca.Attr = "NOTFOUND"
		} else if ctx.PERCENT_FOUND() != nil {
			ca.Attr = "FOUND"
		} else if ctx.PERCENT_ROWCOUNT() != nil {
			ca.Attr = "ROWCOUNT"
		} else if ctx.PERCENT_ISOPEN() != nil {
			ca.Attr = "ISOPEN"
		}
		return ca
	}

	if ctx.TO_NUMBER() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "TO_NUMBER"}}
		arg, ok := v.VisitConcatenation(ctx.Concatenation(0).(*plsql.ConcatenationContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
		}
		expr.Args = append(expr.Args, arg)
		return expr
	}
	if ctx.CAST() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "CAST"}}
		if ctx.Concatenation(0) != nil {
			cast := &semantic.CastExpression{}
			cast.SetLine(ctx.Concatenation(0).GetStart().GetLine())
			cast.SetColumn(ctx.Concatenation(0).GetStart().GetColumn())
			arg, ok := v.VisitConcatenation(ctx.Concatenation(0).(*plsql.ConcatenationContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
				return expr
			}
			cast.Expr = arg
			cast.DataType = ctx.Type_spec().GetText()
			expr.Args = append(expr.Args, cast)
		}
		return expr
	}

	if ctx.Over_clause_keyword() != nil {
		name := ctx.Over_clause_keyword().GetText()
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: name}}
		var arg semantic.Expr
		var ok bool
		for _, child := range ctx.Function_argument_analytic().GetChildren() {
			switch child.(type) {
			case antlr.TerminalNode:
				if child.(antlr.TerminalNode) == ctx.Function_argument_analytic().LEFT_PAREN() {
					continue
				}
				if child.(antlr.TerminalNode) == ctx.Function_argument_analytic().RIGHT_PAREN() {
					expr.Args = append(expr.Args, arg)
					arg = nil
				}
				if child.(antlr.TerminalNode).GetSymbol().GetTokenType() == plsql.PlSqlParserCOMMA {
					expr.Args = append(expr.Args, arg)
					arg = nil
				}
			case *plsql.ArgumentContext:
				context := child.(*plsql.ArgumentContext)
				arg, ok = v.VisitArgument(context).(semantic.Expr)
				if !ok {
					v.ReportError("unsupported expression",
						context.GetStart().GetLine(), context.GetStart().GetColumn())
				}
			case *plsql.Respect_or_ignore_nullsContext:
				context := child.(*plsql.Respect_or_ignore_nullsContext)
				v.ReportError("unsupported expression",
					context.GetStart().GetLine(), context.GetStart().GetColumn())
			case *plsql.Keep_clauseContext:
				context := child.(*plsql.Keep_clauseContext)
				v.ReportError("unsupported expression",
					context.GetStart().GetLine(), context.GetStart().GetColumn())
			}
		}
		return expr
	}

	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitLogical_expression(ctx *plsql.Logical_expressionContext) interface{} {
	if ctx.Unary_logical_expression() != nil {
		return v.VisitUnary_logical_expression(ctx.Unary_logical_expression().(*plsql.Unary_logical_expressionContext))
	} else {
		expr := &semantic.BinaryExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		if ctx.AND() != nil {
			expr.Operator = "AND"
		} else if ctx.OR() != nil {
			expr.Operator = "OR"
		}
		var ok bool
		node := ctx.Logical_expression(0).(*plsql.Logical_expressionContext)
		expr.Left, ok = v.VisitLogical_expression(node).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}

		node = ctx.Logical_expression(1).(*plsql.Logical_expressionContext)
		expr.Right, ok = v.VisitLogical_expression(node).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}

		return expr
	}
}

func (v *exprVisitor) VisitUnary_logical_expression(ctx *plsql.Unary_logical_expressionContext) interface{} {
	result := &semantic.UnaryLogicalExpression{}
	result.SetLine(ctx.GetStart().GetLine())
	result.SetColumn(ctx.GetStart().GetColumn())
	expression := ctx.Multiset_expression().Accept(v)
	if expression != nil {
		var ok bool
		result.Expr, ok = expression.(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				ctx.Multiset_expression().GetStart().GetLine(),
				ctx.Multiset_expression().GetStart().GetColumn())
			return result
		}
	} else {
		name := &semantic.NameExpression{
			Name: ctx.Multiset_expression().GetText(),
		}
		name.SetLine(ctx.GetStart().GetLine())
		name.SetColumn(ctx.GetStart().GetColumn())
		result.Expr = name
	}

	if len(ctx.AllIS()) == 0 && len(ctx.AllNOT()) == 0 {
		return result.Expr
	}

	if len(ctx.AllIS()) == 1 {
		result.Operator = strings.ToUpper(ctx.Logical_operation(0).GetText())
		if len(ctx.AllNOT()) == 1 {
			result.Not = true
		}
	}

	return result
}

func (v *exprVisitor) VisitRelational_expression(ctx *plsql.Relational_expressionContext) interface{} {
	result := &semantic.RelationalExpression{}
	result.SetLine(ctx.GetStart().GetLine())
	result.SetColumn(ctx.GetStart().GetColumn())
	if len(ctx.AllRelational_expression()) > 0 {
		var ok bool
		var node antlr.ParserRuleContext
		node = ctx.Relational_expression(0)
		left := node.Accept(v)
		result.Left, ok = left.(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return result
		}
		node = ctx.Relational_expression(1)
		right := node.Accept(v)
		result.Right, ok = right.(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return result
		}
		result.Operator = strings.ToUpper(ctx.Relational_operator().GetText())
		return result
	} else {
		return v.VisitCompound_expression(ctx.Compound_expression().(*plsql.Compound_expressionContext))
	}
}

func (v *exprVisitor) VisitCompound_expression(ctx *plsql.Compound_expressionContext) interface{} {
	var ok bool
	if ctx.IN() != nil {
		expr := &semantic.InExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		node := ctx.Concatenation(0)
		expr.Expr, ok = v.VisitConcatenation(node.(*plsql.ConcatenationContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}
		right := ctx.In_elements().(*plsql.In_elementsContext)
		elems, ok := v.VisitIn_elements(right).([]semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				right.GetStart().GetLine(),
				right.GetStart().GetColumn())
			return expr
		}
		expr.Elems = elems
		return expr
	}

	if ctx.BETWEEN() != nil {
		expr := &semantic.BetweenExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		node := ctx.Concatenation(0).(*plsql.ConcatenationContext)
		expr.Expr, ok = v.VisitConcatenation(node).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}

		right := ctx.Between_elements().(*plsql.Between_elementsContext)
		expr.Elems, ok = v.VisitBetween_elements(right).([]semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				right.GetStart().GetLine(),
				right.GetStart().GetColumn())
			return expr
		}
		return expr
	}

	if ctx.GetLike_type() != nil {
		expr := &semantic.LikeExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		node := ctx.Concatenation(0).(*plsql.ConcatenationContext)
		expr.Expr, ok = v.VisitConcatenation(node).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}

		node = ctx.Concatenation(1).(*plsql.ConcatenationContext)
		expr.LikeExpr, ok = v.VisitConcatenation(node).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				node.GetStart().GetLine(),
				node.GetStart().GetColumn())
			return expr
		}
		return expr
	}

	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitIn_elements(ctx *plsql.In_elementsContext) interface{} {
	if ctx.AllConcatenation() != nil {
		elems := make([]semantic.Expr, 0, len(ctx.AllConcatenation()))
		for _, c := range ctx.AllConcatenation() {
			elem, ok := v.VisitConcatenation(c.(*plsql.ConcatenationContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					c.GetStart().GetLine(),
					c.GetStart().GetColumn())
				continue
			}
			elems = append(elems, elem)
		}
		return elems
	}

	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitBetween_elements(ctx *plsql.Between_elementsContext) interface{} {
	if ctx.AllConcatenation() != nil {
		elems := make([]semantic.Expr, 0, len(ctx.AllConcatenation()))
		for _, c := range ctx.AllConcatenation() {
			node := c.(*plsql.ConcatenationContext)
			var ok bool
			expr, ok := v.VisitConcatenation(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				continue
			}

			elems = append(elems, expr)
		}
		return elems
	}

	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitConcatenation(ctx *plsql.ConcatenationContext) interface{} {
	if ctx.Model_expression() != nil {
		return ctx.Model_expression().Accept(v)
	} else {
		if ctx.BAR(0) != nil {
			expr := &semantic.BinaryExpression{}
			expr.SetLine(ctx.GetStart().GetLine())
			expr.SetColumn(ctx.GetStart().GetColumn())

			var ok bool
			var node antlr.ParserRuleContext
			node = ctx.Concatenation(0)
			left := v.VisitConcatenation(node.(*plsql.ConcatenationContext))
			expr.Left, ok = left.(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			node = ctx.Concatenation(1)
			right := node.Accept(v)
			expr.Right, ok = right.(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			expr.Operator = "||"
			return expr
		} else if ctx.GetOp() != nil {
			expr := &semantic.BinaryExpression{}
			expr.SetLine(ctx.GetStart().GetLine())
			expr.SetColumn(ctx.GetStart().GetColumn())
			var ok bool
			var node antlr.ParserRuleContext

			node = ctx.Concatenation(0)
			left := v.VisitConcatenation(node.(*plsql.ConcatenationContext))
			expr.Left, ok = left.(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			node = ctx.Concatenation(1)
			right := v.VisitConcatenation(node.(*plsql.ConcatenationContext))
			expr.Right, ok = right.(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}
			expr.Operator = ctx.GetOp().GetText()
			return expr
		}
	}
	panic("unreachable")
}

func (v *exprVisitor) VisitUnary_expression(ctx *plsql.Unary_expressionContext) interface{} {
	if ctx.MINUS_SIGN() != nil || ctx.PLUS_SIGN() != nil {
		expr := &semantic.SignExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		var ok bool
		expr.Expr, ok = v.VisitUnary_expression(ctx.Unary_expression().(*plsql.Unary_expressionContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				ctx.Unary_expression().GetStart().GetLine(),
				ctx.Unary_expression().GetStart().GetColumn())
			return expr
		}
		var sign string
		if ctx.MINUS_SIGN() != nil {
			sign = "-"
		} else if ctx.PLUS_SIGN() != nil {
			sign = "+"
		}
		expr.Sign = sign
		return expr
	}
	if ctx.Case_statement() != nil {
		expr := &semantic.StatementExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		stmt := v.stmtVisitor.VisitCase_statement(ctx.Case_statement().(*plsql.Case_statementContext)).(*semantic.CaseWhenStatement)
		expr.Stmt = stmt
		return expr
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitQuantified_expression(ctx *plsql.Quantified_expressionContext) interface{} {
	if ctx.EXISTS() != nil {
		expr := &semantic.ExistsExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		if ctx.Select_only_statement() != nil {
			stmt := v.stmtVisitor.VisitSelect_only_statement(ctx.Select_only_statement().(*plsql.Select_only_statementContext)).(*semantic.SelectStatement)
			query := &semantic.QueryExpression{Query: stmt}
			expr.Expr = query
		}
		return expr
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitAtom(ctx *plsql.AtomContext) interface{} {
	if ctx.Outer_join_sign() != nil {
		expr := &semantic.OuterJoinExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		var ok bool
		expr.Expr, ok = v.VisitTable_element(ctx.Table_element().(*plsql.Table_elementContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				ctx.Table_element().GetStart().GetLine(),
				ctx.Table_element().GetStart().GetColumn())
			return expr
		}
		return expr
	}
	if ctx.Expressions() != nil {
		expressions, ok := v.VisitExpressions(ctx.Expressions().(*plsql.ExpressionsContext)).([]semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				ctx.Expressions().GetStart().GetLine(),
				ctx.Expressions().GetStart().GetColumn())
			return nil
		}
		// TODO: wrap multiple expressions with Expressions struct
		if len(expressions) == 1 {
			return expressions[0]
		}
		return expressions
	}
	if ctx.Subquery() != nil {
		expr := &semantic.StatementExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		stmt, ok := v.stmtVisitor.VisitSubquery(ctx.Subquery().(*plsql.SubqueryContext)).(semantic.Statement)
		if !ok {
			v.ReportError("unsupported statement",
				ctx.Subquery().GetStart().GetLine(),
				ctx.Subquery().GetStart().GetColumn())
			return expr
		}
		expr.Stmt = stmt
		return expr
	}
	if ctx.Bind_variable() != nil {
		var expr semantic.Expr
		var ok bool = false
		expr, ok = v.VisitBind_variable(ctx.Bind_variable().(*plsql.Bind_variableContext)).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				ctx.Bind_variable().GetStart().GetLine(),
				ctx.Bind_variable().GetStart().GetColumn())
			return nil
		}
		return expr
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitTable_element(ctx *plsql.Table_elementContext) interface{} {
	text := ctx.GetText()
	return v.parseDotExpr(text)
}

func (v *exprVisitor) VisitString_function(ctx *plsql.String_functionContext) interface{} {
	var ok bool
	if ctx.SUBSTR() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "SUBSTR"}}
		for _, arg := range ctx.AllExpression() {
			node, ok := arg.Accept(v).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					arg.GetStart().GetLine(),
					arg.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, node)
		}
		return expr
	}
	if ctx.TO_CHAR() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "TO_CHAR"}}
		var arg semantic.Expr
		if ctx.Table_element() != nil {
			arg, ok = v.VisitTable_element(ctx.Table_element().(*plsql.Table_elementContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Table_element().GetStart().GetLine(),
					ctx.Table_element().GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		} else if ctx.Standard_function() != nil {
			arg, ok = v.Visit(ctx.Standard_function().(*plsql.Standard_functionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Standard_function().GetStart().GetLine(),
					ctx.Standard_function().GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		} else if ctx.Expression(0) != nil {
			arg, ok = v.VisitExpression(ctx.Expression(0).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Expression(0).GetStart().GetLine(),
					ctx.Expression(0).GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		}
		for _, arg := range ctx.AllQuoted_string() {
			node, ok := v.VisitQuoted_string(arg.(*plsql.Quoted_stringContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					arg.GetStart().GetLine(),
					arg.GetStart().GetColumn())
				continue
			}
			expr.Args = append(expr.Args, node)
		}
		return expr
	}
	if ctx.TO_DATE() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "TO_DATE"}}
		var arg semantic.Expr
		if ctx.Table_element() != nil {
			arg, ok = v.VisitTable_element(ctx.Table_element().(*plsql.Table_elementContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Table_element().GetStart().GetLine(),
					ctx.Table_element().GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		} else if ctx.Standard_function() != nil {
			arg, ok = v.Visit(ctx.Standard_function().(*plsql.Standard_functionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Standard_function().GetStart().GetLine(),
					ctx.Standard_function().GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		} else if ctx.Expression(0) != nil {
			arg, ok = v.VisitExpression(ctx.Expression(0).(*plsql.ExpressionContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Expression(0).GetStart().GetLine(),
					ctx.Expression(0).GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		}
		for _, arg := range ctx.AllQuoted_string() {
			node, ok := v.VisitQuoted_string(arg.(*plsql.Quoted_stringContext)).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					arg.GetStart().GetLine(),
					arg.GetStart().GetColumn())
				continue
			}
			expr.Args = append(expr.Args, node)
		}
		return expr
	}
	if ctx.NVL() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "NVL"}}
		for _, arg := range ctx.AllExpression() {
			node, ok := arg.Accept(v).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					arg.GetStart().GetLine(),
					arg.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, node)
		}
		return expr
	}
	if ctx.DECODE() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "DECODE"}}
		{
			node, ok := ctx.Expressions().Accept(v).([]semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Expressions().GetStart().GetLine(),
					ctx.Expressions().GetStart().GetColumn())
				return expr
			}
			expr.Args = node
		}
		return expr
	}
	if ctx.TRIM() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "TRIM"}}
		{
			node, ok := ctx.Concatenation().Accept(v).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.Concatenation().GetStart().GetLine(),
					ctx.Concatenation().GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, node)
		}
		return expr
	}
	_ = ok
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitNumeric_function(ctx *plsql.Numeric_functionContext) interface{} {
	var ok bool
	if ctx.ROUND() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "ROUND"}}
		if arg := ctx.Expression(); arg != nil {
			node := arg.(*plsql.ExpressionContext)
			elems, ok := v.VisitExpression(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, elems)
		}
		return expr
	}
	if ctx.AVG() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "AVG"}}
		if arg := ctx.Expression(); arg != nil {
			node := arg.(*plsql.ExpressionContext)
			elems, ok := v.VisitExpression(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, elems)
		}
		return expr
	}
	if ctx.SUM() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "SUM"}}
		if arg := ctx.Expression(); arg != nil {
			node := arg.(*plsql.ExpressionContext)
			arg, ok := v.VisitExpression(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, arg)
		}
		return expr
	}
	if ctx.MAX() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "MAX"}}
		if arg := ctx.Expression(); arg != nil {
			node := arg.(*plsql.ExpressionContext)
			arg, ok := v.VisitExpression(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}

			expr.Args = append(expr.Args, arg)
		}
		return expr
	}
	if ctx.COUNT() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "COUNT"}}
		if concatenation := ctx.Concatenation(); concatenation != nil {
			node := concatenation.(*plsql.ConcatenationContext)
			arg, ok := v.VisitConcatenation(node).(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					node.GetStart().GetLine(),
					node.GetStart().GetColumn())
				return expr
			}
			expr.Args = append(expr.Args, arg)
		}
		if ctx.ASTERISK() != nil {
			expr.Args = append(expr.Args, &semantic.StringLiteral{Value: "*"})
		}
		return expr
	}
	_ = ok
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitBind_variable(ctx *plsql.Bind_variableContext) interface{} {
	var expr semantic.Expr
	if ctx.BINDVAR(0) == ctx.GetChild(0) {
		bindExpr := &semantic.BindNameExpression{
			Name: &semantic.NameExpression{
				Name: ctx.BINDVAR(0).GetText()}}
		bindExpr.SetLine(ctx.BINDVAR(0).GetSymbol().GetLine())
		bindExpr.SetColumn(ctx.BINDVAR(0).GetSymbol().GetColumn())
		expr = bindExpr
	}
	if len(ctx.AllGeneral_element_part()) > 0 {
		var dotExpr semantic.Expr
		for i, elem := range ctx.AllGeneral_element_part() {
			elemExpr := v.VisitGeneral_element_part(elem.(*plsql.General_element_partContext)).(semantic.Expr)
			if i == 0 {
				dotExpr = elemExpr
				continue
			}
			dotExpr = &semantic.DotExpression{
				Name:   elemExpr,
				Parent: dotExpr,
			}
		}
		switch dotExpr.(type) {
		case *semantic.NameExpression:
			dotExpr = &semantic.DotExpression{
				Name: dotExpr,
				Parent: &semantic.DotExpression{
					Name: expr,
				},
			}
		case *semantic.DotExpression:
			parent := dotExpr.(*semantic.DotExpression)
			for parent.Parent != nil {
				parent = parent.Parent.(*semantic.DotExpression)
			}
			parent.Parent = &semantic.DotExpression{
				Name: expr,
			}
		}
		expr = dotExpr
	}
	return expr
}

func (v *exprVisitor) VisitGeneral_element_part(ctx *plsql.General_element_partContext) interface{} {
	var dotExpr semantic.Expr
	if ctx.Id_expression(0) != nil {
		dotExpr = &semantic.NameExpression{Name: ctx.Id_expression(0).GetText()}
	}
	if len(ctx.AllId_expression()) > 1 {
		for i, id := range ctx.AllId_expression() {
			if i == 0 {
				dotExpr = &semantic.DotExpression{
					Name: dotExpr,
				}
				continue
			}
			dotExpr = &semantic.DotExpression{
				Name:   &semantic.NameExpression{Name: id.GetText()},
				Parent: dotExpr,
			}
		}
	}
	if ctx.Function_argument() != nil {
		args := v.VisitFunction_argument(ctx.Function_argument().(*plsql.Function_argumentContext)).([]interface{})
		expr := &semantic.FunctionCallExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		for _, arg := range args {
			var ok bool
			elems, ok := arg.(semantic.Expr)
			if !ok {
				v.ReportError("unsupported expression",
					ctx.GetStart().GetLine(),
					ctx.GetStart().GetColumn())
				continue
			}

			expr.Args = append(expr.Args, elems)
		}
		expr.Name = dotExpr
		return expr
	} else {
		return dotExpr
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitFunction_argument(ctx *plsql.Function_argumentContext) interface{} {
	args := make([]interface{}, 0)
	if len(ctx.AllArgument()) > 0 {
		for _, arg := range ctx.AllArgument() {
			args = append(args, arg.Accept(v))
		}
	}
	return args
}

func (v *exprVisitor) VisitQuoted_string(ctx *plsql.Quoted_stringContext) interface{} {
	if ctx.Variable_name() != nil {
		return v.VisitChildren(ctx)
	}
	return &semantic.StringLiteral{Value: ctx.GetText()}
}

func (v *exprVisitor) VisitConstant(ctx *plsql.ConstantContext) interface{} {
	if ctx.NULL_() != nil {
		return &semantic.NullExpression{}
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitVariable_name(ctx *plsql.Variable_nameContext) interface{} {
	parts := strings.Split(ctx.GetText(), ".")
	if len(parts) == 1 {
		name := &semantic.NameExpression{
			Name: ctx.GetText(),
		}
		name.SetLine(ctx.GetStart().GetLine())
		name.SetColumn(ctx.GetStart().GetColumn())
		return name
	}

	return v.parseDotExpr(ctx.GetText())
}

func (v *exprVisitor) VisitNumeric(ctx *plsql.NumericContext) interface{} {
	number := &semantic.NumericLiteral{}
	number.SetLine(ctx.GetStart().GetLine())
	number.SetColumn(ctx.GetStart().GetColumn())
	if ctx.UNSIGNED_INTEGER() != nil {
		if v, err := strconv.ParseInt(ctx.GetText(), 10, 64); err == nil {
			number.Value = v
		}
	}
	return number
}

func (v *exprVisitor) VisitRoutine_name(ctx *plsql.Routine_nameContext) interface{} {
	return v.parseDotExpr(ctx.GetText())
}

func (v *exprVisitor) VisitSelect_list_elements(ctx *plsql.Select_list_elementsContext) interface{} {
	if ctx.Column_alias() != nil {
		expr := &semantic.AliasExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		expr.Expr = ctx.Expression().Accept(v).(semantic.Expr)
		if ctx.Column_alias().Quoted_string() != nil {
			expr.Alias = ctx.Column_alias().Quoted_string().GetText()
		} else if ctx.Column_alias().Identifier() != nil {
			expr.Alias = ctx.Column_alias().Identifier().GetText()
		}
		return expr
	}

	if ctx.Tableview_name() != nil {
		tableName := ctx.Tableview_name().GetText()
		return &semantic.SelectField{WildCard: &semantic.WildCardField{Table: tableName, Schema: "*"}}
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitFor_update_options(ctx *plsql.For_update_optionsContext) interface{} {
	expr := &semantic.ForUpdateOptionsExpression{}

	if ctx.SKIP_() != nil {
		expr.SkipLocked = true
	}

	if ctx.NOWAIT() != nil {
		expr.NoWait = true
	}

	if ctx.WAIT() != nil {
		v.ReportError("unsupported expression", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
	}

	return expr
}

func (v *exprVisitor) VisitGeneral_table_ref(ctx *plsql.General_table_refContext) interface{} {
	var expr semantic.Expr
	if ctx.ONLY() != nil {
		v.ReportError(
			"unsupported expression",
			ctx.ONLY().GetSymbol().GetLine(),
			ctx.ONLY().GetSymbol().GetColumn(),
		)
		return expr
	}

	if ctx.Dml_table_expression_clause() != nil {
		object := ctx.Dml_table_expression_clause().Accept(v)
		expr = object.(semantic.Expr)
	}

	if ctx.Table_alias() != nil {
		expr = &semantic.AliasExpression{
			Expr:  expr,
			Alias: ctx.Table_alias().GetText(),
		}
	}
	return expr
}

func (v *exprVisitor) VisitDml_table_expression_clause(ctx *plsql.Dml_table_expression_clauseContext) interface{} {
	if ctx.Table_collection_expression() != nil {
		v.ReportError(
			"unsupported expression",
			ctx.Table_collection_expression().GetStart().GetLine(),
			ctx.Table_collection_expression().GetStart().GetColumn(),
		)
		return nil
	} else if ctx.Json_table_clause() != nil {
		v.ReportError(
			"unsupported expression",
			ctx.Json_table_clause().GetStart().GetLine(),
			ctx.Json_table_clause().GetStart().GetColumn(),
		)
		return nil
	} else if ctx.Select_statement() != nil {
		v.ReportError(
			"unsupported expression",
			ctx.Select_statement().GetStart().GetLine(),
			ctx.Select_statement().GetStart().GetColumn(),
		)
		return nil
	} else {
		expr := &semantic.NameExpression{Name: ctx.Tableview_name().GetText()}
		if ctx.Sample_clause() != nil {
			v.ReportError(
				"unsupported expression",
				ctx.Sample_clause().GetStart().GetLine(),
				ctx.Sample_clause().GetStart().GetColumn(),
			)
		}
		return expr
	}
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitColumn_based_update_set_clause(ctx *plsql.Column_based_update_set_clauseContext) interface{} {
	if ctx.Paren_column_list() != nil {
		v.ReportError(
			"unsupported expression",
			ctx.Paren_column_list().GetStart().GetLine(),
			ctx.Paren_column_list().GetStart().GetColumn(),
		)
		return nil
	} else {
		var ok bool
		expr := &semantic.BinaryExpression{Operator: "="}
		object := v.parseDotExpr(ctx.Column_name().GetText())
		expr.Left, ok = object.(semantic.Expr)
		if !ok {
			v.ReportError(
				"unsupported expression",
				ctx.Column_name().GetStart().GetLine(),
				ctx.Column_name().GetStart().GetColumn(),
			)
		}
		expr.Right, ok = ctx.Expression().Accept(v).(semantic.Expr)
		if !ok {
			v.ReportError(
				"unsupported expression",
				ctx.Expression().GetStart().GetLine(),
				ctx.Expression().GetStart().GetColumn(),
			)
		}
		return expr
	}
}

func (v *exprVisitor) VisitSynonym_name(ctx *plsql.Synonym_nameContext) interface{} {
	return v.parseDotExpr(ctx.GetText())
}

func (v *exprVisitor) VisitSchema_object_name(ctx *plsql.Schema_object_nameContext) interface{} {
	return v.parseDotExpr(ctx.GetText())
}

func (v *exprVisitor) VisitParen_column_list(ctx *plsql.Paren_column_listContext) interface{} {
	return ctx.Column_list().Accept(v)
}

func (v *exprVisitor) VisitColumn_list(ctx *plsql.Column_listContext) interface{} {
	exprs := make([]semantic.Expr, 0)
	for _, col := range ctx.AllColumn_name() {
		name := col.GetText()
		expr, ok := v.parseDotExpr(name).(semantic.Expr)
		if !ok {
			v.ReportError("unsupported expression",
				col.GetStart().GetLine(),
				col.GetStart().GetColumn())
			continue
		}
		exprs = append(exprs, expr)
	}
	return exprs
}

func (v *exprVisitor) parseDotExpr(text string) semantic.Expr {
	parts := strings.Split(text, ".")
	if len(parts) == 1 {
		return &semantic.NameExpression{
			Name: text,
		}
	}

	var dotExpr semantic.Expr
	for _, part := range parts {
		dotExpr = &semantic.DotExpression{
			Name:   &semantic.NameExpression{Name: part},
			Parent: dotExpr,
		}
	}
	return dotExpr
}
