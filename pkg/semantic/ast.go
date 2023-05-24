//go:generate go run internal/gen-ast-types.go -o types.generated.go
package semantic

type Expression interface {
	Accept(visitor ExprVisitor) (result interface{}, err error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) (err error)
}
