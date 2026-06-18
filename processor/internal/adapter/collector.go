package collector

import (
	"context"

	apperror "github.com/homka122/Gomka122/internal/errors"
	"github.com/homka122/Gomka122/processor/internal/domain"
	"google.golang.org/grpc"

	pb "github.com/homka122/Gomka122/proto/collector"
)

type Collector struct {
	conn   *grpc.ClientConn
	client pb.CollectorServiceClient
}

func NewCollector(conn *grpc.ClientConn) Collector {
	return Collector{
		conn:   conn,
		client: pb.NewCollectorServiceClient(conn),
	}
}

func (c Collector) GetRepository(owner, repo string) (domain.Repository, error) {
	repository, err := c.client.GetRepository(context.Background(), &pb.GetRepositoryRequest{Owner: owner, Repo: repo})
	if err != nil {
		return domain.Repository{}, apperror.FromGRPC(err, "collector get repository")
	}

	return domain.Repository{
		Name:        repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreateDate:  repository.CreateDate.AsTime(),
	}, nil
}

func (c Collector) GetSubscribedRepository() ([]domain.Repository, error) {
	repos, err := c.client.GetSubscribedRepository(context.Background(), &pb.GetSubscribedRepositoryRequest{})
	if err != nil {
		return nil, apperror.FromGRPC(err, "collector get subscribed repository")
	}

	result := make([]domain.Repository, len(repos.Repositories))
	for k, repo := range repos.Repositories {
		result[k] = domain.Repository{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreateDate:  repo.CreateDate.AsTime(),
		}
	}

	return result, nil
}
