package http

import (
	"encoding/json"
	"mockium/pkg/model"
	"net/http"
)

type defaultResponser struct{}

func (defaultResponser) Write(w http.ResponseWriter, r *http.Request, resp *model.SetResponse) {
	for k, v := range resp.SetHeaders {
		w.Header().Set(k, v)
	}

	status := resp.SetStatus
	if status == 0 {
		status = http.StatusOK
	}

	switch {
	case resp.SetFile != "":
		w.WriteHeader(status)
		http.ServeFile(w, r, resp.SetFile)
	case resp.SetBody != nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(resp.SetBody)
	default:
		w.WriteHeader(status)
	}
}

func (defaultResponser) Error(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
