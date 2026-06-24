package usecase

import (
	"context"
	"log/slog"

	"github.com/homka122/Gomka122/gateway/internal/domain"
)

type Subscriber interface {
	Subscribe(ctx context.Context, owner, repo string) error
	Unsubscribe(ctx context.Context, owner, repo string) error
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type SubscribeUseCase struct {
	Subscriber Subscriber
}

func NewSubscribeUseCase(subscriber Subscriber, log *slog.Logger) *SubscribeUseCase {
	return &SubscribeUseCase{Subscriber: subscriber}
}

func (u *SubscribeUseCase) Subscribe(ctx context.Context, owner, repo string) error {
	return u.Subscriber.Subscribe(ctx, owner, repo)
}

func (u *SubscribeUseCase) Unsubscribe(ctx context.Context, owner, repo string) error {
	return u.Subscriber.Unsubscribe(ctx, owner, repo)
}

func (u *SubscribeUseCase) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	return u.Subscriber.GetSubscriptions(ctx)
}
