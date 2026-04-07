package formatters

import (
	"code/diff"
	"encoding/json"
)

type jsonNode struct {
	Key      string     `json:"key"`
	Status   string     `json:"status"`
	OldVal   any        `json:"oldValue,omitempty"`
	NewVal   any        `json:"newValue,omitempty"`
	Children []jsonNode `json:"children,omitempty"`
}

func FormatJSON(nodes []diff.DiffNode) string {
	jNodes := toJSONNodes(nodes)
	data, err := json.MarshalIndent(jNodes, "", "    ")
	if err != nil {
		return "[]"
	}
	return string(data)
}

func toJSONNodes(nodes []diff.DiffNode) []jsonNode {
	result := make([]jsonNode, 0, len(nodes))
	for _, n := range nodes {
		jn := jsonNode{Key: n.Key, Status: diffStatusString(n.Status)}
		switch n.Status {
		case diff.StatusAdded:
			jn.NewVal = n.NewVal
		case diff.StatusRemoved:
			jn.OldVal = n.OldVal
		case diff.StatusChanged:
			jn.OldVal = n.OldVal
			jn.NewVal = n.NewVal
		case diff.StatusUnchanged:
			jn.OldVal = n.OldVal
		case diff.StatusNested:
			jn.Children = toJSONNodes(n.Children)
		}
		result = append(result, jn)
	}
	return result
}

func diffStatusString(s diff.DiffStatus) string {
	switch s {
	case diff.StatusAdded:
		return "added"
	case diff.StatusRemoved:
		return "removed"
	case diff.StatusChanged:
		return "changed"
	case diff.StatusNested:
		return "nested"
	default:
		return "unchanged"
	}
}
