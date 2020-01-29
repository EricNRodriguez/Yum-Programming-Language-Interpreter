package semantic

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/parser"
	"fmt"
	"github.com/spf13/afero"
	"testing"
)

func TestSemanticAnalyser(t *testing.T) {
	tCs := []struct {
		input   []byte
		numErrs int
	}{
		{
			[]byte("var x = 3; x = 2;"),
			0,
		},
		{
			[]byte("x = 3;"),
			1, // not declared
		},
		{
			[]byte("var x = 3; var x = 2;"),
			1, // already declared
		},
		{
			[]byte("print(y*22-x);"),
			2, // x and y not declared
		},
		{
			[]byte("var x = \"hello\";x = 230;"),
			0,
		},
		{
			[]byte("if (a < 3) { print(23 + 33);};"),
			1, // a not declared
		},
		{
			[]byte("if (2 < \"hello\") { print(23 + 33);};"),
			0, // no type checks until run time
		},
		{
			[]byte("var x = 33; if (2 < 3) { x = 000;};"),
			0, // x declared in outer scope
		},
		{
			[]byte("if (2 < 3) { x = 000;};"),
			1, // x not declared  in outer scope
		},
		{
			[]byte("var x = 33; if (2 < 3) { var x = 000;};"),
			0, // x not declared in inner scope
		},
		{
			[]byte("var x = 33; func testFunc(a,b) { x = 0001;};"),
			1, // x not declared inside function scope
		},
		{
			[]byte("if (true) { var x = 3; } else { print(x);};"),
			1, // x not declared in else block scope
		},
		{
			[]byte("var x = [1,2,3,4];"),
			0,
		},
		{
			[]byte("var x = [1,2,3,a];"),
			1, // a not declared
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[1];"),
			0,
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[length(x)-1];"),
			0,
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[\"word\"];"),
			1, // strings are not valid indexes
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[22.33];"),
			1, // floating point numbers are not valid indexes
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[true];"),
			1, // booleans are not valid indexes
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[a];"),
			1, // a not declared
		},
		{
			[]byte("var a = \"hello\"; var x = [1,2,3,4]; x = x[a];"),
			0, // a declared, type checks occur during runtime
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[1+2/3];"),
			0,
		},
		{
			[]byte("var x = [1,2,3,4]; var b = [1,2,3,4]; x = x[b[3]];"),
			0,
		},
		{
			[]byte("if (true) {func helloWorld() {print(\"hello world\");};print(\"made hello world\");};" +
				"helloWorld();"),
			0, // functions defined in global scope
		},
		{
			[]byte("func defineStuff() {var x = 1;var y = \"hello\";var z = [1,2,3,4]; var d = 12.33; " +
				"var f = true;}; print(x,y,z,d,f);"),
			5, // variables declared in function scope
		},
		{
			[]byte("func add(a,b,c,d) { return a + b + c + d;}; add(1,2,3);"),
			1, // invalid number of params
		},
		{
			[]byte("func add(a,b,c,d) { return a + b + c + d;}; func add(a,b,c) { return a + b + c;}; add(1,2,3);"),
			2, // add already declared and invalid number of params
		},
		{
			[]byte("if (true) { print(123); return 6;};"),
			1, // return statement outside of function
		},
		{
			[]byte("func hello() { if (true) { print(123); return 6;};};"),
			0, // return statement inside of function
		},
	}

	var (
		fs  afero.Fs
		err error
	)

	fs = afero.NewMemMapFs()
	if err = fs.MkdirAll("test_files/semantic_analysis", 0755); err != nil {
		t.Fatalf(err.Error())
	}

	for i, tC := range tCs {
		var (
			f    afero.File
			l    lexer.Lexer
			p    parser.Parser
			sA   *semanticAnalyser
			fp   string
			prog *ast.Program
			err  error
			errs []error
		)

		fp = fmt.Sprintf("test_files/semantic_analysis/test_%v.txt", i)

		if err = afero.WriteFile(fs, fp, tC.input, 0644); err != nil {
			t.Fatalf(err.Error())
		}

		if _, err = fs.Stat(fp); err != nil {
			t.Fatalf("file \"%s\" does not exist.\n", fp)
		}

		if f, err = fs.Open(fp); err != nil {
			t.Fatalf(err.Error())
		}

		if l, err = lexer.NewLexer(f); err != nil {
			t.Fatalf(err.Error())
		}

		if p, err = parser.NewRecursiveDescentParser(l); err != nil {
			t.Fatalf(err.Error())
		}

		if prog, errs = p.Parse(); errs != nil && len(errs) != 0 {
			t.Fatalf(internal.ErrInvalidSyntaxSemanticAnalysisTestCases, i+1, len(errs))
		}

		sA = NewSemanticAnalyser()

		errs = sA.Analyse(prog)

		if len(errs) != tC.numErrs {
			t.Errorf(internal.ErrInvalidNumberOfErrorsTest, i+1, tC.numErrs, len(errs))
		}

	}
	return
}
