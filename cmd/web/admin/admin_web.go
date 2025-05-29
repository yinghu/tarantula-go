package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type AdminWeb struct {
	*AdminService
}

func (s *AdminWeb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Path
	if path == "/" {
		f, err := os.Open("web/index.html")
		if err != nil {
			p404(w)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	ext := filepath.Ext(path)
	if ext == ".html" {
		f, err := os.Open("web" + path)
		if err != nil {
			p404(w)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	if ext == ".css" {

		f, err := os.Open("web" + path)
		if err != nil {
			w.Write([]byte(""))
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
			w.Write([]byte(""))
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
			w.Write([]byte(""))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	if ext == ".jpg" {
		f, err := os.Open("web" + path)
		if err != nil {
			w.Write([]byte(""))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		return
	}
	w.Write([]byte(""))
}

func p404(w http.ResponseWriter) {
	p404, err := os.Open("web/404.html")
	if err != nil {
		w.Write([]byte(""))
		return
	}
	defer p404.Close()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, p404)
}
