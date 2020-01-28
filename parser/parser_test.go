package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"strings"
	"fmt"
	"github.com/spf13/afero"
	"testing"
)

func TestParser(t *testing.T) {
	tCs := []struct{
		input []byte
		outputNodeTypes []ast.NodeType
		numErrs int
		outputLiteral string
	}{
		{
			[]byte("var x = 3;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = 3;",
		},
		{
			[]byte("x = 3;"),
			[]ast.NodeType{ast.AssignmentStatementNode},
			0,
			"x = 3;",
		},
		{
			[]byte("return 100;"),
			[]ast.NodeType{ast.ReturnStatementNode},
			0,
			"return 100;",
		},
		{
			[]byte("return;"),
			[]ast.NodeType{ast.ReturnStatementNode},
			0,
			"return;",
		},
		{
			[]byte("print(100*22-3);"),
			[]ast.NodeType{ast.FunctionCallStatementNode},
			0,
			"print(((100*22)-3));",
		},
		{
			[]byte("sayHello();"),
			[]ast.NodeType{ast.FunctionCallStatementNode},
			0,
			"sayHello();",
		},
		{
			[]byte("if (true & false | 3 < 2) {print(33);};"),
			[]ast.NodeType{ast.IfStatementNode},
			0,
			"if((true&false)|(3<2)){print(33);};",
		},
		{
			[]byte("if (3 < 2) {print(33);} else {print(\"howdy\");};"),
			[]ast.NodeType{ast.IfStatementNode},
			0,
			"if (3 < 2) {print(33);} else {print(\"howdy\");};",
		},
		{
			[]byte("while (3 < 2) {print(3 & 22);};"),
			[]ast.NodeType{ast.WhileStatementNode},
			0,
			"while(3<2){print((3 & 22));};",
		},
		{
			[]byte("while (true & false) {var x = 3;};"),
			[]ast.NodeType{ast.WhileStatementNode},
			0,
			"while (true & false) {var x = 3;};",
		},
		{
			[]byte("func add(a,b,c) {return a + b + c;};"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			0,
			"func add(a,b,c) {return ((a + b) + c);};",
		},
		{
			[]byte("var x == 3;"),
			[]ast.NodeType{ast.VarStatementNode},
			1, // expected assign token, received =
			"",
		},
		{
			[]byte("x = 3"),
			[]ast.NodeType{ast.AssignmentStatementNode},
			1, // expected ;
			"",
		},
		{
			[]byte("return var x == 2;;"),
			[]ast.NodeType{ast.ReturnStatementNode},
			2, // var not prefix operator, invalid ;
			"",
		},
		{
			[]byte("return"),
			[]ast.NodeType{ast.ReturnStatementNode},
			2, // EOF not valid prefix op and missing ;,
			"",
		},
		{
			[]byte("print(100**22-3);"),
			[]ast.NodeType{ast.FunctionCallStatementNode},
			1, // ** is not valid in an expression
			"",
		},
		{
			[]byte("sayHello);"),
			[]ast.NodeType{ast.FunctionCallStatementNode},
			1, // missing (
			"",
		},
		{
			[]byte("if (true & false | 3 % 2 == 0) {print(33);};"),
			[]ast.NodeType{ast.IfStatementNode},
			1, // % not a valid infix op, expected )
			"",
		},
		{
			[]byte("if (3 < 2 {print(33);} else2 {print(\"howdy\");};"),
			[]ast.NodeType{ast.IfStatementNode},
			2, // missing ) and expected ; instead of else 2
			"",
		},
		{
			[]byte("func add(a,b,3) {return a + b + c;};"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			1, // params can only be identifiers, the 3 is not valid syntax
			"func add(a,b,c) {return ((a + b) + c);};",
		},
		{
			[]byte("func add(a,b,c) {return a + b + c};"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			1, // missing semicolon on return statement
			"",
		},
		{
			[]byte("func add(a,b,c,) {return a + b + c;};"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			1, // invalid comma in function sig
			"",
		},
		{
			[]byte("def add(a,b,c,) {return a + b + c;};"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			1, // invalid statement beginning with def, unable to parse
			"",
		},
		{
			[]byte("var x = a + b * c - d;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((a + (b * c)) - d);",

		},
		{
			[]byte("var x = a - b / c - d;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((a - (b / c)) - d);",

		},
		{
			[]byte("var x = a * b / c - d;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = (((a * b) / c) - d);",

		},
		{
			[]byte("var x = a * b / c * d;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = (((a * b) / c) * d);",

		},
		{
			[]byte("var x = false & true & false;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((false & true) & false);",

		},
		{
			[]byte("var x = false & true | false;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((false & true) | false);",

		},
		{
			[]byte("var x = false | true & 3 <= 5 | false;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((false | (true & (3 <= 5))) | false);",

		},
		{
			[]byte("var x = 1 < 2 & false;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((1 < 2) & false);",

		},
		{
			[]byte("var x = 1 < 2 & 3 >= 4;"),
			[]ast.NodeType{ast.VarStatementNode},
			0,
			"var x = ((1 < 2) & (3 >= 4));",

		},
		{
			[]byte("func (2 < 3) { var x = 000;}; print(x);"),
			[]ast.NodeType{ast.FunctionDeclarationStatementNode},
			1, // function name not given
			"",
		},

	}

	var (
		fs  afero.Fs
		err error
	)

	fs = afero.NewMemMapFs()
	if err = fs.MkdirAll("test_files/parser_test", 0755); err != nil {
		t.Fatalf(err.Error())
	}

	for i, tC := range tCs {
		var (
			f   afero.File
			l   lexer.Lexer
			p Parser
			fp  string
			prog *ast.Program
			err error
			errs []error
		)

		fp = fmt.Sprintf("test_files/parser_test/test_%v.txt", i)

		if err = afero.WriteFile(fs, fp, tC.input, 0644); err != nil {
			t.Fatalf(err.Error())
		}

		if _, err = fs.Stat(fp); err != nil {
			t.Errorf("file \"%s\" does not exist.\n", fp)
		}

		if f, err = fs.Open(fp); err != nil {
			t.Errorf(err.Error())
		}

		if l, err = lexer.NewLexer(f); err != nil {
			t.Errorf(err.Error())
		}

		if p, err = NewRecursiveDescentParser(l); err != nil {
			t.Errorf(err.Error())
		}

		prog, errs = p.Parse()

		if len(errs) != tC.numErrs {
			t.Errorf(internal.ErrInvalidNumberOfErrorsTest, i+1, tC.numErrs, len(errs))
			return
		}

		if tC.numErrs != 0 {
			continue
		}

		if len(tC.outputNodeTypes) != len(prog.Statements) {
			t.Errorf(internal.ErrInvalidNumberOfASTNodesTest, i+1, len(tC.outputNodeTypes), len(prog.Statements))
			return
		}

		for j := 0; j < len(tC.outputNodeTypes); j++ {

			if prog.Statements[j].Type() != tC.outputNodeTypes[j] {
				t.Errorf(internal.ErrInvalidASTNodeTypeTest, i+1, prog.Statements[j].Type(), tC.outputNodeTypes[j])
				return
			}

			if prog.Statements[j].Type() != tC.outputNodeTypes[j] {
				t.Errorf(internal.ErrInvalidASTNodeTypeTest, i+1, prog.Statements[j].Type(), tC.outputNodeTypes[j])
				return
			}
		}

		expectedStr := tC.outputLiteral
		expectedStr = strings.Replace(expectedStr, " ", "", -1)
		expectedStr = strings.Replace(expectedStr, "\n", "", -1)

		parsedStr := prog.String()
		parsedStr = strings.Replace(parsedStr, " ", "", -1)
		parsedStr = strings.Replace(parsedStr, "\n", "", -1)


		if expectedStr != parsedStr {
			t.Errorf(internal.ErrInvalidASTNodeLiteralTest, i+1, parsedStr, expectedStr)
			return
		}


	}
	return
}
