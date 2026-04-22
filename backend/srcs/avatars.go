package main

import (
	"fmt"
	"os"
)

func (app *App) extractPathAvatarsFiles() error {
	path := "../../frontend/srcs/assets/avatars/"
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	defer directory.Close()
	list, err := directory.Readdir(-1)
	if err != nil {
		return err
	}
	app.avatars = make([]Avatars, 0, len(list))
    for i, file := range list {
        avatar := Avatars{
            Id:   i,
            Name: file.Name(),
            Path: path + file.Name(),
        }
		app.avatars = append(app.avatars, avatar)
    }
	return nil
}

func (app *App) manageAvatarsInsertError(imagePath string) (bool, error) {
    var exists bool
    query := `
        SELECT EXISTS(
            SELECT 1 FROM avatars 
            WHERE image_path = $1
        )`
    
    err := app.dataBase.QueryRow(query, imagePath).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("erreur lors de la vérification du avatar: %v", err)
    }
    return exists, nil
}

func (app *App) InsertAvatars() error {
	var exists bool
	err := app.extractPathAvatarsFiles()
	if err != nil {
		return err
	}
	for i, avatar := range app.avatars {
		exists, err = app.manageAvatarsInsertError(avatar.Path)
		if (!exists) {
			fmt.Println(Red + "INSERT AVATAR")
			query := `
				INSERT INTO avatars (id, name, image_path)
				VALUES ($1, $2, $3)`
			_, err := app.dataBase.Exec(query, i, avatar.Name, avatar.Path)
			if err != nil {
				return fmt.Errorf("impossible d'insérer le avatar %s: %v", avatar.Name, err)
			}
		} else {
			fmt.Errorf("Avatars déjà existant %s: %v", avatar.Name, err)
		}
    }

	return nil
}
