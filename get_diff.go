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
	fileContent2, err := GetFileData(file2)
	if err != nil {
		return "", err
	}

	if err := compareJsons(fileContent1, fileContent2); err != nil {
		return "", err
	}

	return "вот и результат", nil
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

func compareJsons(fileContent1, fileContent2 map[string]any) error {
	result := make(map[string]any)
	var keysSlice []string
	//var uniqueValue1 map[string]any
	//var uniqueValue2 map[string]any
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
		return normalize(keysSlice[i]) < normalize(keysSlice[j])
	})
	fmt.Println("{")
	for _, key := range keysSlice {
		fmt.Println("	", key, ": ", result[key])
	}
	fmt.Println("}")
	return nil
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "+ ") || strings.HasPrefix(s, "- ") {
		return s[2:]
	}
	return s
}
