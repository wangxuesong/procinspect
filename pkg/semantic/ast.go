//go:generate go run procinspect/pkg/semantic/internal -o types.generated.go
package semantic

type Expression interface {
	ExprAccept(visitor ExprVisitor) (result interface{}, err error)
}

type Stmt interface {
	StmtAccept(visitor StmtVisitor) (err error)
}

type AstNode interface {
	Accept(visitor NodeVisitor) (err error)
}
