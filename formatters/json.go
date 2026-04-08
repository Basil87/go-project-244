package formatters

import (
	"code/diff"
	"encoding/json"
)

func FormatJSON(nodes []diff.DiffNode) string {
	m := toJSONMap(nodes)
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func toJSONMap(nodes []diff.DiffNode) map[string]any {
	result := make(map[string]any, len(nodes))
	for _, n := range nodes {
		switch n.Status {
		case diff.StatusAdded:
			result[n.Key] = map[string]any{"type": "added", "value": n.NewVal}
		case diff.StatusRemoved:
			result[n.Key] = map[string]any{"type": "removed", "value": n.OldVal}
		case diff.StatusChanged:
			result[n.Key] = map[string]any{"type": "changed", "from": n.OldVal, "to": n.NewVal}
		case diff.StatusUnchanged:
			result[n.Key] = map[string]any{"type": "unchanged", "value": n.OldVal}
		case diff.StatusNested:
			result[n.Key] = map[string]any{"type": "nested", "children": toJSONMap(n.Children)}
		}
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
