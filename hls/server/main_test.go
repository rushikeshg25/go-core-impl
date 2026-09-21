package main

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestUploadProcessingAndRestart(t *testing.T) {
	s, e := newService(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	s.transcode = func(ctx context.Context, input, dir string) error {
		return os.WriteFile(filepath.Join(dir, "index.m3u8"), []byte("#EXTM3U"), 0600)
	}
	var body bytes.Buffer
	m := multipart.NewWriter(&body)
	f, _ := m.CreateFormFile("file", "video.mp4")
	f.Write([]byte("sample"))
	m.Close()
	r := httptest.NewRequest("POST", "/api/videos", &body)
	r.Header.Set("Content-Type", m.FormDataContentType())
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 202 {
		t.Fatal(w.Code, w.Body.String())
	}
	var v video
	json.Unmarshal(w.Body.Bytes(), &v)
	s.wg.Wait()
	w = httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/hls/"+v.ID+"/index.m3u8", nil))
	if w.Code != 200 || w.Body.String() != "#EXTM3U" {
		t.Fatal(w.Code, w.Body.String())
	}
	restored, e := newService(s.root)
	if e != nil {
		t.Fatal(e)
	}
	defer restored.Close()
	if restored.videos[v.ID].Status != "completed" {
		t.Fatal("lost status")
	}
}
func TestFFmpegPipeline(t *testing.T) {
	if _, e := exec.LookPath("ffmpeg"); e != nil {
		t.Skip("ffmpeg unavailable")
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "input.mp4")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	b, e := exec.CommandContext(ctx, "ffmpeg", "-nostdin", "-y", "-f", "lavfi", "-i", "color=c=blue:s=64x64:d=0.3", "-c:v", "libx264", input).CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	if e = transcode(ctx, input, dir); e != nil {
		t.Fatal(e)
	}
	playlist, e := os.ReadFile(filepath.Join(dir, "index.m3u8"))
	if e != nil || !bytes.Contains(playlist, []byte("#EXT-X-ENDLIST")) {
		t.Fatal(string(playlist), e)
	}
	segments, _ := filepath.Glob(filepath.Join(dir, "*.ts"))
	if len(segments) == 0 {
		t.Fatal("no segments")
	}
}
func TestRejectInvalidUpload(t *testing.T) {
	s, _ := newService(t.TempDir())
	defer s.Close()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("POST", "/api/videos", bytes.NewReader(nil)))
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}
