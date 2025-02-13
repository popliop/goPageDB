package upload

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

type FileUpload struct {
	FileName   string
	LoadedTime time.Time
	Data       [][]string
}

func NewFileUpload(fileName string, loadedTime time.Time, data [][]string) FileUpload {
	return FileUpload{
		FileName:   fileName,
		LoadedTime: loadedTime,
		Data:       data,
	}
}

func HandleForm(r *http.Request) error {
	file, header, err := ParseForm(r)
	if err != nil {
		return fmt.Errorf("handleForm: %w", err)
	}
	defer file.Close()

	if err := ValidateForm(header); err != nil {
		return fmt.Errorf("validateForm: %w", err)
	}

	data, err := ReadFile(file)
	if err != nil {
		return fmt.Errorf("readFile: %w", err)
	}

	NewFileUpload(header.Filename, time.Now(), data)

	return nil
}

func ReadFile(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	data := [][]string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		data = append(data, record)
	}
	return data, nil
}

func ParseForm(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	// Limit the uploaded file size to 10MB.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, nil, fmt.Errorf("error parsing multipart form: %w", err)
	}

	file, header, err := r.FormFile("datafile")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve form file from request: %w", err)
	}
	return file, header, nil
}

func ValidateForm(header *multipart.FileHeader) error {
	if header.Size == 0 {
		return fmt.Errorf("uploaded file is empty")
	}
	if ext := filepath.Ext(header.Filename); ext != ".csv" {
		return fmt.Errorf("invalid extension: expected .csv, got %s", ext)
	}
	return nil
}

func GetFilename(fu FileUpload) (string, error) {
	if fu.FileName == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	return fu.FileName, nil
}
