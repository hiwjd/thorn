package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`你好`))
	})
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`你好a`))
	})
	http.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`你好b`))
	})
	http.ListenAndServe("0.0.0.0:8080", nil)
}
