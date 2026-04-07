package code

import (
	"code/diff"
	// "code/formatters"
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

type Formatter func([]diff.DiffNode) string

func GetDiff(file1, file2 string, format Formatter) (string, error) {
	fileContent1, err := GetFileData(file1)
	if err != nil {
		return "", err
	}
	fileContent2, err := GetFileData(file2)
	if err != nil {
		return "", err
	}
	return format(buildDiff(fileContent1, fileContent2)), nil
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
	defer func() {
		_ = f.Close()
	}()

	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

func buildDiff(m1, m2 map[string]any) []diff.DiffNode {
	keys := allKeys(m1, m2)
	sort.Strings(keys)
	var nodes []diff.DiffNode
	for _, key := range keys {
		v1, in1 := m1[key]
		v2, in2 := m2[key]
		switch {
		case in1 && in2:
			sub1, ok1 := v1.(map[string]any)
			sub2, ok2 := v2.(map[string]any)
			if ok1 && ok2 {
				nodes = append(nodes, diff.DiffNode{Key: key, Status: diff.StatusNested, Children: buildDiff(sub1, sub2)})
			} else if reflect.DeepEqual(v1, v2) {
				nodes = append(nodes, diff.DiffNode{Key: key, Status: diff.StatusUnchanged, OldVal: v1})
			} else {
				nodes = append(nodes, diff.DiffNode{Key: key, Status: diff.StatusChanged, OldVal: v1, NewVal: v2})
			}
		case in1:
			nodes = append(nodes, diff.DiffNode{Key: key, Status: diff.StatusRemoved, OldVal: v1})
		default:
			nodes = append(nodes, diff.DiffNode{Key: key, Status: diff.StatusAdded, NewVal: v2})
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
