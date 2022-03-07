package model

import "github.com/google/uuid"

type Playlist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Songs   []Song `json:"songs"`
}

func (p *Playlist) Create(userId string) error {
	id := uuid.New().String()
	_, err := db.Exec(`
		insert into playlist 
		(id, name, user_id, enabled) 
		values 
		(uuid_to_bin(?), ?, uuid_to_bin(?), ?);`,
		id, p.Name, userId, p.Enabled)

	p.Id = id
	p.Songs = []Song{}

	return err
}

func (p *Playlist) Update() error {
	_, err := db.Exec(`
		update playlist 
		set name = ?, enabled = ? 
		where id = uuid_to_bin(?);`,
		p.Name, p.Enabled, p.Id)

	if p.Songs == nil {
		p.Songs = []Song{}
	}

	return err
}

func (p *Playlist) Delete() error {
	_, err := db.Exec(`
		delete from playlist 
		where id = uuid_to_bin(?);`,
		p.Id)

	return err
}
