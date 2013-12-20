package md

// xgo/md/emph.go

import ()

type EmphSpan struct {
	emphDoubled bool
	runes       []rune
}

func (e *EmphSpan) Get() (out []rune) {
	if e.emphDoubled {
		out = append(out, OPEN_STRONG...)
		out = append(out, e.runes...)
		out = append(out, CLOSE_STRONG...)
	} else {
		out = append(out, OPEN_EM...)
		out = append(out, e.runes...)
		out = append(out, CLOSE_EM...)
	}
	return
}

// Attempt to parse out an EmphSpan, returning a SpanI reference
// to it on success and nil and possibly an error on failure.  If the parse
// fails but there is no input error, leave the line offset unchanged and
// return a nil SpanI.  If the parse succeeds, return a SpanI and advance
// the offset accordingly.
//
func (q *Line) parseEmphSpan() (span SpanI, err error) {

	emphChar := q.runes[q.offset]
	offset := q.offset
	emphDoubled := false
	firstChar := true
	found := false

	// determine whether the emphasis is doubled, then look for the end
	// of the span
	for offset++; offset < len(q.runes); offset++ {
		ch := q.runes[offset]
		if firstChar {
			firstChar = false
			if ch == emphChar {
				emphDoubled = true
				continue
			}
		}
		if ch == emphChar {
			if emphDoubled {
				if offset+1 < len(q.runes) &&
					q.runes[offset+1] == emphChar {
					offset++
					found = true
					break
				}
			} else {
				found = true
				break
			}
		}
	}
	if found {
		var start, end int
		if emphDoubled {
			start = q.offset + 2
			end = offset - 1
		} else {
			start = q.offset + 1
			end = offset
		}
		q.offset = offset + 1
		span = &EmphSpan{
			emphDoubled: emphDoubled,
			runes:       q.runes[start:end],
		}
	}
	return
}
