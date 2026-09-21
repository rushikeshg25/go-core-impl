package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
const maxMessage = 1 << 20

func token(h, want string) bool {
	for _, v := range strings.Split(h, ",") {
		if strings.EqualFold(strings.TrimSpace(v), want) {
			return true
		}
	}
	return false
}
func computeWebSocketAcceptKey(key string) string {
	h := sha1.Sum([]byte(key + magicGUID))
	return base64.StdEncoding.EncodeToString(h[:])
}
func WsHandler(w http.ResponseWriter, r *http.Request) {
	key, e := base64.StdEncoding.DecodeString(r.Header.Get("Sec-WebSocket-Key"))
	if r.Method != "GET" || !token(r.Header.Get("Connection"), "upgrade") || !token(r.Header.Get("Upgrade"), "websocket") || e != nil || len(key) != 16 {
		http.Error(w, "invalid upgrade", 400)
		return
	}
	if r.Header.Get("Sec-WebSocket-Version") != "13" {
		w.Header().Set("Sec-WebSocket-Version", "13")
		http.Error(w, "version 13 required", 426)
		return
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || !strings.EqualFold(u.Host, r.Host) {
			http.Error(w, "origin rejected", 403)
			return
		}
	}
	h, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "upgrade unavailable", 500)
		return
	}
	conn, rw, err := h.Hijack()
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Fprintf(rw, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", computeWebSocketAcceptKey(r.Header.Get("Sec-WebSocket-Key")))
	if rw.Flush() != nil {
		return
	}
	var message []byte
	var messageOp byte
	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		frame, err := readFrame(rw.Reader)
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			if !errors.Is(err, io.EOF) {
				writeClose(rw.Writer, 1002)
			}
			return
		}
		switch frame.op {
		case 8:
			if !validClose(frame.data) {
				writeClose(rw.Writer, 1002)
			} else {
				writeFrame(rw.Writer, 8, frame.data)
			}
			return
		case 9:
			if writeFrame(rw.Writer, 10, frame.data) != nil {
				return
			}
			continue
		case 10:
			continue
		case 0:
			if messageOp == 0 {
				writeClose(rw.Writer, 1002)
				return
			}
		case 1, 2:
			if messageOp != 0 {
				writeClose(rw.Writer, 1002)
				return
			}
			messageOp = frame.op
		default:
			writeClose(rw.Writer, 1002)
			return
		}
		if len(message)+len(frame.data) > maxMessage {
			writeClose(rw.Writer, 1009)
			return
		}
		message = append(message, frame.data...)
		if frame.fin {
			if messageOp == 1 && !utf8.Valid(message) {
				writeClose(rw.Writer, 1007)
				return
			}
			if writeFrame(rw.Writer, messageOp, message) != nil {
				return
			}
			message = nil
			messageOp = 0
		}
	}
}

type frame struct {
	fin  bool
	op   byte
	data []byte
}

func readFrame(r io.Reader) (frame, error) {
	var f frame
	var h [2]byte
	if _, e := io.ReadFull(r, h[:]); e != nil {
		return f, e
	}
	f.fin = h[0]&128 != 0
	f.op = h[0] & 15
	if h[0]&112 != 0 || h[1]&128 == 0 {
		return f, errors.New("reserved bits or unmasked client")
	}
	n := uint64(h[1] & 127)
	switch n {
	case 126:
		var b [2]byte
		if _, e := io.ReadFull(r, b[:]); e != nil {
			return f, e
		}
		n = uint64(binary.BigEndian.Uint16(b[:]))
		if n < 126 {
			return f, errors.New("noncanonical length")
		}
	case 127:
		var b [8]byte
		if _, e := io.ReadFull(r, b[:]); e != nil {
			return f, e
		}
		n = binary.BigEndian.Uint64(b[:])
		if n < 65536 || n>>63 != 0 {
			return f, errors.New("invalid length")
		}
	}
	if n > maxMessage || f.op >= 8 && (!f.fin || n > 125) {
		return f, errors.New("frame limit")
	}
	if f.op != 0 && f.op != 1 && f.op != 2 && f.op != 8 && f.op != 9 && f.op != 10 {
		return f, errors.New("reserved opcode")
	}
	var mask [4]byte
	if _, e := io.ReadFull(r, mask[:]); e != nil {
		return f, e
	}
	f.data = make([]byte, int(n))
	if _, e := io.ReadFull(r, f.data); e != nil {
		return f, e
	}
	for i := range f.data {
		f.data[i] ^= mask[i%4]
	}
	return f, nil
}
func writeFrame(w *bufio.Writer, op byte, data []byte) error {
	h := []byte{128 | op}
	n := len(data)
	if n < 126 {
		h = append(h, byte(n))
	} else if n <= 65535 {
		h = append(h, 126, byte(n>>8), byte(n))
	} else {
		h = append(h, 127, 0, 0, 0, 0, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
	}
	if _, e := w.Write(h); e != nil {
		return e
	}
	if _, e := w.Write(data); e != nil {
		return e
	}
	return w.Flush()
}
func writeClose(w *bufio.Writer, code uint16) {
	_ = writeFrame(w, 8, []byte{byte(code >> 8), byte(code)})
}
func validClose(b []byte) bool {
	if len(b) == 0 {
		return true
	}
	if len(b) == 1 || !utf8.Valid(b[2:]) {
		return false
	}
	c := binary.BigEndian.Uint16(b)
	return c >= 3000 && c <= 4999 || c >= 1000 && c <= 1014 && c != 1004 && c != 1005 && c != 1006
}
func Health(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Healthy")) }
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/ws", WsHandler)
	log.Fatal((&http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}).ListenAndServe())
}
