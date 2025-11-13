package lexer

import (
	"gojs/token"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.STRICT_EQ, Literal: "===", Line: tok.Line, Column: tok.Column}
			} else {
				tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
			}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "=>", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ASSIGN, l.ch, tok.Line, tok.Column)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: "++", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: "+=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.PLUS, l.ch, tok.Line, tok.Column)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: "--", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: "-=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.MINUS, l.ch, tok.Line, tok.Column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.STRICT_NOT_EQ, Literal: "!==", Line: tok.Line, Column: tok.Column}
			} else {
				tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
			}
		} else {
			tok = newToken(token.BANG, l.ch, tok.Line, tok.Column)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.DIV_ASSIGN, Literal: "/=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.SLASH, l.ch, tok.Line, tok.Column)
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.TIMES_ASSIGN, Literal: "*=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ASTERISK, l.ch, tok.Line, tok.Column)
		}
	case '%':
		tok = newToken(token.PERCENT, l.ch, tok.Line, tok.Column)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: "<=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.LT, l.ch, tok.Line, tok.Column)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: ">=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.GT, l.ch, tok.Line, tok.Column)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: "&&", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: "||", Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, tok.Line, tok.Column)
	case ':':
		tok = newToken(token.COLON, l.ch, tok.Line, tok.Column)
	case '?':
		tok = newToken(token.QUESTION, l.ch, tok.Line, tok.Column)
	case ',':
		tok = newToken(token.COMMA, l.ch, tok.Line, tok.Column)
	case '.':
		tok = newToken(token.DOT, l.ch, tok.Line, tok.Column)
	case '(':
		tok = newToken(token.LPAREN, l.ch, tok.Line, tok.Column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, tok.Line, tok.Column)
	case '{':
		tok = newToken(token.LBRACE, l.ch, tok.Line, tok.Column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, tok.Line, tok.Column)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, tok.Line, tok.Column)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, tok.Line, tok.Column)
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
	case '`':
		tok.Type = token.STRING
		tok.Literal = l.readTemplateString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber(tok.Line, tok.Column)
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // skip '/'
	l.readChar() // skip '*'
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip '*'
			l.readChar() // skip '/'
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber(line, column int) token.Token {
	position := l.position
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // skip '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Handle scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		isFloat = true
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	literal := l.input[position:l.position]
	tokenType := token.INT
	if isFloat {
		tokenType = token.FLOAT
	}

	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // skip escape character
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTemplateString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_' || ch == '$'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}
