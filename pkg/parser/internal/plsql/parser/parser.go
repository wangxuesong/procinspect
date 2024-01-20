package parser

import (
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type (
	SyntaxError struct {
		Line    int
		Column  int
		Stack   []string
		Message string
		Origin  string
	}
	MyErrorListener struct {
		antlr.DefaultErrorListener
		err []error
	}

	SqlParser struct {
		*PlSqlParser
		listener *MyErrorListener
	}
)

func NewParser(source string) *SqlParser {
	input := antlr.NewInputStream(source)
	lexer := NewPlSqlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewPlSqlParser(stream)
	parser.RemoveErrorListeners()
	listener := NewMyErrorListener()
	parser.AddErrorListener(listener)
	return &SqlParser{
		PlSqlParser: parser,
		listener:    listener,
	}
}

func (p *SqlParser) Error() error {
	return p.listener.Error()
}

func NewMyErrorListener() *MyErrorListener {
	return &MyErrorListener{err: nil}
}

func (el *MyErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	p := recognizer.(antlr.Parser)
	stack := p.GetRuleInvocationStack(p.GetParserRuleContext())
	el.err = append(el.err, SyntaxError{
		Line:    line,
		Column:  column,
		Stack:   stack,
		Message: msg,
		Origin:  fmt.Sprintf("%v", offendingSymbol),
	})
}

func (el *MyErrorListener) Error() error {
	var errs error
	errs = errors.Join(el.err...)
	// }
	return errs
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at line %d, column %d: %v", e.Line, e.Column, e.Stack)
}
