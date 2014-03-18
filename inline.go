package md

// xgo/md/inline.go

// So far, just a list of names.  Of these, only <q> may be nested in
// an element of its own type.

// This is not an acceptable Go const 
var (
	INLINE_ELM = [...]string { 
	// Anchor -------------------------------------------------------
	"a",
	// Phrase Elements: General -------------------------------------
	"abbr",	
	"dfn",	
	"em",	
	"strong",	
	// Computer Phrase Elements -------------------------------------
	"code",	
	"samp",	
	"kbd",	
	"var",	
	// Presentation -------------------------------------------------
	"b",	
	"i",	
	"u",	
	"small",	
	"s",	
	// Span ---------------------------------------------------------
	"span",	
	// Others -------------------------------------------------------
	"br",			// need not be closed
	"bdo",	
	"cite",	
	"del",	
	"ins",	
	"q",	
	"sub",	
	"wbr",	
}
)
