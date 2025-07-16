package main

import (
	"fmt"
	"strconv"
	"path/filepath"
	// "encoding/hex"
	// "encoding/json"
	"net/http"
	"io"
	"time"
	// "mime/multipart"
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
    file, fileHeader, err := request.FormFile("image")
    if err != nil {
        http.Error(writer, "Erreur récupération image: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close();
	timeStamp := request.FormValue("timestamp")
	userId := request.FormValue("id")
	tmpId, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(writer, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			http.Error(writer, "Erreur création dossier: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("user_%d_%d_%s", tmpId, timestamp, fileHeader.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(writer, "Erreur création fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(writer, "Erreur sauvegarde image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var imageId int
	err = app.dataBase.QueryRow("INSERT INTO images (image_path, userId, uploadTime) VALUES ($1, $2, $3) RETURNING id", 
		filepath, tmpId, timeStamp).Scan(&imageId)
	if err != nil {
		fmt.Println(Red + "Error : insert image in DB" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// response := map[string]interface{}{
	// 	"success": true,
	// 	"message": "Image sauvegardée avec succès",
	// 	"imageId": imageId,
	// 	"path": filepath,
	// }
	
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, `{"success": true, "message": "Image sauvegardée", "imageId": %d, "path": "%s"}`, imageId, filepath)

	fmt.Println(Green + "Image sauvegardée: " + filepath + Reset)
	err = app.dataBase.QueryRow("INSERT INTO images (image_path, userId, uploadTime) VALUES ($1, $2, $3) RETURNING id", "path", tmpId, timeStamp).Scan(&userId);
	if err != nil {
		fmt.Println(Red + "Error : insert image in DB" + Reset);
		http.Error(writer, err.Error(), http.StatusInternalServerError);
		return;
	}
}
