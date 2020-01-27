package main

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/eval"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/parser"
	"Yum-Programming-Language-Interpreter/semantic"
	"fmt"
	"os"
)

func main() {
	var (
		l lexer.Lexer
		f *os.File
		p parser.Parser
		prog ast.Node
		sA semantic.SemanticAnalyser
		e eval.Evaluator
		err error
		errs []error
	)

	if f, err = os.Open("test_files/progressive.txt"); err != nil {
		fmt.Println(fmt.Sprintf(internal.ErrFailedToReadFile, "test_files/progressive.txt", err))
		os.Exit(0)
	}

	if l, err = lexer.NewLexer(f); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer l.Close()

	if p, err = parser.NewRecursiveDescentParser(l); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if prog, errs = p.Parse(); errs != nil && len(errs) != 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(0)
	}

	sA = semantic.NewSemanticAnalyser()

	if errs = sA.Analyse(prog); errs != nil && len(errs) != 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
	}

	e = eval.NewEvaluator()
	e.Evaluate(prog)


}
