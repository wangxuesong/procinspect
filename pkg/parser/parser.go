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
