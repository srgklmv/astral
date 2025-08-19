package controller

type controller struct{}

type usecase interface {
	authUsecase
	docsUsecase
}

func New() *controller {
	return &controller{}
}
