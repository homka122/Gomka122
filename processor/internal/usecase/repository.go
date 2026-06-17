package usecase

import "github.com/homka122/Gomka122/processor/internal/domain"

type Collector interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type RepositoryUseCase struct {
	Collector Collector
}

func NewRepositoryUsecase(collector Collector) *RepositoryUseCase {
	return &RepositoryUseCase{Collector: collector}
}

func (r *RepositoryUseCase) GetRepository(owner, repo string) (domain.Repository, error) {
	return r.Collector.GetRepository(owner, repo)
}
