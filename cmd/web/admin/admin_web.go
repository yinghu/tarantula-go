package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func handleWeb(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Path
	if path == "/" || path == "/index.html" {
		f, err := os.Open("web/index.html")
		if err != nil {
			w.Write([]byte("no page"))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	ext := filepath.Ext(path)
	fmt.Printf("Ext %s , %s\n", ext, "web"+path)
	if ext == ".css" {

		f, err := os.Open("web" + path)
		if err != nil {
			w.Write([]byte("no page"))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	if ext == ".js" {
		f, err := os.Open("web" + path)
		if err != nil {
			w.Write([]byte("no page"))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	if ext == ".ico" {
		f, err := os.Open("web" + path)
		if err != nil {
			w.Write([]byte("no page"))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	w.Write([]byte("no page" + ext))
}
