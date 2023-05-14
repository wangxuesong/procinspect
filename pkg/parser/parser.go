package parser

import "procinspect/pkg/parser/internal/plsql/parser"

func Parse(text string) error {
	p := parser.NewParser(text)
	_ = p.Sql_script()
	return p.Error()
}
