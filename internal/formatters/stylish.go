package formatters

import (
	"code/internal/diff"
	"fmt"
	"sort"
	"strings"
)

const (
	fmtUnchanged = "%s  %s: %s\n"
	fmtRemoved   = "%s- %s: %s\n"
	fmtAdded     = "%s+ %s: %s\n"
)

// FormatStylish formats the diff in the stylish tree format with +/- prefixes and indentation.
func FormatStylish(nodes []diff.DiffNode) string {
	return renderDiff(nodes, 1)
}

func renderDiff(nodes []diff.DiffNode, depth int) string {
	indent := strings.Repeat(" ", (depth-1)*4)
	signPrefix := strings.Repeat(" ", depth*4-2)
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, node := range nodes {
		switch node.Status {
		case diff.StatusNested:
			fmt.Fprintf(&sb, fmtUnchanged, signPrefix, node.Key, renderDiff(node.Children, depth+1))
		case diff.StatusUnchanged:
			fmt.Fprintf(&sb, fmtUnchanged, signPrefix, node.Key, formatValue(node.OldVal, depth))
		case diff.StatusRemoved:
			fmt.Fprintf(&sb, fmtRemoved, signPrefix, node.Key, formatValue(node.OldVal, depth))
		case diff.StatusAdded:
			fmt.Fprintf(&sb, fmtAdded, signPrefix, node.Key, formatValue(node.NewVal, depth))
		case diff.StatusChanged:
			fmt.Fprintf(&sb, fmtRemoved, signPrefix, node.Key, formatValue(node.OldVal, depth))
			fmt.Fprintf(&sb, fmtAdded, signPrefix, node.Key, formatValue(node.NewVal, depth))
		}
	}
	sb.WriteString(indent + "}")
	return sb.String()
}

func renderMap(m map[string]any, depth int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	indent := strings.Repeat(" ", (depth-1)*4)
	signPrefix := strings.Repeat(" ", depth*4-2)
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, k := range keys {
		fmt.Fprintf(&sb, fmtUnchanged, signPrefix, k, formatValue(m[k], depth))
	}
	sb.WriteString(indent + "}")
	return sb.String()
}

func formatScalar(v any) string {
	if v == nil {
		return "null"
	}
	if f, ok := v.(float64); ok {
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%g", f)
	}
	return fmt.Sprintf("%v", v)
}

func formatValue(v any, depth int) string {
	if m, ok := v.(map[string]any); ok {
		return renderMap(m, depth+1)
	}
	return formatScalar(v)
}
