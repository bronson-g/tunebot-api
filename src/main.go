package main

import (
	"database/sql"
	"fmt"
	"net/http"

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
	Password  string     `json:"password"`
	Playlists []Playlist `json:"playlists"`
	Blacklist Playlist   `json:"blacklist"`
}

func CreateUser(username string, password string) (User, error) {
	id := uuid.New()
	_, err := DB.Exec("insert into `user` (`id`, `username`, `password`) values (uuid_to_bin(?), ?, ?);", id, username, password)
	return User{id, username, password, []Playlist{}, Playlist{}}, err
}

func GetUser(username string, password string) (User, error) {
	user := User{"", username, password, []Playlist{}, Playlist{}}
	result, err := DB.Query("select bin_to_uuid(`id`) as `id` from `user` where `username` = ? and `password` = ?;", username, password)

	if err != nil && result.Next() {
		result.Scan(&user.Id)
	} else {
		// user not found or err
	}

	return user, getPlaylists(&user)
}

func getPlaylists(user *User) error {
	result, err := DB.Query("select bin_to_uuid(`playlist`.`id`) as `playlist_id`, `playlist`.`enabled`, `playlist`.`is_blacklist`, bin_to_uuid(`song`.`id`) as `song_id`, `song`.`url` from `playlist` inner join `playlist_song` on `playlist`.`id` = `playlist_song`.`playlist_id` inner join `song` on `song`.`id` = `playlist_song`.`song_id` where `playlist`.`user_id` = ?;", user.Id)

	if err != nil {
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

func register(w http.ResponseWriter, r *http.Request) {
	// make sure this device isn't already registered
	fmt.Fprintf(w, "register")
}

func main() {
	DB, err := sql.Open("mysql", "root:Se4Q2Lp-3587@tcp(localhost)/tunebot")
	if err != nil {
		//todo handle it (fatal, can't continue.)
	}
	defer DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/register/", register).Methods("POST")
	http.ListenAndServe(":8080", router)
}
