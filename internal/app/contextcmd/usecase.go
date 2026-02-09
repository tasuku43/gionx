package contextcmd

import (
	"slices"
	"strings"
)

type Entry struct {
	ContextName string
	RootPath    string
	LastUsedAt  int64
}

type Port interface {
	ResolveCurrentRoot(cwd string) (string, error)
	ResolveCurrentName(root string) (string, bool, error)
	ListEntries() ([]Entry, error)
	ResolveUseRootByName(name string) (string, bool, error)
	RenameContext(oldName string, newName string) (string, error)
	RemoveContext(name string) (string, error)
	CreateContext(name string, rawPath string) (string, error)
	WriteCurrent(root string) error
}

type Service struct {
	port Port
}

func NewService(port Port) *Service {
	return &Service{port: port}
}

func (s *Service) Current(cwd string) (string, error) {
	root, err := s.port.ResolveCurrentRoot(cwd)
	if err != nil {
		return "", err
	}
	if name, ok, err := s.port.ResolveCurrentName(root); err == nil && ok {
		return name, nil
	}
	return root, nil
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

func (s *Service) Use(name string) (string, error) {
	root, ok, err := s.port.ResolveUseRootByName(name)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errContextNotFound(name)
	}
	if err := s.port.WriteCurrent(root); err != nil {
		return "", err
	}
	return root, nil
}

func (s *Service) ResolveRootByName(name string) (string, bool, error) {
	return s.port.ResolveUseRootByName(name)
}

func (s *Service) Rename(oldName string, newName string) (string, error) {
	root, err := s.port.RenameContext(oldName, newName)
	if err != nil {
		return "", err
	}
	return root, nil
}

func (s *Service) Remove(name string) (string, error) {
	root, err := s.port.RemoveContext(name)
	if err != nil {
		return "", err
	}
	return root, nil
}

func (s *Service) Create(name string, rawPath string, useNow bool) (string, error) {
	root, err := s.port.CreateContext(name, rawPath)
	if err != nil {
		return "", err
	}
	if useNow {
		if err := s.port.WriteCurrent(root); err != nil {
			return "", err
		}
	}
	return root, nil
}

func errContextNotFound(name string) error {
	return &contextNotFoundError{Name: strings.TrimSpace(name)}
}

type contextNotFoundError struct {
	Name string
}

func (e *contextNotFoundError) Error() string {
	return "context not found: " + e.Name
}
