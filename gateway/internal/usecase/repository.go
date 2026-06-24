package usecase

import (
	"context"
	"log/slog"

	"github.com/homka122/Gomka122/gateway/internal/domain"
)

type Processor interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscribedRepository(ctx context.Context) ([](*domain.Repository), error)
}

type RepositoryUseCase struct {
	Processor Processor
	log       *slog.Logger
}

func NewRepositoryUseCase(processor Processor, log *slog.Logger) *RepositoryUseCase {
	return &RepositoryUseCase{Processor: processor, log: log}
}

func (r *RepositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return r.Processor.GetRepository(ctx, owner, repo)
}

func (r *RepositoryUseCase) GetSubscribedRepository(ctx context.Context) ([](*domain.Repository), error) {
	return r.Processor.GetSubscribedRepository(ctx)
}
