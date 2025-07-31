package main

import (
	"fmt"
	"image"
	"image/png"
	"image/draw"
	"image/color"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"image/jpeg"
	"strings"

)

// Il faut créer un dossier qui contiendra toutes mes images, s'il n'existe pas.
// Regarder si le path mon image n'existe pas déjà.
// Générer un nom d'image puis créer le fichier sur le serveur.
// Copier le contenu de mon image dans le file.

type Pixel struct {
	Point image.Point
	Color color.Color
}

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

func openAndDecode(filepath string) (image.Image, error) {
	file, err := os.Open(filepath);
	if err != nil {
		panic(err);
	}
	defer file.Close();
	img, _, err := image.Decode(file);
	if err != nil {
		panic(err);
	}
	return img, nil;
}

func createImage(file multipart.File, fileHeader *multipart.FileHeader, tmpId int, uploadsDir string) (string, error) {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("user_%d_%d_%s", tmpId, timestamp, fileHeader.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return filepath, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return filepath, err
	}
	return filepath, err;	
}

func DecodePixelsFromImage(img image.Image, offsetX, offsetY int) []*Pixel {
    pixels := []*Pixel{}
    for y := 0; y <= img.Bounds().Max.Y; y++ {
        for x := 0; x <= img.Bounds().Max.X; x++ {
            p := &Pixel{
                Point: image.Point{x + offsetX, y + offsetY},
                Color: img.At(x, y),
            }
            pixels = append(pixels, p)
        }
    }
    return pixels
}

func concatImage(imgPath string, stickerPath string, posX int, posY int) {
	img, err := openAndDecode(imgPath)
	if err != nil {
		panic(err)
	}
	
	sticker, err := openAndDecode(stickerPath)
	if err != nil {
		panic(err)
	}

	finImage := image.NewRGBA(img.Bounds())
	fmt.Println(posX, posY)	
	draw.Draw(finImage, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	stickerPos := image.Point{posX, posY}
	stickerRect := image.Rectangle{
		Min: stickerPos,
		Max: stickerPos.Add(sticker.Bounds().Size()),
	}
	
	draw.Draw(finImage, stickerRect, sticker, image.Point{0, 0}, draw.Over)

	out, err := os.Create(imgPath) 
	if err != nil {
		panic(err)
	}
	defer out.Close()
	
	ext := strings.ToLower(filepath.Ext(imgPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, finImage, &jpeg.Options{Quality: 95})
	case ".png":
		err = png.Encode(out, finImage)
	default:
		err = png.Encode(out, finImage)
	}
	
	if err != nil {
		panic(err)
	}
}

func (app *App) downloadImage(writer http.ResponseWriter, request *http.Request) {

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
	tmpStickerId := request.FormValue("imageId");
	posX := request.FormValue("posX");
	posY := request.FormValue("posY");
	x, err := strconv.Atoi(posX);
	y, err := strconv.Atoi(posY);
	tmpId, err := strconv.Atoi(userId)
	stickerId, err := strconv.Atoi(tmpStickerId);
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
		
	var imageId int
	filepath, err := createImage(file, fileHeader, tmpId, uploadsDir);
	stickerPath := app.getStickerPathById(stickerId);
	concatImage(filepath, stickerPath, x, y);
	err = app.dataBase.QueryRow("INSERT INTO images (image_path, userId, uploadTime) VALUES ($1, $2, $3) RETURNING id", 
		filepath, tmpId, timeStamp).Scan(&imageId)
	if err != nil {
		fmt.Println(Red + "Error : insert image in DB" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, `{"success": true, "message": "Image sauvegardée", "imageId": %d, "path": "%s"}`, imageId, filepath)

	fmt.Println(Green + "Image sauvegardée: " + filepath + Reset)
}
