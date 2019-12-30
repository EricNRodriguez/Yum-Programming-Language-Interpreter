package main

import (
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"os"
)

func main() {
	var (
		l lexer.LexerInterface
	)
	f, _ := os.Open("test_files/progressive.txt")
	l, _ = lexer.NewLexer(f)

	fmt.Println("=============== ")
	tok, err := l.NextToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(tok.Type())
	for tok.Type() != token.EOF {
		tok, _ = l.NextToken()
		fmt.Println(tok.Type())
	}

}
