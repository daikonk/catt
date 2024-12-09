package lexer

import "go_interpreter/token"

// Defining our Lexer Obj
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// Our Lexer Constructor
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// reads the character at readPosition, sets the new position and increases the readPosition by 1
// sets ch to after we reach the end of the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// handles the "var map" for the lexer (token types defined in token.go)
func (l *Lexer) NextToken() token.Token {
	var tkn token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '%':
		tkn = newToken(token.MODULO, l.ch)
	case '"':
		tkn.Type = token.STRING
		tkn.Literal = l.ReadString()
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(token.BANG, l.ch)
		}
	case ';':
		tkn = newToken(token.SEMICOLON, l.ch)
	case '(':
		tkn = newToken(token.LPAR, l.ch)
	case ')':
		tkn = newToken(token.RPAR, l.ch)
	case '{':
		tkn = newToken(token.LBRAC, l.ch)
	case '}':
		tkn = newToken(token.RBRAC, l.ch)
	case ',':
		tkn = newToken(token.COMMA, l.ch)
	case '+':
		tkn = newToken(token.PLUS, l.ch)
	case '-':
		tkn = newToken(token.MINUS, l.ch)
	case '/':
		tkn = newToken(token.SLASH, l.ch)
	case '*':
		tkn = newToken(token.ASTR, l.ch)
	case '<':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.CHAN_OP, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(token.LT, l.ch)
		}
	case '>':
		tkn = newToken(token.GT, l.ch)
	case 0:
		tkn.Literal = ""
		tkn.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tkn.Literal = l.readIdentifier()
			tkn.Type = token.LookUpIdent(tkn.Literal)
			return tkn
		} else if isDigit(l.ch) {
			tkn.Literal = l.readNumber()
			tkn.Type = token.INT
			return tkn
		}
		tkn = newToken(token.NOT_ALLOWED, l.ch)
	}
	l.readChar()
	return tkn
}

func (l *Lexer) ReadString() string {
	position := l.readPosition
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// new token constructor
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// reads/parses input for identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// checks if input is a valid letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
