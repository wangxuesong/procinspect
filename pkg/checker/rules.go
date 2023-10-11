package checker

import (
	"reflect"

	"github.com/antonmedv/expr"

	"procinspect/pkg/semantic"
)

type rule struct {
	name      string
	target    semantic.Node
	checkFunc checkFunc
	message   string
}

type checkFunc func(r rule, node semantic.Node) error

var rules = []rule{{
	name:   "nest table type declaration",
	target: &semantic.NestTableTypeDeclaration{},
	checkFunc: func(r rule, node semantic.Node) error {
		return SqlValidationError{Line: node.Line(), Msg: r.message}
	},
	message: "unsupported: nest table type declaration",
}, {
	name:   "create nest table type",
	target: &semantic.CreateNestTableStatement{},
	checkFunc: func(r rule, node semantic.Node) error {
		return SqlValidationError{Line: node.Line(), Msg: r.message}
	},
	message: "unsupported: nest table type declaration",
}, {
	name:   "update set multiple columns with select",
	target: &semantic.UpdateStatement{},
	checkFunc: func(r rule, node semantic.Node) error {
		stmt := node.(*semantic.UpdateStatement)
		binary, ok := stmt.SetExprs[0].(*semantic.BinaryExpression)
		if ok {
			_, left := binary.Left.(*semantic.ExprListExpression)
			_, right := binary.Right.(*semantic.StatementExpression)
			if left && right {
				return SqlValidationError{Line: node.Line(), Msg: r.message}
			}
		}
		return nil
	},
	message: "unsupported: update set multiple columns with select",
},
}
var ruleMap map[reflect.Type]rule

func validDblinkFunc(validRule string) checkFunc {
	return func(r rule, node semantic.Node) error {
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
			return SqlValidationError{Line: node.Line(), Msg: r.message}
		}
		return nil
	}
}
