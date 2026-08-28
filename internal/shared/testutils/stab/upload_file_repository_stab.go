package stab

import (
	"io"

	repository_interface "growlink/internal/domain/interface/repository"
)

var _ repository_interface.UploadFileRepository = (*UploadFileRepository)(nil)

type UploadFileRepository struct {
	ExistsFunc     func(name string) (bool, error)
	NewFileFunc    func(name string, content io.Reader) error
	DeleteFileFunc func(path string) error

	NewFileNames    []string
	ExistsNames     []string
	DeleteFileNames []string
}

func (s *UploadFileRepository) Exists(name string) (bool, error) {
	s.ExistsNames = append(s.ExistsNames, name)
	if s.ExistsFunc == nil {
		return false, nil
	}
	return s.ExistsFunc(name)
}

func (s *UploadFileRepository) NewFile(name string, content io.Reader) error {
	s.NewFileNames = append(s.NewFileNames, name)
	if s.NewFileFunc == nil {
		return nil
	}
	return s.NewFileFunc(name, content)
}

func (s *UploadFileRepository) DeleteFile(path string) error {
	s.DeleteFileNames = append(s.DeleteFileNames, path)
	if s.DeleteFileFunc == nil {
		return nil
	}
	return s.DeleteFileFunc(path)
}
