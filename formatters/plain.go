package formatters

import (
	"code/diff"
	"fmt"
	"strings"
)

func FormatPlain(nodes []diff.DiffNode) string {
	return strings.TrimRight(renderPlain(nodes, ""), "\n")
}

func renderPlain(nodes []diff.DiffNode, prefix string) string {
	var sb strings.Builder
	for _, node := range nodes {
		path := node.Key
		if prefix != "" {
			path = prefix + "." + node.Key
		}
		switch node.Status {
		case diff.StatusAdded:
			sb.WriteString(fmt.Sprintf("Property '%s' was added with value: %s\n", path, plainValue(node.NewVal)))
		case diff.StatusRemoved:
			sb.WriteString(fmt.Sprintf("Property '%s' was removed\n", path))
		case diff.StatusChanged:
			sb.WriteString(fmt.Sprintf("Property '%s' was updated. From %s to %s\n", path, plainValue(node.OldVal), plainValue(node.NewVal)))
		case diff.StatusNested:
			sb.WriteString(renderPlain(node.Children, path))
		}
	}
	return sb.String()
}

func plainValue(v any) string {
	if v == nil {
		return "null"
	}
	if _, ok := v.(map[string]any); ok {
		return "[complex value]"
	}
	if s, ok := v.(string); ok {
		return fmt.Sprintf("'%s'", s)
	}
	if f, ok := v.(float64); ok {
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%g", f)
	}
	return fmt.Sprintf("%v", v)
}
