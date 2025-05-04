package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"gameclustering.com/internal/auth"
)

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

type actionHandler struct {
	mu sync.Mutex // guards n
	n  int
}


func (h *countHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func (h *actionHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func debugging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Println(r.URL.Path, time.Since(start))
		}()
		f(w, r)
	}
}

func bootstrap(host string) {
	http.Handle("/count", new(countHandler))
	http.Handle("/action", new(actionHandler))
	http.Handle("/auth",http.HandlerFunc(debugging(auth.AuthHandler)))
	log.Fatal(http.ListenAndServe(host, nil))
}
