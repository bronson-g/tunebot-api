package endpoint

import (
	"net/http"
	"strings"

	"github.com/bronson-g/tunebot-api/log"
)

func successResponse(data []byte, w http.ResponseWriter) {
	log.Println(log.Green(string(data)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func errorResponse(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(log.Red(err.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + strings.Replace(err.Error(), "\"", "\\\"", -1) + "\"}"))
	}
}
