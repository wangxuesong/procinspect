//go:generate go run procinspect/pkg/semantic/internal -o types.generated.go
package semantic

import (
	"errors"
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

func VisitChildren(v NodeVisitor, node AstNode) (err error) {
	visitor := reflect.ValueOf(v)
	tn := reflect.TypeOf(node)
	vn := reflect.ValueOf(node)
	if tn.Kind() == reflect.Ptr {
		tn = tn.Elem()
		vn = vn.Elem()
	}
	for i := 0; i < tn.NumField(); i++ {
		field := tn.Field(i)
		value := vn.Field(i)
		if field.Name == "SyntaxNode" {
			continue
		}
		if field.Type.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if field.Type.Kind() == reflect.Slice {
			for j := 0; j < value.Len(); j++ {
				element := value.Index(j)
				if element.Kind() == reflect.Interface {
					element = element.Elem()
				}
			try:
				acceptMethod := element.MethodByName("Accept")

				if acceptMethod.IsValid() {
					vs := acceptMethod.Call([]reflect.Value{visitor})
					if vs[0].IsNil() {
						continue
					}
					err = errors.Join(err, vs[0].Interface().(error))
				} else {
					element = element.Addr()
					goto try
				}
			}
		} else if field.Type.Kind() == reflect.Struct {
			element := value
			if element.Kind() == reflect.Interface {
				element = element.Elem()
			}
		retry:
			acceptMethod := element.MethodByName("Accept")

			if acceptMethod.IsValid() {
				vs := acceptMethod.Call([]reflect.Value{visitor})
				if vs[0].IsNil() {
					continue
				}
				err = errors.Join(err, vs[0].Interface().(error))
			} else {
				element = element.Addr()
				goto retry
			}
		}
	}
	return
}
