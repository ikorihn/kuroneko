package controller

import (
	"os"
	"os/exec"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

const favoritesFile = "kuroneko/favorites.toml"

// EditBody edits request body using $EDITOR.
// If no EDITOR is specified, vim will open.
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
