package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"goji.io/pat"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ensureUniqueBird(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(BucketName).C(ColectionType)

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
func listBirds(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(BucketName).C(ColectionType)

		var birds []Bird
		err := c.Find(bson.M{"visible": true}).All(&birds)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all birds: ", err)
			return
		}
		//fmt.Println(birds)

		respBody, err := json.MarshalIndent(birds, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
func addBird(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var bird Bird
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&bird)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}
		if len(bird.Name) == 0 || len(bird.Family) == 0 || len(bird.Continents) == 0 {
			ErrorWithJSON(w, "Mandatory Parameters Missing ", http.StatusBadRequest)
			return
		}
		if len(bird.Added) == 0 {
			t := time.Now().UTC()
			bird.Added = t.Format("2006-01-02")
		}

		c := session.DB(BucketName).C(ColectionType)

		err = c.Insert(bird)
		if err != nil {
			//fmt.Println(err)
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Bird is already Present.Try with different One",
					http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert bird: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+bird.Name)
		w.WriteHeader(http.StatusCreated)
	}
}
func deleteBird(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		c := session.DB(BucketName).C(ColectionType)

		err := c.Remove(bson.M{"id": id})
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete bird: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Bird not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func getBirdById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		c := session.DB(BucketName).C(ColectionType)

		var bird Bird
		err := c.Find(bson.M{"id": id}).One(&bird)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find bird: ", err)
			return
		}

		if bird.Id == "" {
			ErrorWithJSON(w, "bird not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(bird, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
