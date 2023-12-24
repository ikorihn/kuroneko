package controller

import "net/http"

type Controller struct {
	httpClient *http.Client

	favorites Favorite
}

func NewController() *Controller {
	return &Controller{
		httpClient: http.DefaultClient,
	}
}
