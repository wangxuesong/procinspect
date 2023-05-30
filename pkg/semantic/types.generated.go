// Code generated by gen-ast-types. DO NOT EDIT.

package semantic

import (
	"errors"
)

type ExprVisitor interface {
	VisitNumericLiteral(v *NumericLiteral) (result interface{}, err error)
	VisitNameExpression(v *NameExpression) (result interface{}, err error)
}

type StubExprVisitor struct{}

var _ ExprVisitor = StubExprVisitor{}

func (s StubExprVisitor) VisitNumericLiteral(_ *NumericLiteral) (interface{}, error) {
	return nil, errors.New("visit func for NumericLiteral is not implemented")
}

func (s StubExprVisitor) VisitNameExpression(_ *NameExpression) (interface{}, error) {
	return nil, errors.New("visit func for NameExpression is not implemented")
}

func (b *NumericLiteral) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitNumericLiteral(b)
}

func (b *NameExpression) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitNameExpression(b)
}

type StmtVisitor interface {
	VisitScript(v *Script) (err error)
	VisitCreateProcedureStatement(v *CreateProcedureStatement) (err error)
	VisitBlockStatement(v *BlockStatement) (err error)
	VisitBody(v *Body) (err error)
	VisitAssignmentStatement(v *AssignmentStatement) (err error)
	VisitProcedureCall(v *ProcedureCall) (err error)
	VisitVariableDeclaration(v *VariableDeclaration) (err error)
}

type StubStmtVisitor struct{}

var _ ExprVisitor = StubExprVisitor{}

func (s StubExprVisitor) VisitScript(_ *Script) error {
	return errors.New("visit func for Script is not implemented")
}

func (s StubExprVisitor) VisitCreateProcedureStatement(_ *CreateProcedureStatement) error {
	return errors.New("visit func for CreateProcedureStatement is not implemented")
}

func (s StubExprVisitor) VisitBlockStatement(_ *BlockStatement) error {
	return errors.New("visit func for BlockStatement is not implemented")
}

func (s StubExprVisitor) VisitBody(_ *Body) error {
	return errors.New("visit func for Body is not implemented")
}

func (s StubExprVisitor) VisitAssignmentStatement(_ *AssignmentStatement) error {
	return errors.New("visit func for AssignmentStatement is not implemented")
}

func (s StubExprVisitor) VisitProcedureCall(_ *ProcedureCall) error {
	return errors.New("visit func for ProcedureCall is not implemented")
}

func (s StubExprVisitor) VisitVariableDeclaration(_ *VariableDeclaration) error {
	return errors.New("visit func for VariableDeclaration is not implemented")
}

func (b *Script) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitScript(b)
}

func (b *CreateProcedureStatement) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitCreateProcedureStatement(b)
}

func (b *BlockStatement) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitBlockStatement(b)
}

func (b *Body) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitBody(b)
}

func (b *AssignmentStatement) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitAssignmentStatement(b)
}

func (b *ProcedureCall) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitProcedureCall(b)
}

func (b *VariableDeclaration) Accept(visitor StmtVisitor) (err error) {
	return visitor.VisitVariableDeclaration(b)
}
