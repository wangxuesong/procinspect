package parser

import (
	"procinspect/pkg/parser/internal/plsql/parser"
	"procinspect/pkg/semantic"
)

func Parse(text string) error {
	p := parser.NewParser(text)
	_ = p.Sql_script()
	return p.Error()
}

func ParseScript(src string) (*semantic.Script, error) {
	p := parser.NewParser(src)
	root := p.Sql_script()
	if p.Error() != nil {
		return nil, p.Error()
	}
	visitor := &plsqlVisitor{}
	script := visitor.VisitSql_script(root.(*parser.Sql_scriptContext)).(*semantic.Script)

	return script, nil
}

func ParseBlock(src string) (*semantic.Script, error) {
	p := parser.NewParser(src)
	root := p.Block()
	if p.Error() != nil {
		return nil, p.Error()
	}
	visitor := &plsqlVisitor{}
	block := visitor.VisitBlock(root.(*parser.BlockContext)).(*semantic.BlockStatement)

	script := &semantic.Script{}

	script.Statements = append(script.Statements, block)

	return script, nil
}
