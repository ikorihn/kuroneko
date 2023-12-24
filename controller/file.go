package controller

import (
	"os"
	"os/exec"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

const favoritesFile = "kuroneko/favorites.toml"

func (c *Controller) EditBody() ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	tempFile, err := os.CreateTemp("", "kuroneko-*")
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	defer os.Remove(tempFile.Name())

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	body, err := os.ReadFile(tempFile.Name())

	return body, err
}

func (c *Controller) LoadFavorite() (Favorite, error) {
	c.favorites = Favorite{
		Request: make([]Request, 0),
	}
	f, err := xdg.SearchDataFile(favoritesFile)
	if err != nil {
		return c.favorites, nil
	}

	b, err := os.ReadFile(f)
	if err != nil {
		return c.favorites, err
	}

	if err := toml.Unmarshal(b, &c.favorites); err != nil {
		return c.favorites, err
	}

	return c.favorites, nil
}

func (c *Controller) SaveFavorite(request Request) error {

	favoritesFile, err := xdg.DataFile(favoritesFile)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(favoritesFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c.favorites.Request = append(c.favorites.Request, request)
	b, err := toml.Marshal(c.favorites)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil

}
