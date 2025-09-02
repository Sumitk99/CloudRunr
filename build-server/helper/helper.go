package helper

import (
	"fmt"
	"github.com/Sumitk99/build-server/constants"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//func GetFileType(file *os.File) *string {
//	buffer := make([]byte, 512)
//	_, err := file.Read(buffer)
//	if err != nil {
//		fmt.Println("Error reading file:", err)
//		return nil
//	}
//	contentType := http.DetectContentType(buffer)
//	return &contentType
//}

func GetFileType(filePath string) *string {
	ext := strings.ToLower(filepath.Ext(filePath))
	mimeTypes := map[string]string{
		".js":    "application/javascript",
		".css":   "text/css",
		".json":  "application/json",
		".png":   "image/png",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".gif":   "image/gif",
		".svg":   "image/svg+xml",
		".ico":   "image/x-icon",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".eot":   "application/vnd.ms-fontobject",
		".html":  "text/html",
		".map":   "application/json",
	}
	if mime, ok := mimeTypes[ext]; ok {
		return &mime
	}

	// Fallback: detect from file content
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
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

func DetectBuildCommand(framework string) string {
	switch framework {
	case constants.REACT:
		return constants.REACT_BUILD_COMMAND
	case constants.ANGULAR:
		return constants.ANGULAR_BUILD_COMMAND
	}

	return ""
}
