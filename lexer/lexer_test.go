package lexer

import (
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"github.com/spf13/afero"
	"testing"
)

func TestLexer(t *testing.T) {

	tCs := []struct {
		input               []byte
		outputTokenTypes    []token.TokenType
		outputTokenLiterals []string
	}{
		{
			[]byte("var x = 22;"),
			[]token.TokenType{token.VarToken, token.IdentifierToken, token.AssignToken,
				token.IntegerToken, token.SemicolonToken},
			[]string{"var", "x", "=", "22", ";"},
		},
		{
			[]byte("var u = (a + b) / 22;"),
			[]token.TokenType{token.VarToken, token.IdentifierToken, token.AssignToken,
				token.LeftParenToken, token.IdentifierToken, token.AddToken, token.IdentifierToken, token.RightParenToken,
				token.DivToken, token.IntegerToken, token.SemicolonToken},
			[]string{"var", "u", "=", "(", "a", "+", "b", ")", "/", "22", ";"},
		},
		{
			[]byte("var hello = [a,2,\"hello\",4];"),
			[]token.TokenType{token.VarToken, token.IdentifierToken, token.AssignToken, token.LeftBracketToken,
				token.IdentifierToken, token.CommaToken, token.IntegerToken, token.CommaToken, token.QuotationMarkToken,
				token.IdentifierToken, token.QuotationMarkToken, token.CommaToken, token.IntegerToken, token.RightBracketToken,
				token.SemicolonToken},
			[]string{"var", "hello", "=", "[", "a", ",", "2", ",", "\"", "hello", "\"", ",", "4", "]",
				";"},
		},
		{
			[]byte(`func testFunc(a,b,c,d,e) {
print(22);
d = a * b + c;
return a + b - c / d * e;
};`),
			[]token.TokenType{token.FuncToken, token.IdentifierToken, token.LeftParenToken, token.IdentifierToken,
				token.CommaToken, token.IdentifierToken, token.CommaToken, token.IdentifierToken, token.CommaToken, token.IdentifierToken,
				token.CommaToken, token.IdentifierToken, token.RightParenToken, token.LeftBraceToken, token.IdentifierToken,
				token.LeftParenToken, token.IntegerToken, token.RightParenToken, token.SemicolonToken, token.IdentifierToken,
				token.AssignToken, token.IdentifierToken, token.MultToken, token.IdentifierToken, token.AddToken,
				token.IdentifierToken, token.SemicolonToken, token.ReturnToken, token.IdentifierToken, token.AddToken,
				token.IdentifierToken, token.SubToken, token.IdentifierToken, token.DivToken, token.IdentifierToken,
				token.MultToken, token.IdentifierToken, token.SemicolonToken, token.RightBraceToken, token.SemicolonToken},
			[]string{"func", "testFunc", "(", "a", ",", "b", ",", "c", ",", "d", ",", "e", ")", "{",
				"print", "(", "22", ")", ";", "d", "=", "a", "*", "b", "+", "c", ";", "return", "a", "+", "b", "-", "c",
				"/", "d", "*", "e", ";", "}", ";"},
		},
		{
			[]byte(`if (true & false | false) {
testFunc(1,2,3,4,5);
} else {
testFunc(5,4,3,2,1);
};`),
			[]token.TokenType{token.IfToken, token.LeftParenToken, token.BooleanToken, token.AndToken,
				token.BooleanToken, token.OrToken, token.BooleanToken, token.RightParenToken, token.LeftBraceToken,
				token.IdentifierToken, token.LeftParenToken, token.IntegerToken, token.CommaToken, token.IntegerToken,
				token.CommaToken, token.IntegerToken, token.CommaToken, token.IntegerToken, token.CommaToken,
				token.IntegerToken, token.RightParenToken, token.SemicolonToken, token.RightBraceToken, token.ElseToken,
				token.LeftBraceToken, token.IdentifierToken, token.LeftParenToken, token.IntegerToken, token.CommaToken,
				token.IntegerToken, token.CommaToken, token.IntegerToken, token.CommaToken, token.IntegerToken,
				token.CommaToken, token.IntegerToken, token.RightParenToken, token.SemicolonToken, token.RightBraceToken,
				token.SemicolonToken},
			[]string{"if", "(", "true", "&", "false", "|", "false", ")", "{", "testFunc", "(", "1",
				",", "2", ",", "3", ",", "4", ",", "5", ")", ";", "}", "else", "{", "testFunc", "(", "5", ",", "4",
				",", "3", ",", "2", ",", "1", ")", ";", "}", ";"},
		},
		{
			[]byte("x = [1,2];"),
			[]token.TokenType{token.IdentifierToken, token.AssignToken, token.LeftBracketToken,
				token.IntegerToken, token.CommaToken, token.IntegerToken, token.RightBracketToken, token.SemicolonToken},
			[]string{"x", "=", "[", "1", ",", "2", "]", ";"},
		},
		{
			[]byte("x = x[1];"),
			[]token.TokenType{token.IdentifierToken, token.AssignToken, token.IdentifierToken,
				token.LeftBracketToken, token.IntegerToken, token.RightBracketToken, token.SemicolonToken},
			[]string{"x", "=", "x", "[", "1", "]", ";"},
		},
		{
			[]byte(" var isEqual = hello == 1;"),
			[]token.TokenType{token.VarToken, token.IdentifierToken, token.AssignToken,
				token.IdentifierToken, token.EqualToken, token.IntegerToken, token.SemicolonToken},
			[]string{"var", "isEqual", "=", "hello", "==", "1", ";"},
		},
		{
			[]byte("isEqual = true != (3 == 2);"),
			[]token.TokenType{token.IdentifierToken, token.AssignToken, token.BooleanToken, token.NotEqualToken,
				token.LeftParenToken, token.IntegerToken, token.EqualToken, token.IntegerToken, token.RightParenToken,
				token.SemicolonToken},
			[]string{"isEqual", "=", "true", "!=", "(", "3", "==", "2", ")", ";"},
		},
		{
			[]byte(`while (a < b & a <= b & a > b & a >= b) {
x = x+1;
};`),
			[]token.TokenType{token.WhileToken, token.LeftParenToken, token.IdentifierToken, token.LThanToken,
				token.IdentifierToken, token.AndToken, token.IdentifierToken, token.LThanEqualToken, token.IdentifierToken,
				token.AndToken, token.IdentifierToken, token.GThanToken, token.IdentifierToken, token.AndToken, token.IdentifierToken,
				token.GThanEqualToken, token.IdentifierToken, token.RightParenToken, token.LeftBraceToken, token.IdentifierToken,
				token.AssignToken, token.IdentifierToken, token.AddToken, token.IntegerToken, token.SemicolonToken,
				token.RightBraceToken, token.SemicolonToken},
			[]string{"while", "(", "a", "<", "b", "&", "a", "<=", "b", "&", "a", ">", "b", "&", "a",
				">=", "b", ")", "{", "x", "=", "x", "+", "1", ";", "}", ";"},
		},
	}

	var (
		fs  afero.Fs
		err error
	)

	fs = afero.NewMemMapFs()
	if err = fs.MkdirAll("test_files/lexer_tests", 0755); err != nil {
		t.Fatalf(err.Error())
	}

	for i, tC := range tCs {
		var (
			f   afero.File
			l   Lexer
			cT  token.Token
			err error
			fp  string
		)

		fp = fmt.Sprintf("test_files/lexer_tests/test_%v.txt", i)

		if err = afero.WriteFile(fs, fp, tC.input, 0644); err != nil {
			t.Fatalf(err.Error())
		}

		if _, err = fs.Stat(fp); err != nil {
			t.Fatalf("file \"%s\" does not exist.\n", fp)
		}

		if f, err = fs.Open(fp); err != nil {
			t.Fatalf(err.Error())
		}

		if l, err = NewLexer(f); err != nil {
			t.Fatalf(err.Error())
		}

		for j := 0; j < len(tC.outputTokenLiterals); j++ {
			if cT, err = l.NextToken(); err != nil {
				t.Errorf(err.Error())
				continue
			}

			if cT.Type() != tC.outputTokenTypes[j] {
				t.Errorf(fmt.Sprintf(internal.ErrInvalidTokenTypeTest, i+1, cT.Type(), tC.outputTokenTypes[j]))
				continue
			}

			if cT.Literal() != tC.outputTokenLiterals[j] {
				t.Errorf(fmt.Sprintf(internal.ErrInvalidTokenLiteralTest, i+1, cT.Literal(), tC.outputTokenLiterals[j]))
				continue
			}

		}

	}

}
