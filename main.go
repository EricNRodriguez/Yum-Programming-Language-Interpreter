package main

import (
	"Yum-Programming-Language-Interpreter/eval"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/parser"
	"fmt"
	"os"
)

func main() {
	var (
		l lexer.Lexer
	)
	f, _ := os.Open("test_files/progressive.txt")
	l, _ = lexer.NewLexer(f)
	defer l.Close()

	p, err := parser.NewRecursiveDescentParser(l)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prog := p.Parse()
	fmt.Println(prog.String())
	fmt.Println("----")
	fmt.Println(eval.Evaluate(prog).Literal())


}

