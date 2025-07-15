package main

import (
	"fmt"
	"strconv"
	"path/filepath"
	// "encoding/hex"
	// "encoding/json"
	"net/http"
	"mime/multipart"
	"os"
)

// Il faut créer un dossier qui contiendra toutes mes images, s'il n'existe pas. 
// Regarder si le path mon image n'existe pas déjà.
// Générer un nom d'image puis créer le fichier sur le serveur.
// Copier le contenu de mon image dans le file. 

func createDirectory() error {
	path := "photosDirectory"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
        return fmt.Errorf("directory already exists: %s", path);
    }
	err := os.MkdirAll(path, 0755)
    if err != nil {
		return fmt.Errorf("failed to create directory: %v", err);
	}
	return (nil);
}

func createImage(file multipart.File, id int, timeStamp string) {
	fmt.Println(file, id, timeStamp);
	filename := fmt.Sprintf("user_%d_%d_%s", tmpId, timestamp, fileHeader.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	// Créer le fichier sur le serveur
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(writer, "Erreur création fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copier le contenu de l'image dans le fichier
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(writer, "Erreur sauvegarde image: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) downloadImage(writer http.ResponseWriter, request *http.Request) {
	// var image Pictures;

	fmt.Println(Yellow + "Download image" + Reset)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    err := request.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(writer, "Erreur parsing formulaire: "+err.Error(), http.StatusBadRequest)
        return
    }
    file, _, err := request.FormFile("image")
    if err != nil {
        http.Error(writer, "Erreur récupération image: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
	timeStamp := request.FormValue("timestamp");
	userId := request.FormValue("id");
	tmpId, _ := strconv.Atoi(userId);
	fmt.Println(timeStamp, userId);
	createImage(file, tmpId, timeStamp);
	err = app.dataBase.QueryRow("INSERT INTO images (image_path, userId, uploadTime) VALUES ($1, $2, $3) RETURNING id", "path", tmpId, timeStamp).Scan(&userId);
	if err != nil {
		fmt.Println(Red + "Error : insert image in DB" + Reset);
		http.Error(writer, err.Error(), http.StatusInternalServerError);
		return;
	}
}
