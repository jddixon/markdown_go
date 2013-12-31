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

	BLOCKQUOTE_OPEN  = []rune("<blockquote>\n")
	BLOCKQUOTE_CLOSE = []rune("</blockquote>\n")
	EM_OPEN          = []rune("<em>")
	EM_CLOSE         = []rune("</em>")
	H_RULE           = []rune("<hr />")
	LI_OPEN          = []rune("<li>")
	LI_CLOSE         = []rune("</li>\n")
	OL_OPEN          = []rune("<ol>\n")
	OL_CLOSE         = []rune("</ol>\n")
	PARA_OPEN        = []rune("<p>")
	PARA_CLOSE       = []rune("</p>")
	STRONG_OPEN      = []rune("<strong>")
	STRONG_CLOSE     = []rune("</strong>")
	UL_OPEN          = []rune("<ul>\n")
	UL_CLOSE         = []rune("</ul>\n")
)
