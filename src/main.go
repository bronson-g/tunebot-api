package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/valyala/fastjson"
)

var DB *sql.DB

type Song struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Playlist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
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
	blacklist := Playlist{"", "Blacklist", false, []Song{}}
	_, err := DB.Exec(`
		insert into user 
		(id, username, password) 
		values 
		(uuid_to_bin(?), ?, ?);`,
		id, username, password)

	if err == nil {
		blacklist, err = createPlaylist(uuid.New().String(), "Blacklist", true)
	}

	return User{id, username, password, []Playlist{}, blacklist}, err
}

func getUser(username string, password string) (User, error) {
	user := User{"", username, password, []Playlist{}, Playlist{"", "Blacklist", false, []Song{}}}
	result, err := DB.Query(`
		select bin_to_uuid(id) as id 
		from user 
		where username = ? and password = cast(? as binary(60));`,
		username, password)

	if err == nil && result.Next() {
		result.Scan(&user.Id)
		err = getPlaylists(&user)
	} else if err == nil {
		err = errors.New("invalid login credentials")
	}

	return user, err
}

func createPlaylist(userId string, name string, enabled bool) (Playlist, error) {
	id := uuid.New().String()
	_, err := DB.Exec(`
		insert into playlist 
		(id, name, user_id, enabled) 
		values 
		(uuid_to_bin(?), ?, uuid_to_bin(?), ?);`,
		id, name, userId, enabled)

	return Playlist{id, name, enabled, []Song{}}, err
}

func updatePlaylist(playlist *Playlist) error {
	_, err := DB.Exec(`
		update playlist 
		set name = ?, enabled = ? 
		where id = uuid_to_bin(?);`,
		playlist.Name, playlist.Enabled, playlist.Id)

	return err
}

func deletePlaylist(playlistId string) error {
	_, err := DB.Exec(`
		delete from playlist 
		where id = uuid_to_bin(?);`,
		playlistId)

	return err
}

func addSongToPlaylist(playlistId string, url string) (Song, error) {
	songId := uuid.New().String()
	_, err := DB.Exec(`
		insert into song 
		(id, url) 
		values 
		(uuid_to_bin(?), ?);`,
		songId, url)

	if err != nil {
		result, err := DB.Query(`
			select bin_to_uuid(id) as id 
			from song 
			where url = ?;`,
			url)

		if err == nil && result.Next() {
			result.Scan(&songId)
		}
	}

	_, err = DB.Exec(`
		insert into playlist_song 
		(id, playlist_id, song_id) 
		values 
		(uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`,
		uuid.New().String(), playlistId, songId)

	return Song{songId, url}, err
}

func removeSongFromPlaylist(playlistId string, songId string) error {
	_, err := DB.Exec(`
		delete from playlist_song 
		where playlist_id = uuid_to_bin(?)
		and song_id = uuid_to_bin(?);`,
		playlistId, songId)

	return err
}

func addPlaylist(user *User, playlist Playlist) {
	if playlist.Id == "" {
		return
	}

	if playlist.Name == "Blacklist" {
		for i := 0; i < len(playlist.Songs); i++ {
			addSong(&user.Blacklist, playlist.Songs[i])
		}

		user.Blacklist.Id = playlist.Id
		user.Blacklist.Name = playlist.Name
		user.Blacklist.Enabled = playlist.Enabled
	} else {
		for i := 0; i < len(user.Playlists); i++ {
			if playlist.Id == user.Playlists[i].Id {
				for j := 0; j < len(playlist.Songs); j++ {
					addSong(&user.Playlists[i], playlist.Songs[j])
				}
				return
			}
		}

		user.Playlists = append(user.Playlists, playlist)
	}
}

func addSong(playlist *Playlist, song Song) {
	if song.Id == "" {
		return
	}

	for i := 0; i < len(playlist.Songs); i++ {
		if playlist.Id == playlist.Songs[i].Id {
			return
		}
	}

	playlist.Songs = append(playlist.Songs, song)
}

