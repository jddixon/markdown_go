package md

// xgo/md/entitySpan.go

import ()

type EntitySpan struct {
	runes []rune
}

func (e *EntitySpan) String() (s string) {
	return string(e.GetHtml())
}

func (e *EntitySpan) GetHtml() (out []rune) {
	out = append(out, '&')
	out = append(out, e.runes...)
	out = append(out, ';')
	return
}

// Attempt to parse out an EntitySpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
func (q *Line) parseEntitySpan() (span SpanI, err error) {

	// we are guaranteed that on entry the character at offset is '&'
	offset := q.offset
	eol := uint(len(q.runes))
	maxOffset := offset + MAX_ENTITY_CHARS + 1
	if maxOffset > eol {
		maxOffset = eol
	}
	semiAt := uint(0) // location of the closing semicolon
	for offset++; offset < maxOffset; offset++ {
		if q.runes[offset] == ';' {
			semiAt = offset
			break
		}
	}
	if semiAt > 0 {

		var body []rune
		for i := q.offset + 1; i < semiAt; i++ {
			body = append(body, q.runes[i])
		}
		span = &EntitySpan{
			runes: body,
		}
		q.offset = semiAt + 1
	}
	return
}
