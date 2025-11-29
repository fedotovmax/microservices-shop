package app

import "context"

type lifesycleService interface {
	Stop(ctx context.Context) error
}

type service struct {
	name string
	lifesycleService
}

func (s *service) Stop(ctx context.Context) error {
	return s.lifesycleService.Stop(ctx)
}

func (s *service) Name() string {
	return s.name
}

func newService(name string, ls lifesycleService) *service {
	return &service{
		name:             name,
		lifesycleService: ls,
	}
}
