package main

import (
	"fmt"
	// "log"
	"os"
)

func (app *App) extractPathFiles() error {
	path := "../../frontend/srcs/assets/stickers/"
	directory, err := os.Open(path)
	if err != nil {
		fmt.Println(Red + "HERE 1" + Reset)
		return err
	}
	defer directory.Close()
	list, err := directory.Readdir(-1)
	if err != nil {
		fmt.Println(Red + "HERE 2" + Reset)
		return err
	}
	app.stickers = make([]Stickers, 0, len(list))
    for i, file := range list {
        sticker := Stickers{
            Id:   i,
            Name: file.Name(),
            Path: path + file.Name(),
        }
		app.stickers = append(app.stickers, sticker)
    }
	return nil
}

func (app *App) manageStickersInsertError(imagePath string) (bool, error) {
    var exists bool
    query := `
        SELECT EXISTS(
            SELECT 1 FROM stickers 
            WHERE image_path = $1
        )`
    
    err := app.dataBase.QueryRow(query, imagePath).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("erreur lors de la vérification du sticker: %v", err)
    }
    return exists, nil
}

func (app *App) InsertSticker() error {
	var exists bool
	err := app.extractPathFiles()
	if err != nil {
		return err
	}
	for i, sticker := range app.stickers {
		exists, err = app.manageStickersInsertError(sticker.Path)
		if (!exists) {
			fmt.Println("Insertion")
			query := `
				INSERT INTO stickers (id, name, image_path)
				VALUES ($1, $2, $3)`
			_, err := app.dataBase.Exec(query, i, sticker.Name, sticker.Path)
			if err != nil {
				return fmt.Errorf("impossible d'insérer le sticker %s: %v", sticker.Name, err)
			}
		} else {
			fmt.Errorf("Stickers déjà existant %s: %v", sticker.Name, err)
		}
    }



	return nil
}