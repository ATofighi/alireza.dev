package main

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPHandler struct {
}

func (h *HTTPHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, r.Method, r.URL)
	for k, v := range r.Header {
		fmt.Fprintf(rw, "%s: %+v\n", k, v)
	}
	fmt.Fprintln(rw)
}

func main() {
	err := http.ListenAndServe("127.0.0.1:9000", &HTTPHandler{})
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
