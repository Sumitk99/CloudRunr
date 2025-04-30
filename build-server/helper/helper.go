package helper

import (
	"fmt"
	"net/http"
	"os"
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
