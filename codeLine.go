package md

// xgo/md/codeLine.go

type CodeLine struct {
	runes []rune
}

func NewCodeLine(runes []rune) (t *CodeLine) {
	txt := make([]rune, len(runes))
	copy(txt, runes)
	return &CodeLine{runes: txt}
}

func (cl *CodeLine) String() string {
	return string(cl.runes) + "\n"
}

// Any backticks in CodeLine.Runes are literal backticks.
// In the current implementation, < and > are automatically 'escaped',
// in the sense that they are converted to character entities here.
func (p *CodeLine) GetHtml() (out []rune) {

	for i := 0; i < len(p.runes); i++ {
		r := p.runes[i]
		if r == '&' {
			out = append(out, AMP_ENTITY...)
		} else if r == '<' {
			out = append(out, LT_ENTITY...)
		} else if r == '>' {
			out = append(out, GT_ENTITY...)
		} else {
			out = append(out, r)
		}
	}
	return
}
