package adapter

import (
	"github.com/homka122/Gomka122/collector/internal/domain"
	"github.com/homka122/Gomka122/internal/github"
)

type GithubRepositoryAdapter struct {
	GithubClient github.Client
}

func NewGithubRepositoryAdapter(client github.Client) GithubRepositoryAdapter {
	return GithubRepositoryAdapter{GithubClient: client}
}

func (gra GithubRepositoryAdapter) GetRepository(owner, repo string) (domain.GithubRepository, error) {
	dto, err := gra.GithubClient.GetRepository(owner, repo)
	if err != nil {
		return domain.GithubRepository{}, err
	}

	return domain.GithubRepository{
		Name:        dto.Name,
		Description: dto.Description,
		Stars:       int32(dto.Stars),
		Forks:       int32(dto.Forks),
		Create_date: dto.CreationDate,
	}, nil
}
