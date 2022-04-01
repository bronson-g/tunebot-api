package model

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	Id        string     `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password,omitempty"` // remember to ALWAYS set this to = "" before writing out data!
	Playlists []Playlist `json:"playlists"`
	Blacklist Playlist   `json:"blacklist"`
}

func (u *User) Create() error {
	id := uuid.New().String()
	blacklist := Playlist{"", "Blacklist", false, []Song{}}
	_, err := db.Exec(`
		insert into user 
		(id, username, password) 
		values 
		(uuid_to_bin(?), ?, ?);`,
		id, u.Username, u.Password)

	if err == nil {
		blacklist = Playlist{"", "Blacklist", true, []Song{}}
		blacklist.Create(id)
	}

	u.Id = id
	u.Blacklist = blacklist
	u.Playlists = []Playlist{}

	return err
}

func (u *User) Read() error {
	u.Playlists = []Playlist{}
	u.Blacklist = Playlist{"", "Blacklist", false, []Song{}}

	result, err := db.Query(`
		select bin_to_uuid(id) as id 
		from user 
		where username = ? and password = cast(? as binary(60));`,
		u.Username, u.Password)

	if err == nil && result.Next() {
		result.Scan(&u.Id)
		err = u.getPlaylists()
	} else if err == nil {
		err = errors.New("invalid login credentials")
	}

	return err
}

func (u *User) addPlaylist(playlist Playlist) {
	if playlist.Id == "" {
		return
	}

	if playlist.Name == "Blacklist" {
		for i := 0; i < len(playlist.Songs); i++ {
			addSong(&u.Blacklist, playlist.Songs[i])
		}

		u.Blacklist.Id = playlist.Id
		u.Blacklist.Name = playlist.Name
		u.Blacklist.Enabled = playlist.Enabled
	} else {
		for i := 0; i < len(u.Playlists); i++ {
			if playlist.Id == u.Playlists[i].Id {
				for j := 0; j < len(playlist.Songs); j++ {
					addSong(&u.Playlists[i], playlist.Songs[j])
				}
				return
			}
		}

		u.Playlists = append(u.Playlists, playlist)
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

func (u *User) getPlaylists() error {
	u.Playlists = []Playlist{}
	u.Blacklist = Playlist{"", "Blacklist", false, []Song{}}

	result, err := db.Query(`
		select bin_to_uuid(p.id) as playlist_id, p.name, bin_to_uuid(s.id) as song_id, s.url, cast(p.enabled as signed) as enabled
		from playlist as p 
		left join playlist_song 
			on p.id = playlist_song.playlist_id 
		left join song as s 
			on s.id = playlist_song.song_id 
		where p.user_id = uuid_to_bin(?)
		order by playlist_id;`,
		u.Id)

	if err == nil {
		for result.Next() {
			playlist := Playlist{"", "", false, []Song{}}
			song := Song{"", ""}

			result.Scan(&playlist.Id, &playlist.Name, &song.Id, &song.Url, &playlist.Enabled)
			addSong(&playlist, song)
			u.addPlaylist(playlist)
		}
	}

	return err
}
