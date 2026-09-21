package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func masked(op byte, fin bool, b []byte) []byte {
	h := op
	if fin {
		h |= 128
	}
	out := []byte{h, 128 | byte(len(b)), 1, 2, 3, 4}
	for i, v := range b {
		out = append(out, v^byte(i%4+1))
	}
	return out
}
func TestEchoFragmentPingClose(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(WsHandler))
	defer s.Close()
	c, e := net.Dial("tcp", strings.TrimPrefix(s.URL, "http://"))
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(3 * time.Second))
	fmt.Fprintf(c, "GET / HTTP/1.1\r\nHost: %s\r\nConnection: keep-alive, Upgrade\r\nUpgrade: websocket\r\nSec-WebSocket-Version: 13\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n\r\n", strings.TrimPrefix(s.URL, "http://"))
	r := bufio.NewReader(c)
	resp, e := http.ReadResponse(r, nil)
	if e != nil || resp.StatusCode != 101 {
		t.Fatal(resp, e)
	}
	if resp.Header.Get("Sec-WebSocket-Accept") != "s3pPLMBiTxaQ9kYGzzhZRbK+xOo=" {
		t.Fatal("bad accept")
	}
	c.Write(masked(1, false, []byte("hel")))
	c.Write(masked(9, true, []byte("ping")))
	read := func(op byte, want string) {
		t.Helper()
		h := make([]byte, 2)
		io.ReadFull(r, h)
		b := make([]byte, int(h[1]))
		_, e := io.ReadFull(r, b)
		if e != nil || h[0] != 128|op || string(b) != want {
			t.Fatalf("frame %x %q %v", h, b, e)
		}
	}
	read(10, "ping")
	c.Write(masked(0, true, []byte("lo")))
	read(1, "hello")
	c.Write(masked(8, true, []byte{3, 232}))
	read(8, string([]byte{3, 232}))
}
func TestRejectProtocol(t *testing.T) {
	for _, b := range [][]byte{{0x81, 0}, {0xC1, 0x80, 0, 0, 0, 0}, {9, 0x80, 0, 0, 0, 0}, {0x81, 0xfe, 0, 1}, {0x83, 0x80, 0, 0, 0, 0}} {
		if _, e := readFrame(bytes.NewReader(b)); e == nil {
			t.Fatalf("accepted %x", b)
		}
	}
	for _, code := range []uint16{999, 1005, 1015, 2000, 5000} {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, code)
		if validClose(b) {
			t.Fatal(code)
		}
	}
}
func TestRejectUpgrade(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example/ws", nil)
	w := httptest.NewRecorder()
	WsHandler(w, r)
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}
func FuzzFrame(f *testing.F) {
	f.Add(masked(1, true, []byte("hi")))
	f.Fuzz(func(t *testing.T, b []byte) {
		if len(b) > maxMessage+14 {
			t.Skip()
		}
		readFrame(bytes.NewReader(b))
	})
}
