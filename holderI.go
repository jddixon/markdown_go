package md

// Must be kept in sync with holder.go

type HolderI interface {
	AddChild(child BlockI) (err error)
	Size() int
	GetChild(n int) (child BlockI, err error)
	BlockI
}
