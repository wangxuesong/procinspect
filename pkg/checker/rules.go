package checker

import (
	"context"

	"github.com/go-playground/validator/v10"

	"procinspect/pkg/semantic"
)

type rule struct {
	name      string
	target    semantic.Node
	validFunc warpFunc
	message   string
}

type warpFunc func(r rule) validator.StructLevelFuncCtx

var rules = []rule{
	{
		"nest table type declaration",
		&semantic.NestTableTypeDeclaration{},
		validateNestTableTypeDeclaration,
		"unsupported: nest table type declaration",
	},
	{
		"create nest table type",
		&semantic.CreateNestTableStatement{},
		validateCreateNestTableType,
		"unsupported: nest table type declaration",
	},
	{
		"update set multiple columns with select",
		&semantic.UpdateStatement{},
		validateUpdateStatement,
		"unsupported: update set multiple columns with select",
	},
}

func validateUpdateStatement(r rule) validator.StructLevelFuncCtx {
	return func(ctx context.Context, sl validator.StructLevel) {
		cur := sl.Current().Interface().(semantic.UpdateStatement)
		binary, ok := cur.SetExprs[0].(*semantic.BinaryExpression)
		if ok {
			_, left := binary.Left.(*semantic.ExprListExpression)
			_, right := binary.Right.(*semantic.StatementExpression)
			if left && right {
				sl.ReportError(
					cur.SetExprs[0],
					"SetExpres",
					"SetExpres",
					r.name,
					r.message,
				)
			}
		}
	}
}

func validateNestTableTypeDeclaration(r rule) validator.StructLevelFuncCtx {
	return func(ctx context.Context, sl validator.StructLevel) {
		cur := sl.Current().Interface().(semantic.NestTableTypeDeclaration)
		sl.ReportError(cur, cur.Name, cur.Name, r.name, r.message)
	}
}

func validateCreateNestTableType(r rule) validator.StructLevelFuncCtx {
	return func(ctx context.Context, sl validator.StructLevel) {
		cur := sl.Current().Interface().(semantic.CreateNestTableStatement)
		sl.ReportError(cur, cur.Name, cur.Name, r.name, r.message)
	}
}
