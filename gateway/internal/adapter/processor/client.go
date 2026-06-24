package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/homka122/Gomka122/gateway/internal/config"
	"github.com/homka122/Gomka122/gateway/internal/domain"
	apperror "github.com/homka122/Gomka122/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbProcessor "github.com/homka122/Gomka122/proto/processor"
)

type Processor struct {
	conn   *grpc.ClientConn
	client pbProcessor.ProcessorServiceClient
	log    *slog.Logger
}

func NewProcessor(cfg config.Config, log *slog.Logger) Processor {
	conn, error := grpc.NewClient(cfg.ProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}

	client := pbProcessor.NewProcessorServiceClient(conn)

	return Processor{
		conn:   conn,
		client: client,
		log:    log,
	}
}

func (p Processor) GetRepository(owner, repoName string) (*domain.Repository, error) {
	res, err := p.client.GetRepository(context.Background(), &pbProcessor.GetRepositoryRequest{Owner: owner, Repo: repoName})
	if err != nil {
		return nil, apperror.FromGRPC(err, err.Error())
	}

	repo := res.Repository

	switch repo.Status {
	case pbProcessor.RepositoryStatus_REPOSITORY_STATUS_READY:
		return &domain.Repository{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreateDate:  repo.CreateDate.AsTime(),
		}, nil
	case pbProcessor.RepositoryStatus_REPOSITORY_STATUS_NOT_FOUND:
		return nil, apperror.New(apperror.CodeNotFound, "repo not found")
	case pbProcessor.RepositoryStatus_REPOSITORY_STATUS_PREPARING:
		return nil, nil
	default:
		return nil, apperror.New(apperror.CodeInternal, "get repo internal error")
	}
}

func (p Processor) GetSubscribedRepository() ([](*domain.Repository), error) {
	repos, error := p.client.GetSubscribedRepository(context.Background(), &pbProcessor.GetSubscribedRepositoryRequest{})
	if error != nil {
		return nil, apperror.FromGRPC(error, "processor get subscribed repository")
	}

	result := make([](*domain.Repository), len(repos.Repositories))
	for k, repo := range repos.Repositories {
		switch repo.Status {
		case pbProcessor.RepositoryStatus_REPOSITORY_STATUS_READY:
			result[k] = &domain.Repository{
				Name:        repo.Name,
				Description: repo.Description,
				Stars:       repo.Stars,
				Forks:       repo.Forks,
				CreateDate:  repo.CreateDate.AsTime(),
			}

			continue
		case pbProcessor.RepositoryStatus_REPOSITORY_STATUS_PREPARING:
			result[k] = nil
			continue
		default:
			panic(fmt.Sprintf("unexpected error %v", repo))
		}
	}

	return result, nil
}

func (p Processor) Ping() (string, error) {
	pong, err := p.client.Ping(context.Background(), &pbProcessor.PingRequest{})
	if err != nil {
		return "", apperror.New(apperror.CodeUnavailable, "processor unvailable")
	}

	return pong.Reply, nil
}

func (p Processor) Close() {
	p.conn.Close()
}
