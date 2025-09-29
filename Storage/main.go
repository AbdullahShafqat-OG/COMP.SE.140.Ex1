package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const logFilePath = "/app/data/log.txt"

func ensureLogFile() error {
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func postLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Error opening log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := file.Write(body); err != nil {
		http.Error(w, "Error writing to log file", http.StatusInternalServerError)
		return
	}
	
	if _, err := file.WriteString("\n"); err != nil {
		http.Error(w, "Error writing newline to log file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func getLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Error reading log file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func main() {
	if err := ensureLogFile(); err != nil {
		log.Fatal("Failed to ensure log file:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/log", postLogHandler).Methods("POST")
	r.HandleFunc("/log", getLogHandler).Methods("GET")

	log.Println("Storage service starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}