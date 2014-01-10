package md

// Must be kept in sync with holder.go

type HolderI interface {
	AddBlock(block BlockI) (err error)
	Size() int
	GetBlock(n int) (block BlockI, err error)
	BlockI
}
