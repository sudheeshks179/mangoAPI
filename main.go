package main

import (
	"fmt"
	"net/http"

	goji "goji.io"
	"goji.io/pat"
	mgo "gopkg.in/mgo.v2"
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

var BucketName string = "cage"
var ColectionType string = "birds"

type Bird struct {
	Name       string   `json:"name"`
	Family     string   `json:"family"`
	Continents []string `json:"continents"`
	Added      string   `json:"added"`
	Visible    bool     `json:"visible"`
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureUniqueBird(session)

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/birds"), listBirds(session))
	mux.HandleFunc(pat.Post("/birds"), addBird(session))
	mux.HandleFunc(pat.Get("/birds/:id"), getBirdById(session))
	mux.HandleFunc(pat.Delete("/birds/:id"), deleteBird(session))
	http.ListenAndServe("localhost:8080", mux)
}
