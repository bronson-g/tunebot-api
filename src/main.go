package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var DB *sql.DB

type Song struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Playlist struct {
	Id      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Songs   []Song `json:"songs"`
}

type User struct {
	Id        string     `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password,omitempty"` // remember to ALWAYS set this to = "" before writing out data!
	Playlists []Playlist `json:"playlists"`
	Blacklist Playlist   `json:"blacklist"`
}

func createUser(username string, password string) (User, error) {
	id := uuid.New().String()
	_, err := DB.Exec("insert into `user` (`id`, `username`, `password`) values (uuid_to_bin(?), ?, ?);", id, username, password)
	return User{id, username, password, []Playlist{}, Playlist{}}, err
}

func getUser(username string, password string) (User, error) {
	user := User{"", username, password, []Playlist{}, Playlist{"", false, []Song{}}}
	result, err := DB.Query("select bin_to_uuid(`id`) as `id` from `user` where `username` = ? and `password` = ?;", username, password)

	if err == nil && result.Next() {
		result.Scan(&user.Id)
	} else {
		// user not found or err
		fmt.Println(err.Error())
	}

	return user, getPlaylists(&user)
}

func getPlaylists(user *User) error {
	//todo remember to insert is_blacklist == false as NULL
	result, err := DB.Query("select bin_to_uuid(`playlist`.`id`) as `playlist_id`, `playlist`.`enabled`, coalesce(`playlist`.`is_blacklist`, 0) as `is_blacklist`, bin_to_uuid(`song`.`id`) as `song_id`, `song`.`url` from `playlist` inner join `playlist_song` on `playlist`.`id` = `playlist_song`.`playlist_id` inner join `song` on `song`.`id` = `playlist_song`.`song_id` where `playlist`.`user_id` = ?;", user.Id)

	if err == nil {
		playlist := Playlist{"", false, []Song{}}
		tempPlaylist := Playlist{"", false, []Song{}}
		wasBlacklist := false

		for result.Next() {
			song := Song{}

			isBlacklist := false
			result.Scan(&tempPlaylist.Id, &tempPlaylist.Enabled, &isBlacklist, &song.Id, &song.Url)

			if tempPlaylist.Id != playlist.Id {
				if playlist.Id != "" {
					if wasBlacklist {
						user.Blacklist = playlist
					} else {
						user.Playlists = append(user.Playlists, playlist)
					}
				}
				tempPlaylist.Songs = append(tempPlaylist.Songs, song)
				playlist = tempPlaylist
			}
			wasBlacklist = isBlacklist
		}
	}

	return err
}

func register(w http.ResponseWriter, req *http.Request) {
	user := User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		errorResponse(err, w)
		return
	}

	user, err = createUser(user.Username, user.Password)

	if err != nil {
		errorResponse(err, w)
	} else {
		user.Password = ""
		data, err := json.Marshal(user)

		if err != nil {
			errorResponse(err, w)
		} else {
			w.Write(data)
		}
	}
}

func login(w http.ResponseWriter, req *http.Request) {
	user := User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		errorResponse(err, w)
		return
	}

	user, err = getUser(user.Username, user.Password)

	if err != nil {
		errorResponse(err, w)
	} else {
		user.Password = ""
		data, err := json.Marshal(user)

		if err != nil {
			errorResponse(err, w)
		} else {
			w.Write(data)
		}
	}
}

func main() {
	db, err := sql.Open("mysql", "root:Se4Q2Lp-3587@tcp(localhost)/tunebot")
	if db == nil || err != nil {
		fmt.Println("Failed to connect to database")
		fmt.Println(err.Error())
		return
	}
	DB = db
	defer DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/register/", register).Methods("POST")
	router.HandleFunc("/api/login/", login).Methods("POST")
	http.ListenAndServe(":8080", router)
}

func errorResponse(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte("{\"error\":\"" + strings.Replace(err.Error(), "\"", "\\\"", -1) + "\"}"))
	}
}
