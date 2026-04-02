package code

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"
)

func GetDiff(file1, file2 string) (string, error) {
	fileContent1, err := GetFileData(file1)
	if err != nil {
		return "", err
	}
	fileContent2, err := GetFileData(file2)
	if err != nil {
		return "", err
	}

	result, err := compareJsons(fileContent1, fileContent2)
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetFileData(path string) (map[string]any, error) {
	info, err := os.Stat(path)

	var fileContent map[string]any
	if err != nil {
		return fileContent, fmt.Errorf("file not exists: %s", path)
	}

	if info.IsDir() {
		return fileContent, fmt.Errorf("expected file, not a directory")
	}

	contentType, err := detectFileType(path)
	if err != nil {
		return fileContent, err
	}

	if contentType != "application/json" && contentType != "text/plain; charset=utf-8" {
		return fileContent, fmt.Errorf("unsupported file type: %s", contentType)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fileContent, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &fileContent); err != nil {
			return fileContent, fmt.Errorf("invalid yaml: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &fileContent); err != nil {
			return fileContent, fmt.Errorf("invalid json: %w", err)
		}
	}

	return fileContent, nil
}

func detectFileType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

type diffStatus int

const (
	statusUnchanged diffStatus = iota
	statusRemoved
	statusAdded
	statusChanged
	statusNested
)

type diffNode struct {
	key      string
	status   diffStatus
	oldVal   any
	newVal   any
	children []diffNode
}

func compareJsons(fileContent1, fileContent2 map[string]any) (string, error) {
	nodes := buildDiff(fileContent1, fileContent2)
	return renderDiff(nodes, 1), nil
}

func buildDiff(m1, m2 map[string]any) []diffNode {
	keys := allKeys(m1, m2)
	sort.Strings(keys)
	var nodes []diffNode
	for _, key := range keys {
		v1, in1 := m1[key]
		v2, in2 := m2[key]
		switch {
		case in1 && in2:
			sub1, ok1 := v1.(map[string]any)
			sub2, ok2 := v2.(map[string]any)
			if ok1 && ok2 {
				nodes = append(nodes, diffNode{key: key, status: statusNested, children: buildDiff(sub1, sub2)})
			} else if reflect.DeepEqual(v1, v2) {
				nodes = append(nodes, diffNode{key: key, status: statusUnchanged, oldVal: v1})
			} else {
				nodes = append(nodes, diffNode{key: key, status: statusChanged, oldVal: v1, newVal: v2})
			}
		case in1:
			nodes = append(nodes, diffNode{key: key, status: statusRemoved, oldVal: v1})
		default:
			nodes = append(nodes, diffNode{key: key, status: statusAdded, newVal: v2})
		}
	}
	return nodes
}

func allKeys(m1, m2 map[string]any) []string {
	seen := make(map[string]bool)
	for k := range m1 {
		seen[k] = true
	}
	for k := range m2 {
		seen[k] = true
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
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

func prefixOrder(s string) int {
	if strings.HasPrefix(s, "- ") {
		return 0
	}
	if strings.HasPrefix(s, "+ ") {
		return 2
	}
	return 1
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "+ ") || strings.HasPrefix(s, "- ") {
		return s[2:]
	}
	return s
}
