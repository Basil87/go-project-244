package diff

// DiffStatus represents the change state of a diff node.
type DiffStatus int

const (
	// StatusUnchanged indicates the value was not modified.
	StatusUnchanged DiffStatus = iota
	// StatusRemoved indicates the key was present in the first file only.
	StatusRemoved
	// StatusAdded indicates the key was present in the second file only.
	StatusAdded
	// StatusChanged indicates the value changed between the two files.
	StatusChanged
	// StatusNested indicates the value is an object whose children are compared recursively.
	StatusNested
)

// DiffNode represents a single node in a diff tree.
type DiffNode struct {
	Key      string
	Status   DiffStatus
	OldVal   any
	NewVal   any
	Children []DiffNode
}
