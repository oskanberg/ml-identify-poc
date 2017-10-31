package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/machinebox/sdk-go/tagbox"
	"github.com/matryer/way"
)

// Server is the app server.
type Server struct {
	assets string
	tagbox *tagbox.Client
	router *way.Router
}

// NewServer makes a new Server.
func NewServer(assets string, tagbox *tagbox.Client) *Server {
	srv := &Server{
		assets: assets,
		tagbox: tagbox,
		router: way.NewRouter(),
	}
	srv.router.Handle(http.MethodGet, "/assets/", Static("/assets/", assets))

	srv.router.HandleFunc(http.MethodPost, "/webFaceID", srv.handlewebFaceID)
	srv.router.HandleFunc(http.MethodGet, "/", srv.handleIndex)
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.assets, "index.html"))
}

func (s *Server) handlewebFaceID(w http.ResponseWriter, r *http.Request) {
	img := r.FormValue("imgBase64")
	b64data := img[strings.IndexByte(img, ',')+1:]
	imgDec, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Printf("[ERROR] Error decoding the image %v\n", err)
		http.Error(w, "can not decode the image", http.StatusInternalServerError)
		return
	}
	tags, err := s.tagbox.Check(bytes.NewReader(imgDec))
	if err != nil {
		log.Printf("[ERROR] Error on tagbox %v\n", err)
		http.Error(w, "something went wrong verifying the tags", http.StatusInternalServerError)
		return
	}

	spew.Dump(tags)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Static gets a static file server for the specified path.
func Static(stripPrefix, dir string) http.Handler {
	h := http.StripPrefix(stripPrefix, http.FileServer(http.Dir(dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
