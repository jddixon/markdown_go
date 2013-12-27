package md

// xgo/md/const.go

const (
	BACKSLASH       = '\\'
	CR              = '\r'
	LF              = '\n'
	MAX_ENTITY_CHAR = 6 // between & and ;
	SPACE           = ' '
	TAB             = '\t'
)

var (
	SEP_CHAR    = []rune{CR, LF}
	FOUR_SPACES = []rune("    ")

	OPEN_EM      = []rune("<em>")
	CLOSE_EM     = []rune("</em>")
	H_RULE       = []rune("<hr />")
	OPEN_STRONG  = []rune("<strong>")
	CLOSE_STRONG = []rune("</strong>")
)
