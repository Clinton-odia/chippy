package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHit atomic.Int64
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHit.Add(1)
		next.ServeHTTP(w, r)
	})
}
func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: serveMux,
	}
	cfg := &apiConfig{}
	serveMux.Handle("/", cfg.middlewareMetricsInc(http.FileServer(http.Dir("./app"))))
	serveMux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./app/assets")))))
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Print(err)
	}

}
