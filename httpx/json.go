package httpx

import (
	"encoding/json"
	"net/http"
)

const (
	JSONContentType = "application/json"
)

func WriteJSONResponse(w http.ResponseWriter, code int, v interface{}) error {
	raw, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(code)

	if _, err = w.Write(raw); err != nil {
		return err
	}

	return nil
}
