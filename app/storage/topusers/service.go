package topusers

type Service struct {
	repo *Repository
}

func (this *Service) FindAll(params Params) ([]*TopUser, error) {
	return this.repo.FindAll(params)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
