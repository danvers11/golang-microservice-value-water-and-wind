package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Status struct untuk menyimpan status air dan angin
type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

// StatusText mengembalikan teks status berdasarkan nilai air dan angin
func (s *Status) StatusText() string {
	waterStatus := "aman"
	if s.Water < 5 {
		waterStatus = "aman"
	} else if s.Water >= 6 && s.Water <= 8 {
		waterStatus = "siaga"
	} else {
		waterStatus = "bahaya"
	}

	windStatus := "aman"
	if s.Wind < 6 {
		windStatus = "aman"
	} else if s.Wind >= 7 && s.Wind <= 15 {
		windStatus = "siaga"
	} else {
		windStatus = "bahaya"
	}

	return fmt.Sprintf("Status air: %s, Status angin: %s", waterStatus, windStatus)
}

func main() {
	// Mulai goroutine untuk memperbarui status setiap 15 detik
	go func() {
		for {
			updateStatus()
			time.Sleep(15 * time.Second)
		}
	}()

	// Mengatur server untuk menangani permintaan HTTP
	http.HandleFunc("/", getStatus)
	http.ListenAndServe(":8080", nil)
}

// updateStatus menghasilkan nilai acak untuk air dan angin, kemudian menyimpannya ke file JSON
func updateStatus() {
	status := Status{
		Water: rand.Intn(100) + 1,
		Wind:  rand.Intn(100) + 1,
	}

	jsonData, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = ioutil.WriteFile("status.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

// getStatus menangani permintaan HTTP untuk menampilkan status
func getStatus(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadFile("status.json")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var status Status
	err = json.Unmarshal(jsonData, &status)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	// Menampilkan status sebagai HTML dengan auto reload setiap 5 detik
	fmt.Fprintf(w, "<html><head><meta http-equiv=\"refresh\" content=\"5\"></head><body>")
	fmt.Fprintf(w, "<h1>Status Terkini:</h1>")
	fmt.Fprintf(w, "<p>%s</p>", status.StatusText())
	fmt.Fprintf(w, "<p>Nilai air: %d meter</p>", status.Water)
	fmt.Fprintf(w, "<p>Nilai angin: %d meter per detik</p>", status.Wind)
	fmt.Fprintf(w, "</body></html>")
}
