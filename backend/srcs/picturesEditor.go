package main

import (
	"fmt"
	"image"
	"image/png"
	"image/color"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"math"
	"time"
	"image/jpeg"
	"strings"
	"golang.org/x/image/draw"
)

// Il faut créer un dossier qui contiendra toutes mes images, s'il n'existe pas.
// Regarder si le path mon image n'existe pas déjà.
// Générer un nom d'image puis créer le fichier sur le serveur.
// Copier le contenu de mon image dans le file.

type Pixel struct {
	Point image.Point
	Color color.Color
}

func (app *App) createUploadsDirectory() error {
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
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

func (app *App) resizeImg(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(Red + "Resize" + Reset);

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
	defer file.Close();
	img, _, err := image.Decode(file)
    if err != nil {
        http.Error(writer, "Erreur décodage image: "+err.Error(), http.StatusBadRequest)
        return
    }
	fmt.Println("Before resize :" , img.Bounds().Dx() , img.Bounds().Dy());
    width := int(float64(img.Bounds().Dx()))
    height := int(float64(img.Bounds().Dy()))
	var newWidth, newHeight int
	if width > 500 || height > 700 {
		ratioW := 500.0 / float64(width)
		ratioH := 700.0 / float64(height)
		r := math.Min(ratioW, ratioH)
		newWidth = int(float64(width) * r)
		newHeight = int(float64(height) * r)
		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.BiLinear.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)	
		fmt.Println("After resize :" , resized.Bounds().Dx() , resized.Bounds().Dy());
		writer.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(writer, resized, &jpeg.Options{Quality: 85})
	}
	writer.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(writer, img, &jpeg.Options{Quality: 85})

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
	fmt.Println("Size :", img.Bounds());
	draw.Draw(finImage, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	stickerPos := image.Point{posX, posY}
	stickerRect := sticker.Bounds();

	w := float64(stickerRect.Dx());
	h := float64(stickerRect.Dy());

	if w > h {
		h = h / w * 128;
		w = 128;
	} else {
		w = w / h * 128;
		h = 128;
	}
	stickerRect.Max.X = int(w);
	stickerRect.Max.Y = int(h);
	stickerRect.Max = stickerRect.Max.Add(stickerPos);
	stickerRect.Min = stickerRect.Min.Add(stickerPos);

	draw.BiLinear.Scale(finImage, stickerRect, sticker, sticker.Bounds(), draw.Over, nil)
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

type S struct {
	Path	string `json:"path"`
	PosX	int    `json:"posX"`
	PosY	int    `json:"posY"`
	Id		int	`json:"id"`
}

func (app *App) downloadImage(writer http.ResponseWriter, request *http.Request) {
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
	stickersJSON := request.FormValue("stickers")
	var stickers []S
	err = json.Unmarshal([]byte(stickersJSON), &stickers)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(stickers);
	timeStamp := request.FormValue("timestamp")
	userId := request.FormValue("id")
	tmpId, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(writer, "userId invalide: "+err.Error(), http.StatusBadRequest)
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
	fmt.Println("createImage result:", filepath, err)
	for _, sticker := range stickers {
		path := app.getStickerPathById(sticker.Id)
		concatImage(filepath, path, sticker.PosX, sticker.PosY);
	}
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
