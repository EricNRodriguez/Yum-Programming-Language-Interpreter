package eval

import (
	"github.com/EricNRodriguez/yum/ast"
	"github.com/EricNRodriguez/yum/internal"
	"github.com/EricNRodriguez/yum/object"
	"github.com/EricNRodriguez/yum/symbol_table"
	"github.com/EricNRodriguez/yum/token"
	"github.com/EricNRodriguez/yum/semantic"

	"fmt"
	"github.com/spf13/afero"
	"testing"
)

type symbol struct {
	iden  string
	value string
}

type testCase struct {
	input        []byte
	err          bool
	symbolValues []symbol
}

func TestEvaluator(t *testing.T) {

	tCs := []testCase{
		{
			[]byte("var x = 3;"),
			false,
			[]symbol{
				{
					"x",
					"3",
				},
			},
		},
		{
			[]byte("var x = 3; x = \"hello\";"),
			false,
			[]symbol{
				{
					"x",
					"\"hello\"",
				},
			},
		},
		{
			[]byte("var x = 3; var y = x; x = \"hello\";"),
			false,
			[]symbol{
				{
					"x",
					"\"hello\"",
				},
				{
					"y",
					"3",
				},
			},
		},
		{
			[]byte("var x = 3; if (true) { x = 22; };"),
			false,
			[]symbol{
				{
					"x",
					"22",
				},
			},
		},
		{
			[]byte("var x = 3; if (true) { var x = 22; };"),
			false,
			[]symbol{
				{
					"x",
					"3",
				},
			},
		},
		{
			[]byte("var x = [1,2,3,4]; var y = x[0]; var z = x[length(x)-1];"),
			false,
			[]symbol{
				{
					"x",
					"[1,2,3,4]",
				},
				{
					"y",
					"1",
				},
				{
					"z",
					"4",
				},
			},
		},
		{
			[]byte("var a = +2; var b = +-2; var c = -+2; var d = --+--2; var e = -02; var f = -020;"),
			false,
			[]symbol{
				{
					"a",
					"2",
				},
				{
					"b",
					"-2",
				},
				{
					"c",
					"-2",
				},
				{
					"d",
					"2",
				},
				{
					"e",
					"-2",
				},
				{
					"f",
					"-20",
				},
			},
		},
		{
			[]byte("var a = -1 + 2 + 3 + 4; var b = 1 + (-2 + 3) + 4; var c = 1 + 2 + (3 + 4);"),
			false,
			[]symbol{
				{
					"a",
					"8",
				},
				{
					"b",
					"6",
				},
				{
					"c",
					"10",
				},
			},
		},
		{
			[]byte("var a = 1 * 2 * -3 * 4; var b = 1 * (-2 * -3) * 4; var c = 1 * 2 * (3 * 4); " +
				"var d = 1 * (-(3*4)*5);"),
			false,
			[]symbol{
				{
					"a",
					"-24",
				},
				{
					"b",
					"24",
				},
				{
					"c",
					"24",
				},
				{
					"d",
					"-60",
				},
			},
		},
		{
			[]byte("var a = 1 - -2 - 3 - 4; var b = 1 - -(2 - 3) - 4; var c = 1 - 2 - (3 - 4);"),
			false,
			[]symbol{
				{
					"a",
					"-4",
				},
				{
					"b",
					"-4",
				},
				{
					"c",
					"0",
				},
			},
		},
		{
			[]byte("var a = 1 / -2; var b = (2 / -3) / 4;"),
			false,
			[]symbol{
				{
					"a",
					"0",
				},
				{
					"b",
					"0",
				},
			},
		},
		{
			[]byte("var a = 1.0 / -2 / 3 / 4; var b = 1 / (2.0 / -3) / 4; var c = 1.0 / 2 / (-3.0 / 4);"),
			false,
			[]symbol{
				{
					"a",
					"-0.041667",
				},
				{
					"b",
					"-0.375000",
				},
				{
					"c",
					"-0.666667",
				},
			},
		},
		{
			[]byte("var a = 2 + -(3 * -4) / 9.0 ; var b = 2 - -3 / -(4.0* 2); var c = 2 * 3.0 - -4 / -2.0;"),
			false,
			[]symbol{
				{
					"a",
					"3.333333",
				},
				{
					"b",
					"1.625000",
				},
				{
					"c",
					"4.000000",
				},
			},
		},
		{
			[]byte("var a = \"hello\"; var b = \"world\"; var c = a + b; var d = a + \" \" + b;"),
			false,
			[]symbol{
				{
					"a",
					"\"hello\"",
				},
				{
					"b",
					"\"world\"",
				},
				{
					"c",
					"\"helloworld\"",
				},
				{
					"d",
					"\"hello world\"",
				},
			},
		},
		{
			[]byte("var a = !true; var b = !!true; var c = !(true | false); var d = !!(!(false));"),
			false,
			[]symbol{
				{
					"a",
					"false",
				},
				{
					"b",
					"true",
				},
				{
					"c",
					"false",
				},
				{
					"d",
					"true",
				},
			},
		},
		{
			[]byte("var a = true | false | true; var b = true | !(false | true);"),
			false,
			[]symbol{
				{
					"a",
					"true",
				},
				{
					"b",
					"true",
				},
			},
		},
		{
			[]byte("var a = true & !false & true; var b = (!true | !true) & true;"),
			false,
			[]symbol{
				{
					"a",
					"true",
				},
				{
					"b",
					"false",
				},
			},
		},
		{
			[]byte("var a = true & !false & true; var b = (!true | !true) & true;"),
			false,
			[]symbol{
				{
					"a",
					"true",
				},
				{
					"b",
					"false",
				},
			},
		},
		{
			[]byte("var a = true & !(false | (true | false) & true); var b = !false & (false | true);"),
			false,
			[]symbol{
				{
					"a",
					"false",
				},
				{
					"b",
					"true",
				},
			},
		},
		{
			[]byte("func hello() {return 1 + 2;};var x = [1,2,3,4,5];var a = x[hello()];var b = hello() * 2;" +
				"var c = hello() + hello();"),
			false,
			[]symbol{
				{
					"a",
					"4",
				},
				{
					"b",
					"6",
				},
				{
					"c",
					"6",
				},
			},
		},
		{
			[]byte("func fact(n) {if (n == 1) {return 1;};return n * fact(n-1);};var a = fact(5);"),
			false,
			[]symbol{
				{
					"a",
					"120",
				},
			},
		},
		{
			[]byte("var x = 22;var y = 33;func getNum(n) {return true;};if (true) {x = 100;y = y;} " +
				"else {x = 200;y = getNum(x);};"),
			false,
			[]symbol{
				{
					"x",
					"100",
				},
				{
					"y",
					"33",
				},
			},
		},
		{
			[]byte("var x = 22;var y = 33;func getNum(n) {return true;};if (false) {x = 100;y = y;} " +
				"else {x = 200;y = getNum(x);};"),
			false,
			[]symbol{
				{
					"x",
					"200",
				},
				{
					"y",
					"true",
				},
			},
		},
		{
			[]byte("var x = [1,2,3,4,5,6,7,8,9]; var i = 0; " +
				"var a = -1; while (i < length(x) - 4) {a = x[i]; i = i + 1;};"),
			false,
			[]symbol{
				{
					"a",
					"5",
				},
			},
		},
		{
			[]byte("var x = [1,2,3,4,5,6,7,8,9]; var i = 0; " +
				"var a = -1; while (i < length(x) - 4) {var x  = x[i]; i = i + 1;};"),
			false,
			[]symbol{
				{
					"a",
					"-1",
				},
			},
		},
		{
			[]byte("func returnNull() {}; var x = returnNull();"),
			false,
			[]symbol{
				{
					"x",
					"null",
				},
			},
		},
		{
			[]byte("if (2) {print(2);};"),
			true, // condition not boolean
			[]symbol{},
		},
		{
			[]byte("var x = 10/0;"),
			true, // division by 0
			[]symbol{},
		},
		{
			[]byte("var x = true/true;"),
			true, // invalid infix op for int and bool
			[]symbol{},
		},
		{
			[]byte("var x = 10*true;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 10.0+false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\"-false;"),
			true, // invalid infix
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\"-\"hello\";"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\"*\"hello\";"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\"/\"hello\";"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3]-[1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3]*[1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3]/[1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3]+[1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true + false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true +- false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true / false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true * false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true > false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true >= false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true < false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true <= false;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true + 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true / 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true * 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true - 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" + 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" * 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" / 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" - 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" + true;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" - 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" * 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" / 1;"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" + [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" - [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" * [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = \"hello\" / [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 1 + [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 1 - [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 1 * [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 1 / [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true + [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true - [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true * [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = true / [1,2,3];"),
			true, // invalid infix op
			[]symbol{},
		},
		{
			[]byte("var x = 1 < 2 < 3;"),
			true, // invalid infix op for bool and int
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[-1];"),
			true, // index out of bounds
			[]symbol{},
		},
		{
			[]byte("var x = [1,2,3,4]; x = x[10];"),
			true, // index out of bounds
			[]symbol{},
		},
		{
			[]byte("while (2 + 3 + 4) {};"),
			true, // condition doesnt eval to boolean
			[]symbol{},
		},
		{
			[]byte("func getStr() {return \"hello\";}; var x = [1,2,3,4]; x = x[getStr()];"),
			true, // index type
			[]symbol{},
		},
		{
			[]byte("func getFloat() {return 22.33;}; var x = [1,2,3,4]; x = x[getFloat()];"),
			true, // index type
			[]symbol{},
		},
		{
			[]byte("func getBool() {return true;}; var x = [1,2,3,4]; x = x[getBool()];"),
			true, // index type
			[]symbol{},
		},
		{
			[]byte("var x = -\"hello\";"),
			true, // invalid prefix op
			[]symbol{},
		},
		{
			[]byte("var x = +\"hello\";"),
			true, // invalid prefix op
			[]symbol{},
		},
		{
			[]byte("var x = -[1,2,3];"),
			true, // invalid prefix op
			[]symbol{},
		},
		{
			[]byte("var x = +[1,2,3];"),
			true, // invalid prefix op
			[]symbol{},
		},
		{
			[]byte("var x = -true;"),
			true, // invalid prefix op
			[]symbol{},
		},
		{
			[]byte("var x = +true;"),
			true, // invalid prefix op
			[]symbol{},
		},
	}

	var (
		fs  afero.Fs
		err error
	)

	fs = afero.NewMemMapFs()
	if err = fs.MkdirAll("test_files/evaluation", 0755); err != nil {
		t.Fatalf(err.Error())
	}

	for i, tC := range tCs {
		func() {
			defer func() {
				if r := recover(); r != nil && !tC.err {
					t.Errorf(internal.ErrUnexpectedRuntimeError, i+1, r)
				}
				return
			}()

			var (
				f    afero.File
				l    lexer.Lexer
				p    parser.Parser
				sA   semantic.SemanticAnalyser
				fp   string
				prog *ast.Program
				e    *evaluator
				err  error
				errs []error
			)

			fp = fmt.Sprintf("test_files/evaluation/test_%v.txt", i)

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
				t.Fatalf(internal.ErrInvalidSyntaxEvaluationTestCases, i+1, len(errs))
			}

			sA = semantic.NewSemanticAnalyser()

			if errs = sA.Analyse(prog); errs != nil && len(errs) != 0 {
				t.Fatalf(internal.ErrInvalidSemanticsEvaluationTestCases, i+1, len(errs))
			}

			e = NewEvaluator()
			e.evaluate(prog)

			// test variables
			for _, expectedSymbol := range tC.symbolValues {
				var (
					trueValue object.Object
					ok        bool
				)

				if trueValue, ok = e.symbolTable.GetVar(expectedSymbol.iden); !ok {
					t.Errorf(internal.ErrMissingSymbolTest, i+1, expectedSymbol.iden, expectedSymbol.iden, expectedSymbol.value)
					continue
				}

				if trueValue.Literal() != expectedSymbol.value {
					t.Errorf(internal.ErrInvalidSymbolValueTest, i+1, expectedSymbol.iden, expectedSymbol.value,
						expectedSymbol.iden, trueValue.Literal())
					continue
				}
			}

		}()

	}
	return
}
