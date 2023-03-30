package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
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
	el.err = fmt.Errorf("%d:%d: %s", line, column, msg)
	//panic(msg)
}

func (el *MyErrorListener) Error() error {
	return el.err
}
