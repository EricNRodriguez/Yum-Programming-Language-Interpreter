package main

import (
	"github.com/EricNRodriguez/yum/ast"
	"github.com/EricNRodriguez/yum/eval"
	"github.com/EricNRodriguez/yum/internal"
	"github.com/EricNRodriguez/yum/lexer"
	"github.com/EricNRodriguez/yum/parser"
	"github.com/EricNRodriguez/yum/semantic"
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
)

func main() {
	var (
		l     lexer.Lexer
		appFs afero.Fs
		fp    string
		f     afero.File
		p     parser.Parser
		prog  ast.Node
		sA    semantic.SemanticAnalyser
		e     eval.Evaluator
		err   error
		errs  []error
	)

	if len(os.Args[1:]) == 0 {
		fmt.Println(internal.ErrFileNotProvided)
		os.Exit(0)
	}

	appFs = afero.NewOsFs()

	fp = os.Args[1:][0]
	if _, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(internal.ErrFileNotFound+"\n", fp)
		} else {
			fmt.Printf(internal.ErrLoadFile+"\n", fp, err.Error())
		}
		os.Exit(0)
	}

	if f, err = appFs.Open(fp); err != nil {
		log.Println(fmt.Sprintf(internal.ErrFailedToReadFile+"\n", fp, err))
		os.Exit(0)
	}

	if l, err = lexer.NewLexer(f); err != nil {
		log.Println(err)
		os.Exit(0)
	}

	defer l.Close()

	if p, err = parser.NewRecursiveDescentParser(l); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if prog, errs = p.Parse(); errs != nil && len(errs) != 0 {
		for _, e := range errs {
			log.Println(e)
		}
		os.Exit(0)
	}

	sA = semantic.NewSemanticAnalyser()

	if errs = sA.Analyse(prog); errs != nil && len(errs) != 0 {
		for _, e := range errs {
			log.Println(e)
		}
		os.Exit(0)
	}

	e = eval.NewEvaluator()
	e.Evaluate(prog)

	return
}
