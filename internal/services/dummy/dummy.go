package dummy

import (
	"context"
	"github.com/fredmayer/go-rest-api-template/internal/domain/models"
)

type DummyService struct {
	provider DummyProvider
}

type DummyProvider interface {
	GetById(ctx context.Context, id int) (models.Dummy, error)
}

func New(provider DummyProvider) *DummyService {
	return &DummyService{
		provider,
	}
}

// Найти запись
func (s *DummyService) Find(id int) (models.Dummy, error) {
	return s.provider.GetById(context.Background(), id)
}
