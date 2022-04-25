package endpoint

import (
	"fmt"
	"net/http"
	"strings"
)

func successResponse(data []byte, w http.ResponseWriter) {
	fmt.Println("\\e[32m" + string(data) + "\\e[39m")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func errorResponse(err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("\\e[31m" + err.Error() + "\\e[39m")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + strings.Replace(err.Error(), "\"", "\\\"", -1) + "\"}"))
	}
}
