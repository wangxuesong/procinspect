//go:generate go run ../scripts/pkgreflect.go -nofuncs -novars -norecurs -noconsts -gofile=ast.gen.go ../pkg/semantic
package main

import (
	"fmt"
	"reflect"
	"sort"

	"procinspect/pkg/semantic"
)

func main() {

	nodeType := reflect.TypeOf((*semantic.Declaration)(nil)).Elem()
	names := make([]string, 0)

	for name, t := range AstTypes {
		if t.Kind() == reflect.Struct && reflect.PtrTo(t).Implements(nodeType) {
			//fmt.Println(name)
			names = append(names, name)
		}
	}

	sort.Strings(names)
	fmt.Println("var stmtTypes = Types{")
	for _, n := range names {
		fmt.Println("\t{")
		fmt.Printf("\t\tName:    \"%s\",\n", n)
		fmt.Printf("\t\tFields:  \"%s.%s\",\n", "semantic", n)
		fmt.Printf("\t\tComment: \"\",\n")
		fmt.Println("\t},")
	}
	fmt.Println("}")
	//fmt.Println(len(names))
}

var AstTypes = map[string]reflect.Type{
	"AliasExpression":                  reflect.TypeOf((*semantic.AliasExpression)(nil)).Elem(),
	"Argument":                         reflect.TypeOf((*semantic.Argument)(nil)).Elem(),
	"AssignmentStatement":              reflect.TypeOf((*semantic.AssignmentStatement)(nil)).Elem(),
	"AutonomousTransactionDeclaration": reflect.TypeOf((*semantic.AutonomousTransactionDeclaration)(nil)).Elem(),
	"BetweenExpression":                reflect.TypeOf((*semantic.BetweenExpression)(nil)).Elem(),
	"BinaryExpression":                 reflect.TypeOf((*semantic.BinaryExpression)(nil)).Elem(),
	"BindNameExpression":               reflect.TypeOf((*semantic.BindNameExpression)(nil)).Elem(),
	"BlockStatement":                   reflect.TypeOf((*semantic.BlockStatement)(nil)).Elem(),
	"Body":                             reflect.TypeOf((*semantic.Body)(nil)).Elem(),
	"CaseWhenBlock":                    reflect.TypeOf((*semantic.CaseWhenBlock)(nil)).Elem(),
	"CaseWhenStatement":                reflect.TypeOf((*semantic.CaseWhenStatement)(nil)).Elem(),
	"CastExpression":                   reflect.TypeOf((*semantic.CastExpression)(nil)).Elem(),
	"CloseStatement":                   reflect.TypeOf((*semantic.CloseStatement)(nil)).Elem(),
	"CommitStatement":                  reflect.TypeOf((*semantic.CommitStatement)(nil)).Elem(),
	"ContinueStatement":                reflect.TypeOf((*semantic.ContinueStatement)(nil)).Elem(),
	"CreateFunctionStatement":          reflect.TypeOf((*semantic.CreateFunctionStatement)(nil)).Elem(),
	"CreateNestTableStatement":         reflect.TypeOf((*semantic.CreateNestTableStatement)(nil)).Elem(),
	"CreatePackageBodyStatement":       reflect.TypeOf((*semantic.CreatePackageBodyStatement)(nil)).Elem(),
	"CreatePackageStatement":           reflect.TypeOf((*semantic.CreatePackageStatement)(nil)).Elem(),
	"CreateProcedureStatement":         reflect.TypeOf((*semantic.CreateProcedureStatement)(nil)).Elem(),
	"CreateSynonymStatement":           reflect.TypeOf((*semantic.CreateSynonymStatement)(nil)).Elem(),
	"CreateTypeStatement":              reflect.TypeOf((*semantic.CreateTypeStatement)(nil)).Elem(),
	"CursorAttribute":                  reflect.TypeOf((*semantic.CursorAttribute)(nil)).Elem(),
	"CursorDeclaration":                reflect.TypeOf((*semantic.CursorDeclaration)(nil)).Elem(),
	"Declaration":                      reflect.TypeOf((*semantic.Declaration)(nil)).Elem(),
	"DeleteStatement":                  reflect.TypeOf((*semantic.DeleteStatement)(nil)).Elem(),
	"DotExpression":                    reflect.TypeOf((*semantic.DotExpression)(nil)).Elem(),
	"ElseBlock":                        reflect.TypeOf((*semantic.ElseBlock)(nil)).Elem(),
	"ExceptionDeclaration":             reflect.TypeOf((*semantic.ExceptionDeclaration)(nil)).Elem(),
	"ExecuteImmediateStatement":        reflect.TypeOf((*semantic.ExecuteImmediateStatement)(nil)).Elem(),
	"ExistsExpression":                 reflect.TypeOf((*semantic.ExistsExpression)(nil)).Elem(),
	"ExitStatement":                    reflect.TypeOf((*semantic.ExitStatement)(nil)).Elem(),
	"Expr":                             reflect.TypeOf((*semantic.Expr)(nil)).Elem(),
	"ExprVisitor":                      reflect.TypeOf((*semantic.ExprVisitor)(nil)).Elem(),
	"Expression":                       reflect.TypeOf((*semantic.Expression)(nil)).Elem(),
	"FetchStatement":                   reflect.TypeOf((*semantic.FetchStatement)(nil)).Elem(),
	"FieldList":                        reflect.TypeOf((*semantic.FieldList)(nil)).Elem(),
	"ForUpdateClause":                  reflect.TypeOf((*semantic.ForUpdateClause)(nil)).Elem(),
	"ForUpdateOptionsExpression":       reflect.TypeOf((*semantic.ForUpdateOptionsExpression)(nil)).Elem(),
	"FromClause":                       reflect.TypeOf((*semantic.FromClause)(nil)).Elem(),
	"FunctionCallExpression":           reflect.TypeOf((*semantic.FunctionCallExpression)(nil)).Elem(),
	"FunctionDeclaration":              reflect.TypeOf((*semantic.FunctionDeclaration)(nil)).Elem(),
	"IfStatement":                      reflect.TypeOf((*semantic.IfStatement)(nil)).Elem(),
	"InExpression":                     reflect.TypeOf((*semantic.InExpression)(nil)).Elem(),
	"InsertStatement":                  reflect.TypeOf((*semantic.InsertStatement)(nil)).Elem(),
	"IntoClause":                       reflect.TypeOf((*semantic.IntoClause)(nil)).Elem(),
	"LikeExpression":                   reflect.TypeOf((*semantic.LikeExpression)(nil)).Elem(),
	"LoopStatement":                    reflect.TypeOf((*semantic.LoopStatement)(nil)).Elem(),
	"NameExpression":                   reflect.TypeOf((*semantic.NameExpression)(nil)).Elem(),
	"NestTableTypeDeclaration":         reflect.TypeOf((*semantic.NestTableTypeDeclaration)(nil)).Elem(),
	"Node":                             reflect.TypeOf((*semantic.Node)(nil)).Elem(),
	"NodeType":                         reflect.TypeOf((*semantic.NodeType)(nil)).Elem(),
	"NullExpression":                   reflect.TypeOf((*semantic.NullExpression)(nil)).Elem(),
	"NullStatement":                    reflect.TypeOf((*semantic.NullStatement)(nil)).Elem(),
	"NumericLiteral":                   reflect.TypeOf((*semantic.NumericLiteral)(nil)).Elem(),
	"OpenStatement":                    reflect.TypeOf((*semantic.OpenStatement)(nil)).Elem(),
	"OuterJoinExpression":              reflect.TypeOf((*semantic.OuterJoinExpression)(nil)).Elem(),
	"Parameter":                        reflect.TypeOf((*semantic.Parameter)(nil)).Elem(),
	"ProcedureCall":                    reflect.TypeOf((*semantic.ProcedureCall)(nil)).Elem(),
	"QueryExpression":                  reflect.TypeOf((*semantic.QueryExpression)(nil)).Elem(),
	"RelationalExpression":             reflect.TypeOf((*semantic.RelationalExpression)(nil)).Elem(),
	"ReturnStatement":                  reflect.TypeOf((*semantic.ReturnStatement)(nil)).Elem(),
	"RollbackStatement":                reflect.TypeOf((*semantic.RollbackStatement)(nil)).Elem(),
	"Script":                           reflect.TypeOf((*semantic.Script)(nil)).Elem(),
	"SelectField":                      reflect.TypeOf((*semantic.SelectField)(nil)).Elem(),
	"SelectStatement":                  reflect.TypeOf((*semantic.SelectStatement)(nil)).Elem(),
	"SetPosition":                      reflect.TypeOf((*semantic.SetPosition)(nil)).Elem(),
	"SignExpression":                   reflect.TypeOf((*semantic.SignExpression)(nil)).Elem(),
	"Span":                             reflect.TypeOf((*semantic.Span)(nil)).Elem(),
	"Statement":                        reflect.TypeOf((*semantic.Statement)(nil)).Elem(),
	"StatementDepth":                   reflect.TypeOf((*semantic.StatementDepth)(nil)).Elem(),
	"StatementExpression":              reflect.TypeOf((*semantic.StatementExpression)(nil)).Elem(),
	"Stmt":                             reflect.TypeOf((*semantic.Stmt)(nil)).Elem(),
	"StmtVisitor":                      reflect.TypeOf((*semantic.StmtVisitor)(nil)).Elem(),
	"StringLiteral":                    reflect.TypeOf((*semantic.StringLiteral)(nil)).Elem(),
	"StubExprVisitor":                  reflect.TypeOf((*semantic.StubExprVisitor)(nil)).Elem(),
	"StubStmtVisitor":                  reflect.TypeOf((*semantic.StubStmtVisitor)(nil)).Elem(),
	"TableRef":                         reflect.TypeOf((*semantic.TableRef)(nil)).Elem(),
	"UnaryLogicalExpression":           reflect.TypeOf((*semantic.UnaryLogicalExpression)(nil)).Elem(),
	"UpdateStatement":                  reflect.TypeOf((*semantic.UpdateStatement)(nil)).Elem(),
	"VariableDeclaration":              reflect.TypeOf((*semantic.VariableDeclaration)(nil)).Elem(),
	"WildCardField":                    reflect.TypeOf((*semantic.WildCardField)(nil)).Elem(),
}
