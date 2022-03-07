package model

import "github.com/google/uuid"

type Song struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (s *Song) AddToPlaylist(playlistId string) error {
	songId := uuid.New().String()
	_, err := db.Exec(`
		insert into song 
		(id, url) 
		values 
		(uuid_to_bin(?), ?);`,
		songId, s.Url)

	if err != nil {
		result, err := db.Query(`
			select bin_to_uuid(id) as id 
			from song 
			where url = ?;`,
			s.Url)

		if err == nil && result.Next() {
			result.Scan(&songId)
		}
	}

	_, err = db.Exec(`
		insert into playlist_song 
		(id, playlist_id, song_id) 
		values 
		(uuid_to_bin(?), uuid_to_bin(?), uuid_to_bin(?));`,
		uuid.New().String(), playlistId, songId)

	s.Id = songId

	return err
}

func (s Song) RemoveFromPlaylist(playlistId string) error {
	_, err := db.Exec(`
		delete from playlist_song 
		where playlist_id = uuid_to_bin(?)
		and song_id = uuid_to_bin(?);`,
		playlistId, s.Id)

	return err
}
