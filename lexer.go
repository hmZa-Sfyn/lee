// lexer.go
package main

import (
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	Illegal TokenType = iota
	EOF

	// Identifiers & literals
	Ident
	Int
	Float
	String

	// Keywords
	True
	False
	Let
	Mut
	If
	Else
	While
	Foreach
	In
	Print

	// Operators & delimiters
	Colon     // :
	Pipe      // |
	Arrow     // ->
	Eq        // =
	Plus      // +
	Minus     // -
	Star      // *
	Slash     // /
	Percent   // %
	EqEq      // ==
	NotEq     // !=
	Less      // <
	Greater   // >
	LessEq    // <=
	GreaterEq // >=
	AndAnd    // &&
	OrOr      // ||
	Bang      // !
	LParen    // (
	RParen    // )
	LBracket  // [
	RBracket  // ]
	LBrace    // {
	RBrace    // }
	Comma     // ,
	Dollar    // $
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input   string
	pos     int  // current position (byte index)
	readPos int  // reading position (after current char)
	ch      rune // current char under examination
	line    int  // 1-based
	col     int  // 1-based
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   1, // columns start at 1
	}
	l.readChar() // initialize first character
	return l
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	l.skipComment()

	var tok Token
	tok.Line = l.line
	tok.Col = l.col

	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Value = ""

	case ':':
		tok = newToken(Colon, string(l.ch))

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:  OrOr,
				Value: string(ch) + string(l.ch),
				Line:  l.line,
				Col:   l.col - 1,
			}
		} else {
			tok = newToken(Pipe, string(l.ch))
		}

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EqEq, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok = newToken(Eq, string(l.ch))
		}

	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: Arrow, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok = newToken(Minus, string(l.ch))
		}

	case '+':
		tok = newToken(Plus, string(l.ch))

	case '*':
		tok = newToken(Star, string(l.ch))

	case '/':
		tok = newToken(Slash, string(l.ch))

	case '%':
		tok = newToken(Percent, string(l.ch))

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NotEq, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok = newToken(Bang, string(l.ch))
		}

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LessEq, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok = newToken(Less, string(l.ch))
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GreaterEq, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok = newToken(Greater, string(l.ch))
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AndAnd, Value: string(ch) + string(l.ch), Line: l.line, Col: l.col - 1}
		} else {
			tok.Type = Illegal
			tok.Value = string(l.ch)
		}

	case '(':
		tok = newToken(LParen, string(l.ch))
	case ')':
		tok = newToken(RParen, string(l.ch))
	case '[':
		tok = newToken(LBracket, string(l.ch))
	case ']':
		tok = newToken(RBracket, string(l.ch))
	case '{':
		tok = newToken(LBrace, string(l.ch))
	case '}':
		tok = newToken(RBrace, string(l.ch))
	case ',':
		tok = newToken(Comma, string(l.ch))

	case '"':
		tok.Type = String
		tok.Value = l.readString()
		return tok // early return because readString already advanced

	case '$':
		if l.peekChar() == '"' {
			l.readChar() // consume $
			tok.Type = String
			tok.Value = l.readString()
			return tok
		}
		tok.Type = Dollar
		tok.Value = string(l.ch)

	default:
		if unicode.IsLetter(l.ch) || l.ch == '_' {
			tok.Value = l.readIdentifier()
			tok.Type = lookupIdent(tok.Value)
			return tok // early return
		}
		if unicode.IsDigit(l.ch) || (l.ch == '.' && unicode.IsDigit(l.peekChar())) {
			tok.Value = l.readNumber()
			if len(tok.Value) > 0 && (tok.Value[0] == '.' || tok.Value[len(tok.Value)-1] == '.') {
				tok.Type = Float
			} else {
				tok.Type = Int
			}
			return tok // early return
		}
		tok.Type = Illegal
		tok.Value = string(l.ch)
	}

	l.readChar()
	return tok
}

// ──────────────────────────────────────────────────────────────────────────────
// Helper methods
// ──────────────────────────────────────────────────────────────────────────────

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		r, width := utf8.DecodeRuneInString(l.input[l.readPos:])
		l.ch = r
		l.readPos += width
	}
	l.pos = l.readPos - utf8.RuneLen(l.ch)

	if l.ch == '\n' {
		l.line++
		l.col = 1
	} else if l.ch != 0 {
		l.col++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPos:])
	return r
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	if l.ch == '#' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos
	hasDot := false

	for unicode.IsDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}

	// Allow trailing dot only if followed by digit (already checked in loop)
	return l.input[start:l.pos]
}

func (l *Lexer) readString() string {
	l.readChar() // consume opening " (already checked in caller)
	start := l.pos

	for l.ch != '"' && l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}

	str := l.input[start:l.pos]

	// Consume closing quote if present
	if l.ch == '"' {
		l.readChar()
	}

	return str
}

func newToken(tt TokenType, value string) Token {
	return Token{Type: tt, Value: value}
}

var keywords = map[string]TokenType{
	"true":    True,
	"false":   False,
	"let":     Let,
	"mut":     Mut,
	"if":      If,
	"else":    Else,
	"while":   While,
	"foreach": Foreach,
	"in":      In,
	"print":   Print,
}

func lookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return Ident
}
