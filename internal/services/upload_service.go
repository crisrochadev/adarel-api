package services

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UploadService interface {
	SaveImage(file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	basePath string
}

func NewUploadService(basePath string) UploadService {
	return &uploadService{basePath: basePath}
}

func (s *uploadService) SaveImage(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("file is required")
	}
	if file.Size > 2*1024*1024 {
		return "", errors.New("file too large")
	}

	src, err := file.Open()
	if err != nil {
		return "", errors.New("could not open file")
	}
	defer src.Close()

	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		return "", errors.New("could not read file")
	}
	if _, err := src.Seek(0, 0); err != nil {
		return "", errors.New("could not process file")
	}

	mimeType := http.DetectContentType(buff)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", errors.New("invalid file type")
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	name := uuid.NewString() + ext
	fullPath := filepath.Join(s.basePath, name)
	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return "", errors.New("could not prepare upload directory")
	}

	if err := saveMultipartFile(file, fullPath); err != nil {
		return "", errors.New("could not save file")
	}

	return "/uploads/" + name, nil
}

func saveMultipartFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}
