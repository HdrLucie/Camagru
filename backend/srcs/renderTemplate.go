package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"context"
	"strings"
	"strconv"
	"os"
)

func serveTemplate(templateName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		templateDir := frontendPath("templates")

        tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
        if err != nil {
            http.Error(w, "Could not parse templates", http.StatusInternalServerError)
            fmt.Println("Error parsing templates:", err)
            return
        }

        data := TemplateData{Page: templateName}
        err = tmpl.ExecuteTemplate(w, templateName, data)
        if err != nil {
            http.Error(w, "Could not execute template", http.StatusInternalServerError)
            fmt.Println("Error executing template:", err)
        }
    }
}

func frontendPath(subPath string) string {
    base := os.Getenv("FRONTEND_PATH")
    if base == "" {
        base = "../../frontend/srcs"
    }
    return filepath.Join(base, subPath)
}

func	serveStyleFiles(router *http.ServeMux) {
	styles := http.FileServer(http.Dir(frontendPath("stylesheets")))
	router.Handle("/styles/", http.StripPrefix("/styles", styles))
}

func serveScriptsFiles(router *http.ServeMux) {
	scripts := http.FileServer(http.Dir(frontendPath("scripts")))
	router.Handle("/scripts/", http.StripPrefix("/scripts", scripts))
}

func serveImgFiles(router *http.ServeMux) {
	assets := http.FileServer(http.Dir(frontendPath("assets/img")))
	router.Handle("/assets/", http.StripPrefix("/assets", assets))
}

func serveStickersFiles(router *http.ServeMux) {
	stickers := http.FileServer(http.Dir(frontendPath("assets/stickers")))
	router.Handle("/stickers/", http.StripPrefix("/stickers", stickers))
}

func servePicturesFiles(router *http.ServeMux) {
	stickers := http.FileServer(http.Dir("./uploads/"))
	router.Handle("/uploads/", http.StripPrefix("/uploads", stickers))
}

func (app *App) verifyJWT(JWT string) (*User, error) {
	user, error := app.getUserByJWT(JWT)
	if error != nil {
		return nil, error
	}
	return user, error
}

func (app *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		JWT := extractJWTFromRequest(request)
		if JWT == "" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := app.verifyJWT(JWT)
		if err != nil {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), "user", user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	}
}

func (app *App) viewPhoto(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    id := strings.TrimPrefix(path, "/photo/")
    if id == "" || id == "/" {
        http.Error(w, "Photo ID required", http.StatusBadRequest)
        return
    }
    photoID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid photo ID", http.StatusBadRequest)
        return
    }
    picture, err := app.getPictureById(photoID)
    if err != nil {
        http.Error(w, "Photo not found", http.StatusNotFound)
        fmt.Println("Error getting picture:", err)
        return
    }
    data := struct {
        Page    string
        Picture *Pictures
	}{
        Page:    "photo.html",
        Picture: picture,
    }
	templateDir := frontendPath("templates")

    tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
    if err != nil {
        http.Error(w, "Could not parse templates", http.StatusInternalServerError)
        fmt.Println("Error parsing templates:", err)
        return
    }
    err = tmpl.ExecuteTemplate(w, "photo.html", data)
    if err != nil {
        http.Error(w, "Could not execute template", http.StatusInternalServerError)
        fmt.Println("Error executing template:", err)
    }
}

func (app *App) router(router *http.ServeMux) {
	router.HandleFunc("/", serveTemplate("firstPage.html"))
	router.HandleFunc("/connection", serveTemplate("login.html"))
	router.HandleFunc("/gallery/", serveTemplate("gallery.html"))
	router.HandleFunc("/takePicture", serveTemplate("take-picture.html"))
	router.HandleFunc("/loadImage", serveTemplate("loadImage.html"))
	router.HandleFunc("/authentification", serveTemplate("authentification.html"))
	router.HandleFunc("/forgetPassword", serveTemplate("forgetPassword.html"))
	router.HandleFunc("/verify", serveTemplate("verify.html"))
	router.HandleFunc("/resetPassword", serveTemplate("resetPassword.html"))
	router.HandleFunc("/profile", serveTemplate("profile.html"))
	router.HandleFunc("/photo/", app.viewPhoto)
	router.HandleFunc("/signUp", app.signUp)
	router.HandleFunc("/login", app.login)
	router.HandleFunc("/sendImage", app.authMiddleware(app.downloadImage))
	router.HandleFunc("/logout", app.authMiddleware(app.logout))
	router.HandleFunc("/verifyAccount", app.verifyAccount)
	router.HandleFunc("/sendResetLink", app.sendResetLink)
	router.HandleFunc("/newPassword", app.resetPassword)
	router.HandleFunc("/setUserDatas", app.authMiddleware(app.modifyProfile))
	router.HandleFunc("/editPassword", app.authMiddleware(app.modifyPassword))
	router.HandleFunc("/getUser", app.authMiddleware(app.getUser))
	router.HandleFunc("/getStickers", app.authMiddleware(app.getStickers))
	router.HandleFunc("/getSticker/", app.authMiddleware(app.getStickerById))
	router.HandleFunc("/getPictures/", app.getPage)
	router.HandleFunc("/getAvatars", app.authMiddleware(app.getAvatars))
	router.HandleFunc("/sendLikes", app.sendLikes)
	router.HandleFunc("/getPicture/", app.authMiddleware(app.getPicture))
	router.HandleFunc("/sendComments", app.authMiddleware(app.manageComment))
	router.HandleFunc("/getComments/", app.authMiddleware(app.getComments));
	router.HandleFunc("/deleteImg", app.authMiddleware(app.deleteImg));
	router.HandleFunc("/getLikes/", app.authMiddleware(app.getLikes));
	router.HandleFunc("/resizeImg", app.resizeImg);
}

func renderTemplate(router *http.ServeMux, app *App) {
	serveStyleFiles(router)
	serveScriptsFiles(router)
	serveImgFiles(router)
	serveStickersFiles(router)
	servePicturesFiles(router)
	app.router(router)
}
