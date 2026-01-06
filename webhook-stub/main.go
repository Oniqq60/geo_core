package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := json.MarshalIndent(payload, "", "  ")
		log.Printf("Webhook received:\n%s\n", string(data))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Webhook stub listening on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
