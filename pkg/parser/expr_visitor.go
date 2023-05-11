package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	plsql "procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
	"strconv"
	"strings"
)

type (
	experVisitor struct {
		plsql.BasePlSqlParserVisitor
	}
)

func (v *experVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	children := node.GetChildren()
	nodes := make([]interface{}, 0, len(children))
	for _, child := range children {
		switch child.(type) {
		case *plsql.Other_functionContext:
			c := child.(*plsql.Other_functionContext)
			nodes = append(nodes, v.VisitOther_function(c))
		case *plsql.Unary_logical_expressionContext:
			c := child.(*plsql.Unary_logical_expressionContext)
			nodes = append(nodes, v.VisitUnary_logical_expression(c))
		case *plsql.Relational_expressionContext:
			c := child.(*plsql.Relational_expressionContext)
			nodes = append(nodes, v.VisitRelational_expression(c))
		case *plsql.Variable_nameContext:
			c := child.(*plsql.Variable_nameContext)
			nodes = append(nodes, v.VisitVariable_name(c))
		case *plsql.NumericContext:
			c := child.(*plsql.NumericContext)
			nodes = append(nodes, v.VisitNumeric(c))
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

func (v *experVisitor) VisitCondition(ctx *plsql.ConditionContext) interface{} {
	if ctx.Expression() != nil {
		expr := ctx.Expression().Accept(v)
		switch expr.(type) {
		case []interface{}:
			cc := expr.([]interface{})
			return cc[0]
		default:
			return expr
		}
	} else {
		return nil
	}
}

func (v *experVisitor) VisitOther_function(ctx *plsql.Other_functionContext) interface{} {
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

func (v *experVisitor) VisitUnary_logical_expression(ctx *plsql.Unary_logical_expressionContext) interface{} {
	result := &semantic.UnaryLogicalExpression{}
	result.SetLine(ctx.GetStart().GetLine())
	result.SetColumn(ctx.GetStart().GetColumn())
	expression := ctx.Multiset_expression().Accept(v)
	if expression != nil {
		result.Expr = expression.(semantic.Expr)
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

func (v *experVisitor) VisitRelational_expression(ctx *plsql.Relational_expressionContext) interface{} {
	result := &semantic.RelationalExpression{}
	result.SetLine(ctx.GetStart().GetLine())
	result.SetColumn(ctx.GetStart().GetColumn())
	if len(ctx.AllRelational_expression()) > 0 {
		left := ctx.Relational_expression(0).Accept(v)
		right := ctx.Relational_expression(1).Accept(v)
		result.Left = left.(semantic.Expr)
		result.Right = right.(semantic.Expr)
		result.Operator = strings.ToUpper(ctx.Relational_operator().GetText())
		return result
	} else {
		return ctx.Compound_expression().Accept(v)
	}
}

func (v *experVisitor) VisitVariable_name(ctx *plsql.Variable_nameContext) interface{} {
	name := &semantic.NameExpression{
		Name: ctx.GetText(),
	}
	name.SetLine(ctx.GetStart().GetLine())
	name.SetColumn(ctx.GetStart().GetColumn())
	return name
}

func (v *experVisitor) VisitNumeric(ctx *plsql.NumericContext) interface{} {
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
