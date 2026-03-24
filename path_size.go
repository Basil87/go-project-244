package code

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetParsData(file1, file2 string) (string, error) {
	fileContent1, err := GetFileData(file1)
	fileContent2, err := GetFileData(file2)
	if err != nil {
		return "", err
	}

	fmt.Println(fileContent1, fileContent2)

	return "вот и результат", nil
}

func GetFileData(path string) (string, error) {
	info, err := os.Stat(path)

	if err != nil {
		return "", fmt.Errorf("file not exists: %s", path)
	}

	if info.IsDir() {
		return "", fmt.Errorf("expected file, not a directory")
	}

	contentType, err := detectFileType(path)
	if err != nil {
		return "", err
	}

	if contentType != "application/json" && contentType != "text/plain; charset=utf-8" {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("invalid json: %w", err)
	}

	return fmt.Sprintf("parsed: %v", result), nil
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
