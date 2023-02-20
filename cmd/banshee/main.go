package main

import (
	"12ul/banshee/internal/alerts"
	"12ul/banshee/internal/config"
	"12ul/banshee/internal/routes"
	"log"
	"net/http"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Static files
	fileServer := http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
	http.Handle("/public/", fileServer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	http.HandleFunc("/api/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			routes.HandleAlertFetch(w, r)
			return
		}
		if r.Method == http.MethodPut {
			routes.HandleAlertDismiss(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	alerts.StartFetchInterval()
	if err := http.ListenAndServe(":"+config.GetEnv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
