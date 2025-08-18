package controller

type controller struct{}

type service interface {
	authService
	docsService
}

func New() *controller {
	return &controller{}
}
