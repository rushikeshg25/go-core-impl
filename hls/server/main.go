package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type video struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
	Error  string `json:"error,omitempty"`
}
type service struct {
	root      string
	mu        sync.Mutex
	videos    map[string]video
	slots     chan struct{}
	transcode func(context.Context, string, string) error
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func newService(root string) (*service, error) {
	if e := os.MkdirAll(root, 0700); e != nil {
		return nil, e
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &service{root: root, videos: map[string]video{}, slots: make(chan struct{}, 1), transcode: transcode, ctx: ctx, cancel: cancel}
	dirs, e := os.ReadDir(root)
	if e != nil {
		cancel()
		return nil, e
	}
	for _, dir := range dirs {
		if !validID(dir.Name()) || !dir.IsDir() {
			continue
		}
		b, e := os.ReadFile(filepath.Join(root, dir.Name(), "status.json"))
		if e != nil {
			continue
		}
		var v video
		if json.Unmarshal(b, &v) != nil || v.ID != dir.Name() {
			continue
		}
		if v.Status == "processing" {
			v.Status = "failed"
			v.Error = "processing interrupted by restart"
		}
		s.videos[v.ID] = v
	}
	return s, nil
}
func validID(id string) bool { b, e := hex.DecodeString(id); return e == nil && len(b) == 16 }
func (s *service) save(v video) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	dir := filepath.Join(s.root, v.ID)
	p := filepath.Join(dir, "status.tmp")
	if e = os.WriteFile(p, b, 0600); e != nil {
		return e
	}
	if e = os.Rename(p, filepath.Join(dir, "status.json")); e != nil {
		return e
	}
	s.mu.Lock()
	s.videos[v.ID] = v
	s.mu.Unlock()
	return nil
}
func (s *service) Close() { s.cancel(); s.wg.Wait() }
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}
	switch {
	case r.URL.Path == "/api/videos" && r.Method == "POST":
		s.upload(w, r)
	case r.URL.Path == "/api/videos" && r.Method == "GET":
		s.mu.Lock()
		out := make([]video, 0, len(s.videos))
		for _, v := range s.videos {
			out = append(out, v)
		}
		s.mu.Unlock()
		sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	case strings.HasPrefix(r.URL.Path, "/api/videos/") && r.Method == "GET":
		id := strings.TrimPrefix(r.URL.Path, "/api/videos/")
		s.mu.Lock()
		v, ok := s.videos[id]
		s.mu.Unlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	case strings.HasPrefix(r.URL.Path, "/hls/") && r.Method == "GET":
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/hls/"), "/")
		if len(parts) != 2 || !validID(parts[0]) || filepath.Base(parts[1]) != parts[1] || (filepath.Ext(parts[1]) != ".ts" && parts[1] != "index.m3u8") {
			http.NotFound(w, r)
			return
		}
		s.mu.Lock()
		v := s.videos[parts[0]]
		s.mu.Unlock()
		if v.Status != "completed" {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(parts[1], ".m3u8") {
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		} else {
			w.Header().Set("Content-Type", "video/mp2t")
		}
		http.ServeFile(w, r, filepath.Join(s.root, parts[0], parts[1]))
	default:
		http.NotFound(w, r)
	}
}
func (s *service) upload(w http.ResponseWriter, r *http.Request) {
	select {
	case <-s.ctx.Done():
		http.Error(w, "shutting down", 503)
		return
	default:
	}
	select {
	case s.slots <- struct{}{}:
	default:
		http.Error(w, "processor busy; retry later", 503)
		return
	}
	release := true
	defer func() {
		if release {
			<-s.slots
		}
	}()
	r.Body = http.MaxBytesReader(w, r.Body, 101<<20)
	if e := r.ParseMultipartForm(1 << 20); e != nil {
		http.Error(w, "invalid or oversized multipart upload", 400)
		return
	}
	defer r.MultipartForm.RemoveAll()
	f, h, e := r.FormFile("file")
	if e != nil {
		http.Error(w, "file required", 400)
		return
	}
	defer f.Close()
	if h.Size == 0 || h.Size > 100<<20 {
		http.Error(w, "file must be 1 byte to 100 MiB", 400)
		return
	}
	var idbytes [16]byte
	if _, e = rand.Read(idbytes[:]); e != nil {
		http.Error(w, "ID generation failed", 500)
		return
	}
	id := hex.EncodeToString(idbytes[:])
	dir := filepath.Join(s.root, id)
	if e = os.Mkdir(dir, 0700); e != nil {
		http.Error(w, "storage failed", 500)
		return
	}
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(dir)
		}
	}()
	out, e := os.OpenFile(filepath.Join(dir, "input"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		http.Error(w, "storage failed", 500)
		return
	}
	_, e = io.Copy(out, f)
	e = errors.Join(e, out.Close())
	if e != nil {
		http.Error(w, "upload write failed", 500)
		return
	}
	v := video{ID: id, Name: filepath.Base(h.Filename), Status: "processing"}
	if e = s.save(v); e != nil {
		http.Error(w, "metadata failed", 500)
		return
	}
	cleanup = false
	release = false
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() { <-s.slots }()
		ctx, cancel := context.WithTimeout(s.ctx, 2*time.Minute)
		defer cancel()
		err := s.transcode(ctx, filepath.Join(dir, "input"), dir)
		if err != nil {
			v.Status = "failed"
			v.Error = "transcoding failed"
			log.Printf("transcode %s: %v", id, err)
		} else {
			v.Status = "completed"
			v.URL = "/hls/" + id + "/index.m3u8"
		}
		os.Remove(filepath.Join(dir, "input"))
		if e := s.save(v); e != nil {
			log.Printf("persist status %s: %v", id, e)
			v.Status = "failed"
			v.Error = "status persistence failed"
			s.mu.Lock()
			s.videos[id] = v
			s.mu.Unlock()
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(video{ID: id, Name: filepath.Base(h.Filename), Status: "processing"})
}
func transcode(ctx context.Context, input, dir string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-nostdin", "-y", "-protocol_whitelist", "file,pipe", "-i", input, "-map", "0:v:0", "-map", "0:a:0?", "-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p", "-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2", "-c:a", "aac", "-f", "hls", "-hls_time", "2", "-hls_playlist_type", "vod", "-hls_segment_filename", filepath.Join(dir, "segment%03d.ts"), filepath.Join(dir, "index.m3u8"))
	if b, e := cmd.CombinedOutput(); e != nil {
		return fmt.Errorf("ffmpeg: %w: %.2000s", e, b)
	}
	return nil
}
func main() {
	s, e := newService("./videos")
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	log.Fatal((&http.Server{Addr: ":8080", Handler: s, ReadHeaderTimeout: 5 * time.Second}).ListenAndServe())
}
