package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var monitored []string
var port string

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func callHealthcheck(url string, ch chan<- int16) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR: Healthcheck at %s returned error.\n", url)
		ch <- 0
		return
	}
	ch <- int16(resp.StatusCode)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/healthcheck" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}
	ch := make(chan int16)
	for _, url := range monitored {
		go callHealthcheck(url, ch)
	}
	for range monitored {
		status := <-ch
		if status == 0 {
			http.Error(w, "{\"status\": \"One or more addresses is not responding.\"", http.StatusInternalServerError)
			return
		}
		if status != 200 {
			http.Error(w, "{\"status\": \"One or more endpoint responded with non-200 HTTP status.\"", http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, "{\"status\": \"OK\"}")
}

func init() {
	fmt.Println("INFO: Initializing Simple Healthchecker")
	tmp_monitored := getEnv("MONITORED_URLS", "")
	if tmp_monitored == "" {
		fmt.Println("ERROR: MONITORED_URLS environment variable is empty.")
		os.Exit(1)
	}
	tmp_monitored = strings.ReplaceAll(tmp_monitored, " ", "")
	monitored = strings.Split(tmp_monitored, ";")
	for _, url := range monitored {
		fmt.Printf("INFO: Monitoring %s\n", url)
	}
	port = getEnv("SERVER_PORT", "8080")
}

func main() {
	http.HandleFunc("/healthcheck", healthcheckHandler) // Update this line of code

	fmt.Printf("INFO: Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
