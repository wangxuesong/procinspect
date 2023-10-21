package checker

import (
	"github.com/antonmedv/expr"

	"procinspect/pkg/semantic"
)

type Rule struct {
	Name      string
	Target    semantic.Node
	CheckFunc checkFunc
	Message   string
}

type checkFunc func(r Rule, node semantic.Node) error

var rules = []Rule{{
	Name:   "nest table type declaration",
	Target: &semantic.NestTableTypeDeclaration{},
	CheckFunc: func(r Rule, node semantic.Node) error {
		return SqlValidationError{Line: node.Line(), Msg: r.Message}
	},
	Message: "unsupported: nest table type declaration",
}, {
	Name:   "create nest table type",
	Target: &semantic.CreateNestTableStatement{},
	CheckFunc: func(r Rule, node semantic.Node) error {
		return SqlValidationError{Line: node.Line(), Msg: r.Message}
	},
	Message: "unsupported: nest table type declaration",
}, {
	Name:   "update set multiple columns with select",
	Target: &semantic.UpdateStatement{},
	CheckFunc: func(r Rule, node semantic.Node) error {
		stmt := node.(*semantic.UpdateStatement)
		binary, ok := stmt.SetExprs[0].(*semantic.BinaryExpression)
		if ok {
			_, left := binary.Left.(*semantic.ExprListExpression)
			_, right := binary.Right.(*semantic.StatementExpression)
			if left && right {
				return SqlValidationError{Line: node.Line(), Msg: r.Message}
			}
		}
		return nil
	},
	Message: "unsupported: update set multiple columns with select",
},
}

func ValidExprFunc(validRule string) checkFunc {
	return func(r Rule, node semantic.Node) error {
		env := map[string]any{
			"node": node,
		}
		program, err := expr.Compile(validRule, expr.Env(env), expr.AsBool())
		if err != nil {
			return err
		}
		output, err := expr.Run(program, env)
		if err != nil {
			return err
		}
		if output.(bool) {
			return SqlValidationError{Line: node.Line(), Msg: r.Message}
		}
		return nil
	}
}
