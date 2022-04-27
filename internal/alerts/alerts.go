package alerts

import (
	"12ul/banshee/internal/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	alerts []Alert = []Alert{
		// {
		// 	ID: "DIT_IS_EEN_TEST_ID",
		// 	Labels: AlertLabel{
		// 		Alertname: "service_up",
		// 		Instance:  "172.12.50.22",
		// 		Job:       "blackbox-ping",
		// 		Name:      "client2",
		// 		Service:   "client2",
		// 		Severity:  "gravest",
		// 	},
		// 	Annotations: make(map[string]string),
		// 	State:       "firing",
		// 	ActiveAt:    "2022-04-26T19:40:31.656129563Z",
		// 	Value:       "1e+00",
		// },
	}
	dismissedAlerts []string = []string{}
)

func generateIdForAlert(alert Alert) (string, bool) {
	id := fmt.Sprintf("%s-%s-%s-%s-%s", alert.Labels.Service, alert.Labels.Job, alert.Labels.Name, alert.Labels.Severity, alert.ActiveAt)
	for _, alert := range alerts {
		if alert.ID == id {
			return "", true
		}
	}
	return id, false
}

func fetch() {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/api/v1/alerts", config.GetEnv("PROMETHEUS_HOST"), config.GetEnv("PROMETHEUS_PORT")))

	if err != nil {
		log.Println(`Could not fetch alerts from Prometheus:`, err)
	}

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf(`Could not fetch alerts from Prometheus (code: %d):`, resp.StatusCode), resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(`Could not read response body:`, err)
	}

	var alertRequest AlertRequest
	err = json.Unmarshal(body, &alertRequest)

	if err != nil {
		log.Println(`Could not parse alert reponse body:`, err)
	}

	if alertRequest.Status != "success" {
		log.Println(`Could not fetch alerts from Prometheus, returned status:`, alertRequest.Status)
	}

	newAlerts := alertRequest.Data.Alerts

	idxToDelete := []int{}
	// Assign ids to the alerts
	for i, alert := range newAlerts {
		id, alreadyExists := generateIdForAlert(alert)
		if alreadyExists || alert.State != "firing" {
			idxToDelete = append(idxToDelete, i)
		} else {
			newAlerts[i].ID = id
		}
	}
	// loop over idxToDelete backwards to not mess up the index
	for i := len(idxToDelete) - 1; i >= 0; i-- {
		newAlerts = append(alerts[:idxToDelete[i]], alerts[idxToDelete[i]+1:]...)
	}
	// Add newAlerts to alerts
	alerts = append(alerts, newAlerts...)
	log.Printf("Fetched %d new alerts from Prometheus", len(newAlerts))
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

func GetDismissedAlerts() []string {
	return dismissedAlerts
}

func DismissAlert(id string) {
	log.Printf("Dismissing alert %s", id)
	dismissedAlerts = append(dismissedAlerts, id)
}
