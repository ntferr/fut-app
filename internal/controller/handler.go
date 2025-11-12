package controller

type Controller struct {
	Auth Auth
}

func NewController(auth Auth) *Controller {
	return &Controller{
		Auth: auth,
	}
}
