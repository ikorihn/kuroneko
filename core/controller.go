package core

import "net/http"

type Controller struct {
	httpClient *http.Client

	Favorites Favorite
}

func NewController() (*Controller, error) {
	favorites, err := loadFavorite()
	if err != nil {
		return nil, err
	}

	return &Controller{
		httpClient: http.DefaultClient,
		Favorites:  favorites,
	}, nil
}
