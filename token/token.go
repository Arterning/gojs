package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 123456
	FLOAT  = "FLOAT"  // 123.456
	STRING = "STRING" // "foobar"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"
	UNDEFINED = "UNDEFINED"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"

	LT     = "<"
	GT     = ">"
	LE     = "<="
	GE     = ">="
	EQ     = "=="
	NOT_EQ = "!="
	STRICT_EQ = "==="
	STRICT_NOT_EQ = "!=="

	AND = "&&"
	OR  = "||"

	PLUS_ASSIGN  = "+="
	MINUS_ASSIGN = "-="
	TIMES_ASSIGN = "*="
	DIV_ASSIGN   = "/="

	INCREMENT = "++"
	DECREMENT = "--"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."
	QUESTION  = "?"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	ARROW = "=>"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	VAR      = "VAR"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	WHILE    = "WHILE"
	DO       = "DO"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	NEW      = "NEW"
	THIS     = "THIS"
	TRY      = "TRY"
	CATCH    = "CATCH"
	FINALLY  = "FINALLY"
	THROW    = "THROW"
	TYPEOF   = "TYPEOF"
	DELETE   = "DELETE"
	IN       = "IN"
	OF       = "OF"
	INSTANCEOF = "INSTANCEOF"
	ASYNC    = "ASYNC"
	AWAIT    = "AWAIT"
	CLASS    = "CLASS"
	EXTENDS  = "EXTENDS"
	STATIC   = "STATIC"
	SUPER    = "SUPER"
	IMPORT   = "IMPORT"
	EXPORT   = "EXPORT"
	FROM     = "FROM"
	AS       = "AS"
)

var keywords = map[string]TokenType{
	"function":   FUNCTION,
	"let":        LET,
	"const":      CONST,
	"var":        VAR,
	"return":     RETURN,
	"if":         IF,
	"else":       ELSE,
	"true":       TRUE,
	"false":      FALSE,
	"null":       NULL,
	"undefined":  UNDEFINED,
	"for":        FOR,
	"while":      WHILE,
	"do":         DO,
	"break":      BREAK,
	"continue":   CONTINUE,
	"switch":     SWITCH,
	"case":       CASE,
	"default":    DEFAULT,
	"new":        NEW,
	"this":       THIS,
	"try":        TRY,
	"catch":      CATCH,
	"finally":    FINALLY,
	"throw":      THROW,
	"typeof":     TYPEOF,
	"delete":     DELETE,
	"in":         IN,
	"of":         OF,
	"instanceof": INSTANCEOF,
	"async":      ASYNC,
	"await":      AWAIT,
	"class":      CLASS,
	"extends":    EXTENDS,
	"static":     STATIC,
	"super":      SUPER,
	"import":     IMPORT,
	"export":     EXPORT,
	"from":       FROM,
	"as":         AS,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
