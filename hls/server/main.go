package main

import (
	"io"
	"log"
	"net/http"
	"os"
)


func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		err := r.ParseMultipartForm(10 << 20)
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

		log.Println("File saved" + header.Filename)
		w.Write([]byte("File uploaded successfully"))
		return
	}
	http.NotFound(w, r)
}


func main(){
	http.HandleFunc("/", GetFile)
	port:=":8080"
	log.Printf("Starting server on port %s",port)
	http.ListenAndServe(port,nil)
}