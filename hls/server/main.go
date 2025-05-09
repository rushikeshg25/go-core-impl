package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println("Error retrieving file:", err)
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		err = os.MkdirAll("uploads", os.ModePerm)
		if err != nil {
			log.Println("Error creating uploads directory:", err)
			http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
			return
		}

		outFile, err := os.Create("uploads/" + header.Filename)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Println("Error saving file:", err)
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		log.Println("File saved: " + header.Filename)
		w.Write([]byte("File uploaded successfully"))
		return
	}
	http.NotFound(w, r)
}

func Hls() {
	hlsDir := "./hls"
	_, err := os.Stat(hlsDir)
	if os.IsNotExist(err) {
		log.Fatalf("./hls does not exist: %v", err)
	}
	fs := http.FileServer(http.Dir(hlsDir))
	http.Handle("/hls/", http.StripPrefix("/hls/", fs))
}

func main() {
	http.Handle("/", enableCORS(http.HandlerFunc(GetFile)))

	Hls()

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}