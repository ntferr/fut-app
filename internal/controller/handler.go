package controller

type Controller struct {
	Auth     Auth
	Champion Champion
}

func NewController(auth Auth, champion Champion) *Controller {
	return &Controller{
		Auth:     auth,
		Champion: champion,
	}
}
