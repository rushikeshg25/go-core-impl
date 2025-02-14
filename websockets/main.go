package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func main() {
  http.Handle("/health", enableCORS(http.HandlerFunc(Health)))
  http.Handle("/ws", enableCORS(http.HandlerFunc(WsHandler)))

  port := ":8080"
  log.Printf("Starting server on port %s", port)
  if err := http.ListenAndServe(port, nil); err != nil {
          log.Fatalf("Server failed to start: %v", err)
  }

}


func WsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !strings.Contains(r.Header.Get("Connection"), "Upgrade") ||
		strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		http.Error(w, "Invalid WebSocket handshake", http.StatusBadRequest)
		return
	}

	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	if secWebSocketKey == "" {
		http.Error(w, "Missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}

	acceptKey := computeWebSocketAcceptKey(secWebSocketKey)

	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Accept: %s\r\n\r\n", acceptKey)

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Println("Failed to send WebSocket handshake response:", err)
		return
	}

	log.Println("WebSocket handshake successful!")
}

func Health(w http.ResponseWriter,r *http.Request){
   w.WriteHeader(http.StatusOK)
   w.Write([]byte("Healthy"))
   
}


func computeWebSocketAcceptKey(secWebSocketKey string) string {
	hash := sha1.Sum([]byte(secWebSocketKey + magicGUID))
	return base64.StdEncoding.EncodeToString(hash[:])
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