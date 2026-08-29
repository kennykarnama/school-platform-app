package shared

import (
	"encoding/json"
	"net/http"
)

type Empty struct{}

func ResponseJson(w http.ResponseWriter, payload interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	rep, _ := json.Marshal(payload)
	w.Write(rep)
}