func getPlaylists(user *User) error {
	result, err := DB.Query(`
		select bin_to_uuid(p.id) as playlist_id, p.name, bin_to_uuid(s.id) as song_id, s.url, cast(p.enabled as signed) as enabled
		from playlist as p 
		left join playlist_song 
			on p.id = playlist_song.playlist_id 
		left join song as s 
			on s.id = playlist_song.song_id 
		where p.user_id = uuid_to_bin(?)
		order by playlist_id;`,
		user.Id)

	if err == nil {
		for result.Next() {
			playlist := Playlist{"", "", false, []Song{}}
			song := Song{"", ""}

			result.Scan(&playlist.Id, &playlist.Name, &song.Id, &song.Url, &playlist.Enabled)
			addSong(&playlist, song)
			addPlaylist(user, playlist)
		}
	}

	return err
}

func registerEndpoint(w http.ResponseWriter, req *http.Request) {
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
			successResponse(data, w)
		}
	}
}

func loginEndpoint(w http.ResponseWriter, req *http.Request) {
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
			successResponse(data, w)
		}
	}
}

func createPlaylistEndpoint(w http.ResponseWriter, req *http.Request) {
	var parser fastjson.Parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	raw, err := parser.Parse(buf.String())

	if err != nil {
		errorResponse(err, w)
		return
	}

	playlist, err := createPlaylist(string(raw.GetStringBytes("userId")), string(raw.GetStringBytes("name")), true)
	if err != nil {
		errorResponse(err, w)
		return
	}
	data, err := json.Marshal(playlist)
	successResponse(data, w)
}

func updatePlaylistEndpoint(w http.ResponseWriter, req *http.Request) {
	playlist := Playlist{}
	err := json.NewDecoder(req.Body).Decode(&playlist)

	if err != nil {
		errorResponse(err, w)
		return
	}

	err = updatePlaylist(&playlist)

	if err != nil {
		errorResponse(err, w)
		return
	}

	data, err := json.Marshal(playlist)
	successResponse(data, w)
}

func deletePlaylistEndpoint(w http.ResponseWriter, req *http.Request) {
	var parser fastjson.Parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	raw, err := parser.Parse(buf.String())

	if err != nil {
		errorResponse(err, w)
		return
	}

	err = deletePlaylist(string(raw.GetStringBytes("id")))
	if err != nil {
		errorResponse(err, w)
		return
	}
	successResponse([]byte("{}"), w)
}

func addSongEndpoint(w http.ResponseWriter, req *http.Request) {
	var parser fastjson.Parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	raw, err := parser.Parse(buf.String())

	if err != nil {
		errorResponse(err, w)
		return
	}

	song, err := addSongToPlaylist(string(raw.GetStringBytes("playlistId")), string(raw.GetStringBytes("url")))
	if err != nil {
		errorResponse(err, w)
		return
	}
	data, err := json.Marshal(song)

	successResponse(data, w)
}

func removeSongEndpoint(w http.ResponseWriter, req *http.Request) {
	var parser fastjson.Parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	raw, err := parser.Parse(buf.String())

	if err != nil {
		errorResponse(err, w)
		return
	}

	err = removeSongFromPlaylist(string(raw.GetStringBytes("playlistId")), string(raw.GetStringBytes("songId")))
	if err != nil {
		errorResponse(err, w)
		return
	}

	successResponse([]byte("{}"), w)
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

	router.HandleFunc("/user/register/", registerEndpoint).Methods("POST")
	router.HandleFunc("/user/login/", loginEndpoint).Methods("POST")
	router.HandleFunc("/playlist/create/", createPlaylistEndpoint).Methods("POST")
	router.HandleFunc("/playlist/update/", updatePlaylistEndpoint).Methods("POST")
	router.HandleFunc("/playlist/delete/", deletePlaylistEndpoint).Methods("POST")
	router.HandleFunc("/playlist/song/add/", addSongEndpoint).Methods("POST")
	router.HandleFunc("/playlist/song/remove/", removeSongEndpoint).Methods("POST")
	http.ListenAndServe(":8080", router)
}

func successResponse(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func errorResponse(err error, w http.ResponseWriter) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + strings.Replace(err.Error(), "\"", "\\\"", -1) + "\"}"))
	}
}
