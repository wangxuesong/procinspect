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
	if len(exprs) == 1 {
		return exprs[0]
	}
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
		return v.VisitExpressions(ctx.Expressions().(*plsql.ExpressionsContext))
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
	_ = ok
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) VisitGeneral_element_part(ctx *plsql.General_element_partContext) interface{} {
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
		expr.Name = v.parseDotExpr(ctx.Id_expression(0).GetText())
		return expr
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
	return v.VisitChildren(ctx)
}

func (v *exprVisitor) parseDotExpr(text string) semantic.Expr {
	parts := strings.Split(text, ".")
	if len(parts) == 1 {
		return &semantic.NameExpression{
			Name: text,
		}
	}

	length := len(parts)
	var expr semantic.Expr = &semantic.NameExpression{
		Name: parts[0],
	}

	for i := 1; i < length; i++ {
		expr = &semantic.DotExpression{
			Name:   parts[i],
			Parent: expr,
		}
	}

	return expr
}
