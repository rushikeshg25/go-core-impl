package main

import (
	"log"
	"net/http"
)

func main() {
  http.Handle("/health", enableCORS(http.HandlerFunc(Health)))
  http.Handle("/ws", enableCORS(http.HandlerFunc(WsHandler)))

  port := ":8080"
  log.Printf("Starting server on port %s", port)
  if err := http.ListenAndServe(port, nil); err != nil {
          log.Fatalf("Server failed to start: %v", err)
  }

}


func WsHandler(w http.ResponseWriter,r *http.Request){

}

func Health(w http.ResponseWriter,r *http.Request){
   w.WriteHeader(http.StatusOK)
   w.Write([]byte("Healthy"))
   
}

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