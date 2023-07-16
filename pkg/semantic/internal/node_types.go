package main

var nodeTypes = Types{
	{
		Name:    "AliasExpression",
		Fields:  "semantic.AliasExpression",
		Comment: "",
	},
	{
		Name:    "Argument",
		Fields:  "semantic.Argument",
		Comment: "",
	},
	{
		Name:    "AssignmentStatement",
		Fields:  "semantic.AssignmentStatement",
		Comment: "",
	},
	{
		Name:    "AutonomousTransactionDeclaration",
		Fields:  "semantic.AutonomousTransactionDeclaration",
		Comment: "",
	},
	{
		Name:    "BetweenExpression",
		Fields:  "semantic.BetweenExpression",
		Comment: "",
	},
	{
		Name:    "BinaryExpression",
		Fields:  "semantic.BinaryExpression",
		Comment: "",
	},
	{
		Name:    "BindNameExpression",
		Fields:  "semantic.BindNameExpression",
		Comment: "",
	},
	{
		Name:    "BlockStatement",
		Fields:  "semantic.BlockStatement",
		Comment: "",
	},
	{
		Name:    "Body",
		Fields:  "semantic.Body",
		Comment: "",
	},
	{
		Name:    "CaseWhenBlock",
		Fields:  "semantic.CaseWhenBlock",
		Comment: "",
	},
	{
		Name:    "CaseWhenStatement",
		Fields:  "semantic.CaseWhenStatement",
		Comment: "",
	},
	{
		Name:    "CastExpression",
		Fields:  "semantic.CastExpression",
		Comment: "",
	},
	{
		Name:    "CloseStatement",
		Fields:  "semantic.CloseStatement",
		Comment: "",
	},
	{
		Name:    "CommitStatement",
		Fields:  "semantic.CommitStatement",
		Comment: "",
	},
	{
		Name:    "ContinueStatement",
		Fields:  "semantic.ContinueStatement",
		Comment: "",
	},
	{
		Name:    "CreateFunctionStatement",
		Fields:  "semantic.CreateFunctionStatement",
		Comment: "",
	},
	{
		Name:    "CreateNestTableStatement",
		Fields:  "semantic.CreateNestTableStatement",
		Comment: "",
	},
	{
		Name:    "CreatePackageBodyStatement",
		Fields:  "semantic.CreatePackageBodyStatement",
		Comment: "",
	},
	{
		Name:    "CreatePackageStatement",
		Fields:  "semantic.CreatePackageStatement",
		Comment: "",
	},
	{
		Name:    "CreateProcedureStatement",
		Fields:  "semantic.CreateProcedureStatement",
		Comment: "",
	},
	{
		Name:    "CreateSynonymStatement",
		Fields:  "semantic.CreateSynonymStatement",
		Comment: "",
	},
	{
		Name:    "CreateTypeStatement",
		Fields:  "semantic.CreateTypeStatement",
		Comment: "",
	},
	{
		Name:    "CursorAttribute",
		Fields:  "semantic.CursorAttribute",
		Comment: "",
	},
	{
		Name:    "CursorDeclaration",
		Fields:  "semantic.CursorDeclaration",
		Comment: "",
	},
	{
		Name:    "DeleteStatement",
		Fields:  "semantic.DeleteStatement",
		Comment: "",
	},
	{
		Name:    "DotExpression",
		Fields:  "semantic.DotExpression",
		Comment: "",
	},
	{
		Name:    "ElseBlock",
		Fields:  "semantic.ElseBlock",
		Comment: "",
	},
	{
		Name:    "ExceptionDeclaration",
		Fields:  "semantic.ExceptionDeclaration",
		Comment: "",
	},
	{
		Name:    "ExecuteImmediateStatement",
		Fields:  "semantic.ExecuteImmediateStatement",
		Comment: "",
	},
	{
		Name:    "ExistsExpression",
		Fields:  "semantic.ExistsExpression",
		Comment: "",
	},
	{
		Name:    "ExitStatement",
		Fields:  "semantic.ExitStatement",
		Comment: "",
	},
	{
		Name:    "FetchStatement",
		Fields:  "semantic.FetchStatement",
		Comment: "",
	},
	{
		Name:    "FieldList",
		Fields:  "semantic.FieldList",
		Comment: "",
	},
	{
		Name:    "ForUpdateClause",
		Fields:  "semantic.ForUpdateClause",
		Comment: "",
	},
	{
		Name:    "ForUpdateOptionsExpression",
		Fields:  "semantic.ForUpdateOptionsExpression",
		Comment: "",
	},
	{
		Name:    "FromClause",
		Fields:  "semantic.FromClause",
		Comment: "",
	},
	{
		Name:    "FunctionCallExpression",
		Fields:  "semantic.FunctionCallExpression",
		Comment: "",
	},
	{
		Name:    "FunctionDeclaration",
		Fields:  "semantic.FunctionDeclaration",
		Comment: "",
	},
	{
		Name:    "IfStatement",
		Fields:  "semantic.IfStatement",
		Comment: "",
	},
	{
		Name:    "InExpression",
		Fields:  "semantic.InExpression",
		Comment: "",
	},
	{
		Name:    "InsertStatement",
		Fields:  "semantic.InsertStatement",
		Comment: "",
	},
	{
		Name:    "IntoClause",
		Fields:  "semantic.IntoClause",
		Comment: "",
	},
	{
		Name:    "LikeExpression",
		Fields:  "semantic.LikeExpression",
		Comment: "",
	},
	{
		Name:    "LoopStatement",
		Fields:  "semantic.LoopStatement",
		Comment: "",
	},
	{
		Name:    "NameExpression",
		Fields:  "semantic.NameExpression",
		Comment: "",
	},
	{
		Name:    "NestTableTypeDeclaration",
		Fields:  "semantic.NestTableTypeDeclaration",
		Comment: "",
	},
	{
		Name:    "NullExpression",
		Fields:  "semantic.NullExpression",
		Comment: "",
	},
	{
		Name:    "NullStatement",
		Fields:  "semantic.NullStatement",
		Comment: "",
	},
	{
		Name:    "NumericLiteral",
		Fields:  "semantic.NumericLiteral",
		Comment: "",
	},
	{
		Name:    "OpenStatement",
		Fields:  "semantic.OpenStatement",
		Comment: "",
	},
	{
		Name:    "OuterJoinExpression",
		Fields:  "semantic.OuterJoinExpression",
		Comment: "",
	},
	{
		Name:    "Parameter",
		Fields:  "semantic.Parameter",
		Comment: "",
	},
	{
		Name:    "ProcedureCall",
		Fields:  "semantic.ProcedureCall",
		Comment: "",
	},
	{
		Name:    "QueryExpression",
		Fields:  "semantic.QueryExpression",
		Comment: "",
	},
	{
		Name:    "RelationalExpression",
		Fields:  "semantic.RelationalExpression",
		Comment: "",
	},
	{
		Name:    "ReturnStatement",
		Fields:  "semantic.ReturnStatement",
		Comment: "",
	},
	{
		Name:    "RollbackStatement",
		Fields:  "semantic.RollbackStatement",
		Comment: "",
	},
	{
		Name:    "Script",
		Fields:  "semantic.Script",
		Comment: "",
	},
	{
		Name:    "SelectField",
		Fields:  "semantic.SelectField",
		Comment: "",
	},
	{
		Name:    "SelectStatement",
		Fields:  "semantic.SelectStatement",
		Comment: "",
	},
	{
		Name:    "SignExpression",
		Fields:  "semantic.SignExpression",
		Comment: "",
	},
	{
		Name:    "StatementExpression",
		Fields:  "semantic.StatementExpression",
		Comment: "",
	},
	{
		Name:    "StringLiteral",
		Fields:  "semantic.StringLiteral",
		Comment: "",
	},
	{
		Name:    "TableRef",
		Fields:  "semantic.TableRef",
		Comment: "",
	},
	{
		Name:    "UnaryLogicalExpression",
		Fields:  "semantic.UnaryLogicalExpression",
		Comment: "",
	},
	{
		Name:    "UpdateStatement",
		Fields:  "semantic.UpdateStatement",
		Comment: "",
	},
	{
		Name:    "VariableDeclaration",
		Fields:  "semantic.VariableDeclaration",
		Comment: "",
	},
	{
		Name:    "WildCardField",
		Fields:  "semantic.WildCardField",
		Comment: "",
	},
}