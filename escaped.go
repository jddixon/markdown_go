package md

// xgo/md/escaped.go

import ()

// Returns true if Markdown accepts the character as a valid part of an
// escape sequence.
func escaped(r rune) bool {
	n := int(r)
	return n == 0x21 || n == 0x23 || n == 0x26 || // ! #  &
		0x28 <= n && n <= 0x2B || // ( ) # +
		n == 0x2d || n == 0x2e || // - . (MINUS, DOT)
		0x5b <= n && n <= 0x5d || // [ \ ]
		n == 0x5f || n == 0x60 || // _ `(UNDERSCORE, BACKTICK)
		n == 0x7b || n == 0x7d
}
