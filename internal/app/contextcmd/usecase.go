package contextcmd

import (
	"slices"
	"strings"
)

type Entry struct {
	RootPath   string
	LastUsedAt int64
}

type Port interface {
	ResolveCurrentRoot(cwd string) (string, error)
	ListEntries() ([]Entry, error)
	ResolveUseRoot(raw string) (string, error)
	WriteCurrent(root string) error
}

type Service struct {
	port Port
}

func NewService(port Port) *Service {
	return &Service{port: port}
}

func (s *Service) Current(cwd string) (string, error) {
	return s.port.ResolveCurrentRoot(cwd)
}

func (s *Service) List() ([]Entry, error) {
	entries, err := s.port.ListEntries()
	if err != nil {
		return nil, err
	}
	slices.SortFunc(entries, func(a, b Entry) int {
		if a.LastUsedAt != b.LastUsedAt {
			if a.LastUsedAt > b.LastUsedAt {
				return -1
			}
			return 1
		}
		return strings.Compare(a.RootPath, b.RootPath)
	})
	return entries, nil
}

func (s *Service) Use(raw string) (string, error) {
	root, err := s.port.ResolveUseRoot(raw)
	if err != nil {
		return "", err
	}
	if err := s.port.WriteCurrent(root); err != nil {
		return "", err
	}
	return root, nil
}
