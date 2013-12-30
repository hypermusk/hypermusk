package apiimpl

import (
	"github.com/hypermusk/hypermusk/tests/api"
)

type Service struct{}

func (s *Service) Authorize(name string) (api.UseMap, error) {
	return api.UseMap{}, nil
}
