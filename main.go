package main

import (
	"Yum-Programming-Language-Interpreter/eval"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/parser"
	"Yum-Programming-Language-Interpreter/semantic"
	"fmt"
	"os"
)

func main() {
	var (
		l lexer.Lexer
	)
	f, _ := os.Open("test_files/progressive.txt")
	l, err := lexer.NewLexer(f)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	p, err := parser.NewRecursiveDescentParser(l)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prog := p.Parse()

	sA := semantic.NewSemanticAnalyser()
	sA.Analyse(prog)
	for _, e := range sA.SemanticErrors() {
		fmt.Println(e)
	}

	//fmt.Println(prog.String())
	evalu := eval.NewEvaluator()
	if len(sA.SemanticErrors()) == 0 {
		evalu.Evaluate(prog)
	}

}
