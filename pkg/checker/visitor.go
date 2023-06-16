package checker

import (
	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	Visitor struct {
		semantic.StubExprVisitor
	}
)

func LoadScript(src string) (*semantic.Script, error) {
	script, err := parser.ParseScript(src)
	if err != nil {
		return nil, err
	}

	return script, nil
}
