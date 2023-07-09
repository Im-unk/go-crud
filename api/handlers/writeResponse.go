package api

import (
	"encoding/json"
	"net/http"
)

// writeResponse writes the response in JSON format
func writeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
