package alerts

import (
	"12ul/banshee/internal/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AlertLabel struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Name      string `json:"name"`
	Service   string `json:"service"`
	Severity  string `json:"severity"`
}

type Alert struct {
	ID          string            `json:"id,omitempty"`
	Labels      AlertLabel        `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    string            `json:"activeAt"`
	Value       string            `json:"value"`
}

type DismissedAlert struct {
	Job       string
	Alertname string
	ActiveAt  string
}

type AlertRequest struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []Alert `json:"alerts"`
	} `json:"data"`
}

var (
	alerts          []Alert
	dismissedAlerts []string
)

func getIdForAlert() string {
	id := uuid.New().String()
	// Check id is already in use
	for _, alert := range alerts {
		if alert.ID == id {
			return getIdForAlert()
		}
	}
	return id
}

func fetch() {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/api/v1/alerts", config.GetEnv("PROMETHEUS_HOST"), config.GetEnv("PROMETHEUS_PORT")))

	if err != nil {
		log.Panic(`Could not fetch alerts from Prometheus:`, err)
	}

	if resp.StatusCode != 200 {
		log.Panic(fmt.Sprintf(`Could not fetch alerts from Prometheus (code: %d):`, resp.StatusCode), resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Panic(`Could not read response body:`, err)
	}

	var alertRequest AlertRequest
	err = json.Unmarshal(body, &alertRequest)

	if err != nil {
		log.Panic(`Could not parse alert reponse body:`, err)
	}

	if alertRequest.Status != "success" {
		log.Panic(`Could not fetch alerts from Prometheus, returned status:`, alertRequest.Status)
	}

	alerts = alertRequest.Data.Alerts

	// Assign ids to the alerts
	for i := range alerts {
		alerts[i].ID = getIdForAlert()
	}
	log.Printf("Fetched %d alerts from Prometheus", len(alerts))
}

func isDismissed(alert Alert) bool {
	for _, dismissedAlert := range dismissedAlerts {
		if alert.ID == dismissedAlert {
			return true
		}
	}
	return false
}

func StartFetchInterval() {
	fetch()
	// Start a interval based on a interval
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			fetch()
		}
	}()
}

func GetUndismissedAlerts() []Alert {
	undismissedAlerts := []Alert{}
	for _, alert := range alerts {
		if !isDismissed(alert) {
			undismissedAlerts = append(undismissedAlerts, alert)
		}
	}
	return undismissedAlerts
}

func DismissAlert(id string) {
	dismissedAlerts = append(dismissedAlerts, id)
}
