package helper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetFileType(file *os.File) *string {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	contentType := http.DetectContentType(buffer)
	return &contentType
}

func GetFilePaths(DistFolderPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(DistFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return files, nil
}
