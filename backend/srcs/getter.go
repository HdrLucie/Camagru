package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
)

// ! ||--------------------------------------------------------------------------------||
// ! ||                                  USER GETTERS                                  ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) UserExists(username string) (error) {
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Check if an user exists in DB" + Reset)
	}
	var err error
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`
	app.dataBase.QueryRow(query, username).Scan(&err)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (app *App) getUserByJWT(JWT string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by JWT" + Reset)
	}
	query := "SELECT id, email, username, password, JWT, authToken, authStatus, avatar, notify FROM Users WHERE JWT = $1"
	row := app.dataBase.QueryRow(query, JWT)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.JWT, &user.AuthToken, &user.AuthStatus, &user.Avatar, &user.Notify)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (app *App) getStatus(id int) (bool) {
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get Status" + Reset)
	}	
	for i, _ := range app.users {
		if app.users[i].Id == id {
			return app.users[i].AuthStatus
		}
	}
	return false
}

func (app *App) getUserByUsername(username string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by username" + Reset)
	}
	query := "SELECT id, email, username, password, authToken, authStatus, notify FROM Users WHERE username = $1"
	row := app.dataBase.QueryRow(query, username)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.AuthStatus, &user.Notify)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (app *App) getUserByEmail(email string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by email" + Reset)
	}
	query := "SELECT id, email, username, password, authToken, authStatus FROM Users WHERE email = $1"
	row := app.dataBase.QueryRow(query, email)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.AuthStatus)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                  FRONT GETTERS                                 ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) deserializeUserData(writer http.ResponseWriter, request *http.Request) User {
	var u User

	fmt.Println(Yellow + "deserializeUserData function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err := json.NewDecoder(request.Body).Decode(&u)
	fmt.Println(u)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	return u
}

func extractJWTFromRequest(request *http.Request) string {
	JWT := request.Header.Get("Authorization")
	return strings.TrimPrefix(JWT, "Bearer ")
}

func (app *App) getUser(writer http.ResponseWriter, request *http.Request) {
	token := extractJWTFromRequest(request)
	user, _ := app.getUserByJWT(token)
	user.Password = ""
	user.AuthToken = ""
	fmt.Println(user);
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(user)
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                 GETTER STICKERS                                ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) getStickers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(app.stickers)
}

func (app *App) getStickerById(writer http.ResponseWriter, request *http.Request) {
	var sticker Stickers;
    if request.Method != http.MethodGet {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    path := strings.TrimPrefix(request.URL.Path, "/getSticker/")
    id, err := strconv.Atoi(path)
	fmt.Println(id);
    if err != nil {
        http.Error(writer, "Invalid ID", http.StatusBadRequest)
        return
    }
	query := "SELECT name, image_path FROM Stickers WHERE id = $1"
	row := app.dataBase.QueryRow(query, id)
	err = row.Scan(&sticker.Name, &sticker.Path)
	fmt.Println(sticker)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sticker.Id, sticker.Name)
	json.NewEncoder(writer).Encode(sticker);
}

func (app *App)	getStickerPathById(id int) string {
	for i, _ := range app.stickers {
		if app.stickers[i].Id == id {
			return app.stickers[i].Path;
		}
	}
	return "";
}

func (app *App) getAvatars(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(app.stickers)
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                 PICTURES GETTER                                ||
// ! ||--------------------------------------------------------------------------------||


// ! ||--------------------------------------------------------------------------------||
// ! ||                            HELPER FUNCTIONS - DB QUERIES                       ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) getAllPictures(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Content-Type", "application/json")
    if request.Method != http.MethodGet {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    pictures, err := app.fetchAllPicturesFromDB()
    if err != nil {
        fmt.Println(Red + "Error fetching pictures: " + err.Error() + Reset)
        http.Error(writer, "Error fetching pictures", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(writer).Encode(pictures)
}

func (app *App) fetchAllPicturesFromDB() ([]Pictures, error) {
    var pictures []Pictures
    
    query := "SELECT id, image_path, userId, uploadTime FROM images ORDER BY uploadTime DESC"
    rows, err := app.dataBase.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var pic Pictures
        err := rows.Scan(&pic.Id, &pic.Path, &pic.userId, &pic.uploadTime)
        if err != nil {
            return nil, err
        }
        pictures = append(pictures, pic)
    }
    
    return pictures, nil
}
