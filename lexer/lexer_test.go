package lexer

import (
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	tC := struct {
		inputFilePath  string
		outputTypes    []token.TokenType
		outputLiterals []string
	}{
			"../test_files/lexer_tests.txt",
			[]token.TokenType{
				//var x = 22;
				token.VarToken, token.IdentifierToken, token.AssignToken, token.IntegerToken, token.SemicolonToken,

				//var u = (a + b) / 22;
				token.VarToken, token.IdentifierToken, token.AssignToken, token.LeftParenToken, token.IdentifierToken,
				token.AddToken, token.IdentifierToken, token.RightParenToken, token.DivToken, token.IntegerToken,
				token.SemicolonToken,

				//var hello = [a,2,"hello",4];
				token.VarToken, token.IdentifierToken, token.AssignToken, token.LeftBracketToken, token.IdentifierToken,
				token.CommaToken, token.IntegerToken, token.CommaToken, token.QuotationMarkToken, token.IdentifierToken,
				token.QuotationMarkToken, token.CommaToken, token.IntegerToken, token.RightBracketToken, token.SemicolonToken,

				// x = 100 / [a,2,"hello",22.33];
				token.IdentifierToken, token.AssignToken, token.IntegerToken, token.DivToken, token.LeftBracketToken,
				token.IdentifierToken, token.CommaToken, token.IntegerToken, token.CommaToken, token.QuotationMarkToken,
				token.IdentifierToken, token.QuotationMarkToken, token.CommaToken, token.SubToken,
				token.FloatingPointToken, token.RightBracketToken, token.SemicolonToken,

				// func testFunc(a,b,c,d,e) { print(22); d = a * b + c; return a + b - c / d * e;};
				token.FuncToken, token.IdentifierToken, token.LeftParenToken, token.IdentifierToken, token.CommaToken,
				token.IdentifierToken, token.CommaToken, token.IdentifierToken, token.CommaToken, token.IdentifierToken,
				token.CommaToken, token.IdentifierToken, token.RightParenToken, token.LeftBraceToken, token.IdentifierToken,
				token.LeftParenToken, token.IntegerToken, token.RightParenToken, token.SemicolonToken, token.IdentifierToken,
				token.AssignToken, token.IdentifierToken, token.MultToken, token.IdentifierToken, token.AddToken,
				token.IdentifierToken, token.SemicolonToken, token.ReturnToken, token.IdentifierToken, token.AddToken,
				token.IdentifierToken, token.SubToken, token.IdentifierToken, token.DivToken, token.IdentifierToken,
				token.MultToken, token.IdentifierToken, token.SemicolonToken, token.RightBraceToken, token.SemicolonToken,

				// if (true & false | false) { testFunc(1,2,3,4,5);} else { testFunc(5,4,3,2,1);};
				token.IfToken, token.LeftParenToken, token.BooleanToken, token.AndToken, token.BooleanToken, token.OrToken,
				token.BooleanToken, token.RightParenToken, token.LeftBraceToken, token.IdentifierToken, token.LeftParenToken,
				token.IntegerToken, token.CommaToken, token.IntegerToken, token.CommaToken, token.IntegerToken,
				token.CommaToken, token.IntegerToken, token.CommaToken, token.IntegerToken, token.RightParenToken,
				token.SemicolonToken, token.RightBraceToken, token.ElseToken, token.LeftBraceToken, token.IdentifierToken,
				token.LeftParenToken, token.IntegerToken, token.CommaToken, token.IntegerToken, token.CommaToken,
				token.IntegerToken, token.CommaToken, token.IntegerToken, token.CommaToken, token.IntegerToken,
				token.RightParenToken, token.SemicolonToken, token.RightBraceToken, token.SemicolonToken,

				// x = [1,2];
				token.IdentifierToken, token.AssignToken, token.LeftBracketToken, token.IntegerToken, token.CommaToken,
				token.IntegerToken, token.RightBracketToken, token.SemicolonToken,

				// x = x[1];
				token.IdentifierToken, token.AssignToken, token.IdentifierToken, token.LeftBracketToken, token.IntegerToken,
				token.RightBracketToken, token.SemicolonToken,

				// var isEqual = hello == 1;
				token.VarToken, token.IdentifierToken, token.AssignToken, token.IdentifierToken, token.EqualToken,
				token.IntegerToken, token.SemicolonToken,

				// var isEqual = true != (3 == 2);
				token.VarToken, token.IdentifierToken, token.AssignToken, token.BooleanToken, token.NotEqualToken,
				token.LeftParenToken, token.IntegerToken, token.EqualToken, token.IntegerToken, token.RightParenToken,
				token.SemicolonToken,

				// while (a < b & a <= b & a > b & a >= b) { x = x+1; };
				token.WhileToken, token.LeftParenToken, token.IdentifierToken, token.LThanToken, token.IdentifierToken,
				token.AndToken, token.IdentifierToken, token.LThanEqualToken, token.IdentifierToken, token.AndToken,
				token.IdentifierToken, token.GThanToken, token.IdentifierToken, token.AndToken, token.IdentifierToken,
				token.GThanEqualToken, token.IdentifierToken, token.RightParenToken, token.LeftBraceToken, token.IdentifierToken,
				token.AssignToken, token.IdentifierToken, token.AddToken, token.IntegerToken, token.SemicolonToken,
				token.RightBraceToken, token.SemicolonToken,


			},
			[]string{
				//var x = 22;
				"var", "x", "=", "22", ";",

				//var u = (a + b) / 22;
				"var", "u", "=", "(", "a", "+", "b", ")", "/", "22", ";",

				//var hello = [a,2,"hello",4];
				"var", "hello", "=", "[", "a", ",", "2", ",", "\"", "hello", "\"", ",", "4", "]", ";",

				// x = 100 / [a,2,"hello",22.33];
				"x", "=", "100", "/", "[", "a", ",", "2", ",", "\"", "hello", "\"", ",", "-", "22.33", "]", ";",

				// func testFunc(a,b,c,d,e) { print(22); d = a * b + c; return a + b - c / d * e;};
				"func", "testFunc", "(", "a", ",", "b", ",", "c", ",", "d", ",", "e", ")", "{", "print", "(", "22", ")",
				";", "d", "=", "a", "*", "b", "+", "c", ";", "return", "a", "+", "b", "-", "c", "/", "d", "*", "e", ";",
				"}", ";",

				// if (true & false | false) { testFunc(1,2,3,4,5);} else { testFunc(5,4,3,2,1);};
				"if", "(", "true", "&", "false", "|", "false", ")", "{", "testFunc", "(", "1", ",", "2", ",", "3", ",",
				"4", ",", "5", ")", ";", "}", "else", "{", "testFunc", "(", "5", ",", "4", ",", "3", ",", "2", ",", "1",
				")", ";", "}", ";",

				// x = [1,2];
				"x", "=", "[", "1", ",", "2", "]", ";",

				// x = x[1];
				"x", "=", "x", "[", "1", "]", ";",

				// var isEqual = hello == 1;
				"var", "isEqual", "=", "hello", "==", "1", ";",

				// var isEqual = true != (3 == 2);
				"var", "isEqual", "=", "true", "!=", "(", "3", "==", "2", ")", ";",

				// while (a < b & a <= b & a > b & a >= b) { x = x+1; };
				"while", "(", "a", "<", "b", "&", "a", "<=", "b", "&", "a", ">", "b", "&", "a", ">=", "b", ")", "{", "x",
				"=", "x", "+", "1", ";", "}", ";",








			},
	}

	var (
		f *os.File
		l Lexer
		cT token.Token
		err error
	)

	if f, err = os.Open(tC.inputFilePath); err != nil {
		t.Errorf(err.Error())
	}

	if l, err = NewLexer(f); err != nil {
		t.Errorf(err.Error())
	}

	for i := 0; i < len(tC.outputLiterals); i++ {
		if cT, err = l.NextToken(); err != nil {
			t.Errorf(err.Error())
			return
		}

		if cT.Type() != tC.outputTypes[i] {
			t.Errorf(fmt.Sprintf(internal.ErrType, tC.outputTypes[i], cT.Type()))
			return
		}

		if cT.Literal() != tC.outputLiterals[i] {
			t.Errorf(fmt.Sprintf(internal.ErrInvalidToken, cT.Literal(), tC.outputLiterals[i]))
			return
		}

	}

}