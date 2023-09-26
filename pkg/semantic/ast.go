//go:generate go run procinspect/pkg/semantic/internal -o types.generated.go
package semantic

import (
	"reflect"
)

type Expression interface {
	ExprAccept(visitor ExprVisitor) (result interface{}, err error)
}

type Stmt interface {
	StmtAccept(visitor StmtVisitor) (err error)
}

type AstNode interface {
	Accept(visitor NodeVisitor) (err error)
}

func GetChildren(node AstNode) []AstNode {
	// 创建一个 AstNode 类型的 slice，作为返回值
	var children []AstNode

	// 获取 node（你的 struct） 的 Value
	rv := reflect.ValueOf(node)

	// rv 可能是一个指针，如果是的话我们需要获取它指向的实际对象
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 现在 rv 应该是一个 struct，我们遍历它的所有字段
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)

		if field.IsZero() {
			continue
		}

		if field.Kind() == reflect.Slice { // 如果字段是一个 slice
			// 遍历 slice 的每个元素
			for j := 0; j < field.Len(); j++ {
				// 如果元素是一个 AstNode，我们添加它到 children 中
				if field.Index(j).Type().Implements(reflect.TypeOf((*AstNode)(nil)).Elem()) {
					children = append(children, field.Index(j).Interface().(AstNode))
				}
			}
		} else if field.Kind() == reflect.Ptr { // 如果字段是一个指针
			if field.IsNil() {
				continue
			}
			field = field.Elem() // 获取指针指向的实际对象
		}
		if field.Kind() == reflect.Struct || field.Kind() == reflect.Interface { // 如果字段是一个 struct
			if field.IsValid() { // 字段是否有效
				// 如果字段是一个 AstNode，我们添加它到 children 中
				if field.Type().Implements(reflect.TypeOf((*AstNode)(nil)).Elem()) {
					e := field.Interface()
					if e != nil {
						children = append(children, e.(AstNode))
					}
				}
			}
		}
	}

	// 返回所有的 children
	return children
}
