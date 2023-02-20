package routes

import (
	"12ul/banshee/internal/alerts"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func HandleAlertFetch(w http.ResponseWriter, r *http.Request) {
	resStruct := struct {
		Alerts          []alerts.Alert `json:"alerts"`
		DismissedAlerts []string       `json:"dismissedAlerts"`
	}{
		Alerts:          alerts.GetUndismissedAlerts(),
		DismissedAlerts: alerts.GetDismissedAlerts(),
	}

	res, err := json.Marshal(resStruct)
	if err != nil {
		log.Print(err)
		w.Write([]byte("Failed to transform alerts to valid json"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func HandleAlertDismiss(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.Write([]byte("Failed to read the request body"))
		return
	}

	alerts.DismissAlert(string(bodyBytes))
	w.WriteHeader(http.StatusOK)
}
