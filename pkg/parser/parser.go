package parser

import (
	"procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
)

func Parse(text string) (any, error) {
	p := parser.NewParser(text)
	root := p.Sql_script()
	return root, p.Error()
}

func ParseScript(src string) (*semantic.Script, error) {
	p := parser.NewParser(src)
	root := p.Sql_script()
	if p.Error() != nil {
		return nil, p.Error()
	}
	visitor := newPlSqlVisitor()
	script := visitor.VisitSql_script(root.(*parser.Sql_scriptContext)).(*semantic.Script)

	return script, nil
}

func ParseSql(src string) (func() (*semantic.Script, error), error) {
	p := parser.NewParser(src)
	root := p.Sql_script()
	if p.Error() != nil {
		return nil, p.Error()
	}

	return func() (*semantic.Script, error) {
		visitor := newPlSqlVisitor()
		script := visitor.VisitSql_script(root.(*parser.Sql_scriptContext)).(*semantic.Script)

		return script, nil
	}, nil
}

func ParseBlock(src string) (*semantic.Script, error) {
	p := parser.NewParser(src)
	root := p.Block()
	if p.Error() != nil {
		return nil, p.Error()
	}
	visitor := newPlSqlVisitor()
	block := visitor.VisitBlock(root.(*parser.BlockContext)).(*semantic.BlockStatement)

	script := &semantic.Script{}

	script.Statements = append(script.Statements, block)

	return script, nil
}
