package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
)

type	cInfo struct {
	PId	int	`json:"pId"`
}
type cResponse struct {
	Username	string	`json:"Username"`
	Comment		string	`json:"Comment"`
	PId			int		`json:"pId"`

}

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

func (app *App) getUserById(id int) (*User, error) {
    var user User
    query := "SELECT id, email, username, password, authToken, authStatus, avatar, notify FROM Users WHERE id = $1"
    err := app.dataBase.QueryRow(query, id).Scan(
        &user.Id, &user.Email, &user.Username, &user.Password,
        &user.AuthToken, &user.AuthStatus, &user.Avatar, &user.Notify,
    )
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

func (app *App) getUserByPhotoId(id int) (*User, error) {
	var userId int64;
	var user	User;

	query := "SELECT userId FROM images WHERE id = $1"
	row := app.dataBase.QueryRow(query, id)
	err := row.Scan(&userId)
	if err != nil {
		fmt.Println("Here")
		fmt.Println(err)
		return nil, err
	}
	query = "SELECT id, email, username, password, authToken, authStatus, notify FROM Users WHERE id = $1"
	row = app.dataBase.QueryRow(query, userId);
	err = row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.AuthStatus, &user.Notify);
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil;
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
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(user)
}

func (app *App) getComments(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Content-Type", "application/json")
    parts := strings.Split(request.URL.Path, "/getComments/")
    pId, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        http.Error(writer, "Invalid photo ID", http.StatusBadRequest)
        return
    }
	query := `SELECT comment, post_id, user_id FROM comments WHERE post_id = $1`;
	rows, _ := app.dataBase.Query(query, pId);
	defer rows.Close()
	response := []cResponse{}
    for rows.Next() {
        var c Comments
		var r cResponse
        err = rows.Scan(&c.Comment, &c.PId, &c.Username);
		l, _ := strconv.Atoi(c.Username);
		u, _ := app.getUserById(l);
		r.Username = u.Username;
		r.Comment = c.Comment;
		r.PId = c.PId;
        if err != nil {
            fmt.Println(err)
            continue
        }
        response = append(response, r)
    }
    json.NewEncoder(writer).Encode(response)
}

func (app *App) getLikes(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Content-Type", "application/json")
    parts := strings.Split(request.URL.Path, "/getLikes/")
    pId, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        http.Error(writer, "Invalid photo ID", http.StatusBadRequest)
        return
    }

    rows, err := app.dataBase.Query(`SELECT user_id FROM likes WHERE post_id = $1`, pId)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type LikeEntry struct {
        UId int `json:"uId"`
    }
    var likes []LikeEntry
    for rows.Next() {
        var l LikeEntry
        if err := rows.Scan(&l.UId); err == nil {
            likes = append(likes, l)
        }
    }
    if likes == nil {
        likes = []LikeEntry{}
    }
    json.NewEncoder(writer).Encode(likes)
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
    json.NewEncoder(writer).Encode(app.avatars)
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                 PICTURES GETTER                                ||
// ! ||--------------------------------------------------------------------------------||

type page struct {
	Last		bool       `json:"last"`
    Pictures	[]Pictures `json:"pictures"`
	Usr			User
}


func (app *App) getPage(writer http.ResponseWriter, request *http.Request) {
    var r	page;
	if request.Method != http.MethodGet {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
	path := strings.TrimPrefix(request.URL.Path, "/getPictures/")
    page, err := strconv.Atoi(path)
    if err != nil {
		fmt.Println("Erreur 1")
        http.Error(writer, "Invalid ID", http.StatusBadRequest)
        return
    }
	pictures, err, count := app.fetchAllPicturesFromDB((page - 1) * 6)
    if err != nil {
        fmt.Println(Red + "Error fetching pictures: " + err.Error() + Reset)
        http.Error(writer, "Error fetching pictures", http.StatusInternalServerError)
        return
    }
	id := ((page - 1) * 6) + 6;
	if (id >= count) {
		r.Last = true;
	}
	r.Pictures = pictures;
    json.NewEncoder(writer).Encode(r)
}

type getPictureResponse struct {
	Picture	Pictures
	Usr		User
}

func (app *App) getPicture(writer http.ResponseWriter, request *http.Request) {
	var picture Pictures;
	fmt.Println("HERE")
    if request.Method != http.MethodGet {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    path := strings.TrimPrefix(request.URL.Path, "/getPicture/")
    id, err := strconv.Atoi(path)
    if err != nil {
		fmt.Println("Erreur 2")
        http.Error(writer, "Invalid ID", http.StatusBadRequest)
        return
    }
	query := "SELECT image_path, id, userId, uploadTime, like_count, comment_count FROM images WHERE id = $1"
	row := app.dataBase.QueryRow(query, id)
	err = row.Scan(&picture.Path, &picture.Id, &picture.UserId, &picture.UploadTime, &picture.Likes, &picture.Comments)
	if err != nil {
		fmt.Println(err)
	}
	user, _ := app.getUserById(picture.UserId);
	writer.Header().Set("Content-Type", "application/json")
	var response getPictureResponse;
	response.Picture = picture;
	response.Usr = *user;
	json.NewEncoder(writer).Encode(response);
}

func (app *App) getPictureById(id int) (*Pictures, error) {
	var picture Pictures;

	query := "SELECT image_path, id, userId, uploadTime, like_count, comment_count FROM images WHERE id = $1"
	err := app.dataBase.QueryRow(query, id).Scan(
		&picture.Path,
		&picture.Id,
		&picture.UserId,
		&picture.UploadTime,
		&picture.Likes,
		&picture.Comments,
			
	)
	if err != nil {
        return nil, err
    }

    return &picture, nil
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                            HELPER FUNCTIONS - DB QUERIES                       ||
// ! ||--------------------------------------------------------------------------------||
//

func (app *App) fetchAllPicturesFromDB(i int) ([]Pictures, error, int) {
    var pictures []Pictures
	count := 0;
	query := "SELECT COUNT(*) FROM images";
	app.dataBase.QueryRow(query).Scan(&count);
    query = "SELECT id, image_path, userId, uploadTime, like_count, comment_count FROM images ORDER BY uploadTime DESC LIMIT 6 OFFSET $1";
	rows, err := app.dataBase.Query(query, i)
    if err != nil {
        return nil, err, 0
    }
    defer rows.Close()
    for rows.Next() {
        var pic Pictures
        err := rows.Scan(&pic.Id, &pic.Path, &pic.UserId, &pic.UploadTime, &pic.Likes, &pic.Comments)
		pic.Path = "/" + pic.Path;
        if err != nil {
            return nil, err, 0
        }
        pictures = append(pictures, pic)
    }
    return pictures, nil, count
}
