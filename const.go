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
	LI_OPEN      = []rune("<li>")
	LI_CLOSE     = []rune("</li>\n")
	OL_OPEN      = []rune("<ol>\n")
	OL_CLOSE     = []rune("</ol>\n")
	OPEN_STRONG  = []rune("<strong>")
	CLOSE_STRONG = []rune("</strong>")
	UL_OPEN      = []rune("<ul>\n")
	UL_CLOSE     = []rune("</ul>\n")
)
