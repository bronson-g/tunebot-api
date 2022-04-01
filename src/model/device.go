package model

import (
	"errors"

	"github.com/google/uuid"
)

type Device struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Mac    string `json:"mac"`
}

func (d *Device) Link() error {
	deviceId := uuid.New().String()
	_, err := db.Exec(`
		insert into device 
		(id, user_id, mac) 
		values 
		(uuid_to_bin(?), uuid_to_bin(?), ?);`,
		deviceId, d.UserId, d.Mac)

	d.Id = deviceId

	return err
}

func (d *Device) GetUser() (User, error) {
	user := User{}

	if d.UserId == "" {
		result, err := db.Query("select bin_to_uuid(user_id) as user_id from device where mac = ?;", d.Mac)

		if err != nil {
			return user, err
		}

		if !result.Next() {
			return user, errors.New("There are no users associated with this device.")
		}
		result.Scan(&d.UserId)
	}

	user.Id = d.UserId
	err := user.getPlaylists()

	return user, err
}
