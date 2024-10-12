package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	INT    = "INT"
	STRING = "STRING"

	IF    = "IF"
	ELSE  = "ELSE"
	WHILE = "WHILE"
	IDENT = "IDENT"

	ASSIGN = "="
	PLUS   = "+"
	ASTR   = "*"
	SLASH  = "/"
	MODULO = "%"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"

	COMMA     = ","
	SEMICOLON = ";"
	EOF       = "EOF"

	LPAR  = "("
	RPAR  = ")"
	LBRAC = "{"
	RBRAC = "}"

	VAR         = "VAR"
	FUNCTION    = "FUNCTION"
	NOT_ALLOWED = "NOT_ALLOWED"
	RETURN      = "RETURN"

	BANG  = "!"
	MINUS = "-"
)

// keywords dict for indetifiers
var keywords = map[string]TokenType{
	"var":    VAR,
	"false":  FALSE,
	"true":   TRUE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"fn":     FUNCTION,
	"return": RETURN,
}

// if our input is in our keywords map return the token othewise
// return the IDENT
func LookUpIdent(ident string) TokenType {
	if tkn, ok := keywords[ident]; ok {
		return tkn
	} else {
		return IDENT
	}
}