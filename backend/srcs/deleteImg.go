package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"slices"
)

type dRequest struct {
	Username	string	`json:"Username"`
	UId			int		`json:"uId"`
	PId			int		`json:"pId"`
}

func (app *App) deleteCommentFromApp(pId int) {
	for i, c := range app.comments {
		if c.PId == pId {
			app.comments = slices.Delete(app.comments, i, i+1)
            break
		}
	}
}

func (app *App) deleteCommentFromDB(pId int) (int64, error) {
	sql := `DELETE FROM comments WHERE post_id = $1`
	result, err := app.dataBase.Exec(sql, pId);
	if err != nil {
		return 0, err 
	}
	app.deleteCommentFromApp(pId);
	return result.RowsAffected();
}

func (app *App) deleteImgFromDB(pId int) (int64, error) {
    sql := `DELETE FROM images WHERE id = $1`
    result, err := app.dataBase.Exec(sql, pId)
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

func (app *App) isImageOwner(pId int, uId int) (bool, error) {
	var user *User;

	user, err := app.getUserByPhotoId(pId);
	if err != nil {
		return false, err
	}
	if user.Id != uId {
		return false, nil
	}
	return true, nil	
}

func (app *App) deleteImg(writer http.ResponseWriter, request *http.Request) {
	var r dRequest;

	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&r)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(r)
	defer request.Body.Close()
	isOwner, error := app.isImageOwner(r.PId, r.UId);
	if error != nil {
		fmt.Println(error);
	}
	if !isOwner {
        writer.WriteHeader(http.StatusForbidden)
        json.NewEncoder(writer).Encode(map[string]string{
            "error": "You are not allowed to delete this image",
        })
        return
	}
	res, _ := app.deleteCommentFromDB(r.PId);
	fmt.Println(res);
	app.deleteImgFromDB(r.PId);
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{
		"message": "Image deleted successfully",
	})
}
