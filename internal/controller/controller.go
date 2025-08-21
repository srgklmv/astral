package controller

type controller struct {
	documentsUsecase documentsUsecase
	authUsecase      authUsecase
}

type usecase interface {
	authUsecase
	documentsUsecase
}

func New(usecase usecase) *controller {
	return &controller{
		documentsUsecase: usecase,
		authUsecase:      usecase,
	}
}
