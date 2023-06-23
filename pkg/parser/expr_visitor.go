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
		//plsql.BasePlSqlParserVisitor
	}
)

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
		case *plsql.ConcatenationContext:
			c := child.(*plsql.ConcatenationContext)
			nodes = append(nodes, v.VisitConcatenation(c))
		case *plsql.AtomContext:
			c := child.(*plsql.AtomContext)
			nodes = append(nodes, v.VisitAtom(c))
		case *plsql.String_functionContext:
			c := child.(*plsql.String_functionContext)
			nodes = append(nodes, v.VisitString_function(c))
		case *plsql.General_element_partContext:
			c := child.(*plsql.General_element_partContext)
			nodes = append(nodes, v.VisitGeneral_element_part(c))
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
		expr := v.VisitExpression(e.(*plsql.ExpressionContext)).(semantic.Expr)
		exprs = append(exprs, expr)
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	return exprs
}

func (v *exprVisitor) VisitExpression(ctx *plsql.ExpressionContext) interface{} {
	return ctx.Accept(v)
}

func (v *exprVisitor) VisitOther_function(ctx *plsql.Other_functionContext) interface{} {
	switch ctx.GetChild(0).(type) {
	case *plsql.Cursor_nameContext:
		node := ctx.GetChild(0).(*plsql.Cursor_nameContext)
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
	default:
		return v.VisitChildren(ctx)
	}
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
		expr.Left = v.VisitLogical_expression(ctx.Logical_expression(0).(*plsql.Logical_expressionContext)).(semantic.Expr)
		expr.Right = v.VisitLogical_expression(ctx.Logical_expression(1).(*plsql.Logical_expressionContext)).(semantic.Expr)
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
		return ctx.Compound_expression().Accept(v)
	}
}

func (v *exprVisitor) VisitConcatenation(ctx *plsql.ConcatenationContext) interface{} {
	if ctx.Model_expression() != nil {
		return ctx.Accept(v)
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
	return ctx.Accept(v)
}

func (v *exprVisitor) VisitAtom(ctx *plsql.AtomContext) interface{} {
	if ctx.Expressions() != nil {
		return v.VisitExpressions(ctx.Expressions().(*plsql.ExpressionsContext))
	}
	return ctx.Accept(v)
}

func (v *exprVisitor) VisitString_function(ctx *plsql.String_functionContext) interface{} {
	if ctx.SUBSTR() != nil {
		expr := &semantic.FunctionCallExpression{Name: &semantic.NameExpression{Name: "SUBSTR"}}
		for _, arg := range ctx.AllExpression() {
			expr.Args = append(expr.Args, arg.Accept(v).(semantic.Expr))
		}
		return expr
	}
	return ctx.Accept(v)
}

func (v *exprVisitor) VisitGeneral_element_part(ctx *plsql.General_element_partContext) interface{} {
	if ctx.Function_argument() != nil {
		args := v.VisitFunction_argument(ctx.Function_argument().(*plsql.Function_argumentContext)).([]interface{})
		expr := &semantic.FunctionCallExpression{}
		expr.SetLine(ctx.GetStart().GetLine())
		expr.SetColumn(ctx.GetStart().GetColumn())
		for _, arg := range args {
			expr.Args = append(expr.Args, arg.(semantic.Expr))
		}
		expr.Name = v.parseDotExpr(ctx.Id_expression(0).GetText())
		return expr
	}
	return ctx.Accept(v)
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
		return ctx.Accept(v)
	}
	return &semantic.StringLiteral{Value: ctx.GetText()}
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
