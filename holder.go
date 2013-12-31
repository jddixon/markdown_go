package md

import ()

type Holder struct {
	children []BlockI
}

func NewHolder() *Holder {
	var h = new(Holder)
	return h
}

func (h *Holder) AddChild(child BlockI) (err error) {
	if child == nil {
		err = NilChild
	} else {
		// XXX We don't prevent duplicates
		h.children = append(h.children, child)
	}
	return
}

func (h *Holder) Size() int {
	return len(h.children)
}

func (h *Holder) GetChild(n int) (child BlockI, err error) {
	if n < 0 || h.Size() <= n {
		err = ChildNdxOutOfRange
	} else {
		child = h.children[n]
	}
	return
}
