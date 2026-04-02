package code

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
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

	if err := json.Unmarshal(data, &fileContent); err != nil {
		return fileContent, fmt.Errorf("invalid json: %w", err)
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

func compareJsons(fileContent1, fileContent2 map[string]any) (string, error) {
	result := make(map[string]any)
	var keysSlice []string
	for key, value := range fileContent1 {
		secondValue, ok := fileContent2[key]
		if !ok || secondValue != value {
			result["- "+key] = value
			keysSlice = append(keysSlice, "- "+key)
			continue
		}
		result[key] = value
		keysSlice = append(keysSlice, key)
		delete(fileContent2, key)
	}
	for key, value := range fileContent2 {
		result["+ "+key] = value
		keysSlice = append(keysSlice, "+ "+key)
	}
	sort.Slice(keysSlice, func(i, j int) bool {
		ni := normalize(keysSlice[i])
		nj := normalize(keysSlice[j])
		if ni != nj {
			return ni < nj
		}
		return prefixOrder(keysSlice[i]) < prefixOrder(keysSlice[j])
	})
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, key := range keysSlice {
		if strings.HasPrefix(key, "+ ") || strings.HasPrefix(key, "- ") {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, result[key]))
		} else {
			sb.WriteString(fmt.Sprintf("    %s: %v\n", key, result[key]))
		}
	}
	sb.WriteString("}")
	return sb.String(), nil
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
