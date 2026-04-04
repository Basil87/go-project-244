package diff

type DiffStatus int

const (
	StatusUnchanged DiffStatus = iota
	StatusRemoved
	StatusAdded
	StatusChanged
	StatusNested
)

type DiffNode struct {
	Key      string
	Status   DiffStatus
	OldVal   any
	NewVal   any
	Children []DiffNode
}
