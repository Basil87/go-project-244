package code

import (
	"fmt"
	"sort"
	"strings"
)

func FormatStylish(nodes []diffNode) string {
	return renderDiff(nodes, 1)
}

func renderDiff(nodes []diffNode, depth int) string {
	indent := strings.Repeat(" ", (depth-1)*4)
	signPrefix := strings.Repeat(" ", depth*4-2)
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, node := range nodes {
		switch node.status {
		case statusNested:
			sb.WriteString(fmt.Sprintf("%s  %s: %s\n", signPrefix, node.key, renderDiff(node.children, depth+1)))
		case statusUnchanged:
			sb.WriteString(fmt.Sprintf("%s  %s: %s\n", signPrefix, node.key, formatValue(node.oldVal, depth)))
		case statusRemoved:
			sb.WriteString(fmt.Sprintf("%s- %s: %s\n", signPrefix, node.key, formatValue(node.oldVal, depth)))
		case statusAdded:
			sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", signPrefix, node.key, formatValue(node.newVal, depth)))
		case statusChanged:
			sb.WriteString(fmt.Sprintf("%s- %s: %s\n", signPrefix, node.key, formatValue(node.oldVal, depth)))
			sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", signPrefix, node.key, formatValue(node.newVal, depth)))
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
		sb.WriteString(fmt.Sprintf("%s  %s: %s\n", signPrefix, k, formatValue(m[k], depth)))
	}
	sb.WriteString(indent + "}")
	return sb.String()
}

func formatValue(v any, depth int) string {
	if m, ok := v.(map[string]any); ok {
		return renderMap(m, depth+1)
	}
	return formatScalar(v)
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
