package md

// xgo/md/const.go

const (
	BACKSLASH        = '\\'
	CR               = '\r'
	LF               = '\n'
	MAX_ENTITY_CHARS = 6 // between & and ;
	SPACE            = ' '
	TAB              = '\t'
)

var (
	AMP_ENTITY  = []rune("&amp;")
	COPY_ENTITY = []rune("&copy;")
	GT_ENTITY   = []rune("&gt;")
	LT_ENTITY   = []rune("&lt;")
)

var (
	HARD_SPACE  = rune(0x00a0)
	FOUR_SPACES = []rune("    ")
	MAIL_TO     = []rune("mailto:")
	SEP_RUNES   = []rune{CR, LF}

	// 4 is canonical, 3 is OK, backtick is OK
	FENCE = []rune("~~~~~\n")

	BLOCKQUOTE_OPEN  = []rune("<blockquote>\n")
	BLOCKQUOTE_CLOSE = []rune("</blockquote>\n")
	CODE_OPEN        = []rune("<pre><code>")
	CODE_CLOSE       = []rune("</code></pre>\n")
	EM_OPEN          = []rune("<em>")
	EM_CLOSE         = []rune("</em>")
	H_RULE           = []rune("<hr />\n")
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
