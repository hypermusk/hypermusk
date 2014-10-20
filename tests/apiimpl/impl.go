package apiimpl

import (
	"errors"
	"github.com/hypermusk/hypermusk/tests/api"
)

type Service struct{}

func (s *Service) Authorize(name string) (api.UseMap, error) {
	return api.UseMap{}, nil
}

func (s *Service) PermiessionDenied() (err error) {
	return errors.New("permission denied.")
}

var DefaultService = &Service{}
