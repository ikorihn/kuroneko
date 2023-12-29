package core

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

type Favorite struct {
	Request []Request `toml:"request"`
}

const favoritesFile = "kuroneko/favorites.toml"

// loadFavorite loads favorite request from file
func loadFavorite() (Favorite, error) {
	favorite := Favorite{
		Request: make([]Request, 0),
	}
	f, err := xdg.SearchDataFile(favoritesFile)
	if err != nil {
		return favorite, nil
	}

	b, err := os.ReadFile(f)
	if err != nil {
		return favorite, err
	}

	if err := toml.Unmarshal(b, &favorite); err != nil {
		return favorite, err
	}

	return favorite, nil
}

// SaveFavorite adds favorite request and save to file
func (c *Controller) SaveFavorite(request []Request) error {

	favoritesFile, err := xdg.DataFile(favoritesFile)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(favoritesFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c.Favorites.Request = request

	b, err := toml.Marshal(c.Favorites)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
