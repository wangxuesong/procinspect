package checker

import (
	"reflect"

	"procinspect/pkg/semantic"
)

type rule struct {
	name      string
	target    semantic.Node
	checkFunc checkFunc
	message   string
}

type checkFunc func(r rule, node semantic.Node) error

var ruleMap = map[reflect.Type]rule{
	reflect.TypeOf(semantic.CreateNestTableStatement{}): {
		name:   "nest table type declaration",
		target: &semantic.NestTableTypeDeclaration{},
		checkFunc: func(r rule, node semantic.Node) error {
			return SqlValidationError{Line: node.Line(), Msg: r.message}
		},
		message: "unsupported: nest table type declaration",
	},
	reflect.TypeOf(semantic.NestTableTypeDeclaration{}): {
		name:   "create nest table type",
		target: &semantic.CreateNestTableStatement{},
		checkFunc: func(r rule, node semantic.Node) error {
			return SqlValidationError{Line: node.Line(), Msg: r.message}
		},
		message: "unsupported: nest table type declaration",
	},
	// reflect.TypeOf(semantic.TableRef{}): {
	// 	name:   "select from dblink",
	// 	target: &semantic.SelectStatement{},
	// 	checkFunc: func(r rule, node semantic.Node) error {
	// 		table := node.(*semantic.TableRef)
	// 		if strings.Index(table.Table, "@") >= 0 {
	// 			return SqlValidationError{Line: node.Line(), Msg: r.message}
	// 		}
	// 		return nil
	// 	},
	// 	message: "unsupported: select from dblink",
	// },
	reflect.TypeOf(semantic.UpdateStatement{}): {
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
