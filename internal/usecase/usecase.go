package usecase

type repository interface {
	userRepository
}

type userRepository interface {
	IsLoginTaken(login string) (bool, error)
}

type usecase struct {
	userRepository userRepository
}

func New(repos) *usecase {
	return &usecase{}
}
