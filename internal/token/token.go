package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	NEWLINE
	INDENT
	DEDENT

	IDENT
	NUMBER
	STRING
	LBRACKET
	RBRACKET
	COMMA

	ASSIGN
	PLUS
	GT
	LT
	LPAREN
	RPAREN
	COLON

	FUNCTIA
	DRUKUVATY
	NEKHAY
	VVID
	YAKSHO
	INACKSHE
	KORIN
	LOH10
	ABS
	ARKSYN
	KOSYNUS
	STUPIN
	VERNUTY
	POKY
)

func (t TokenType) String() string {
	names := [...]string{
		"ILLEGAL", "EOF", "NEWLINE", "INDENT", "DEDENT",
		"IDENT", "NUMBER", "STRING", "[", "]", ",",
		"=", "+", ">", "<", "(", ")", ":",
		"functia", "drukuvaty", "nekhay", "vvid", "yaksho", "inackshe",
		"korin", "loh10", "abs", "arksyn", "kosynus", "stupin", "vernuty",
	}
	if t >= 0 && int(t) < len(names) {
		return names[t]
	}
	return "UNKNOWN"
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func LookupIdent(ident string) TokenType {
	switch ident {
	case "functia":
		return FUNCTIA
	case "drukuvaty":
		return DRUKUVATY
	case "nekhay":
		return NEKHAY
	case "vvid":
		return VVID
	case "yaksho":
		return YAKSHO
	case "inackshe":
		return INACKSHE
	case "korin":
		return KORIN
	case "loh10":
		return LOH10
	case "abs":
		return ABS
	case "arksyn":
		return ARKSYN
	case "kosynus":
		return KOSYNUS
	case "stupin":
		return STUPIN
	case "vernuty":
		return VERNUTY
	case "poky":
		return POKY
	default:
		return IDENT
	}
}
