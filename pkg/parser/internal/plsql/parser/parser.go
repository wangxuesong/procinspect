package parser

import (
	"errors"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
)

type (
	MyErrorListener struct {
		antlr.DefaultErrorListener
		err error
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
	el.err = errors.Join(el.err, fmt.Errorf("stack: %v; %d:%d at %v: %s", stack[0], line, column, offendingSymbol, ""))
}

func (el *MyErrorListener) Error() error {
	return el.err
}
