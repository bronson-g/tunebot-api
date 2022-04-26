package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/bronson-g/tunebot-api/log"
	"github.com/bronson-g/tunebot-api/model"
)

func Register(w http.ResponseWriter, req *http.Request) {
	log.Println(log.Cyan("user.register"))

	user := model.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		errorResponse(err, w)
		return
	}

	err = user.Create()

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

func Login(w http.ResponseWriter, req *http.Request) {
	log.Println(log.Cyan("user.login"))

	user := model.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		errorResponse(err, w)
		return
	}

	err = user.Read()

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
