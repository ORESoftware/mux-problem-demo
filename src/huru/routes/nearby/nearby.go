package nearby

import (
	"encoding/json"
	"errors"
	"huru/models/nearby"
	"huru/mw"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// NearbyHandler just what it says
type Handler struct{}

type Injection struct {
	Nearby nearby.Map
}

var (
	mtx sync.Mutex
)

func (h Handler) Mount(router *mux.Router, v Injection) {
	mwList := mw.List(mw.Error(), mw.Auth("x-huru-api-token"))
	router.HandleFunc("/nearby", mw.Middleware(mwList, h.makeGetMany(v))).Methods("GET")
	router.HandleFunc("/nearby/{id}", mw.Middleware(mwList, h.makeGetOne(v))).Methods("GET")
	router.HandleFunc("/nearby/{id}", mw.Middleware(mwList, h.makeCreate(v))).Methods("POST")
	router.HandleFunc("/nearby/{id}", mw.Middleware(mwList, h.makeDelete(v))).Methods("DELETE")
	router.HandleFunc("/nearby/{id}", mw.Middleware(mwList, h.makeUpdate(v))).Methods("PUT")
}

// GetMany Display all from the people var
func (h Handler) makeGetMany(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("we are getting many nearby's");
		json.NewEncoder(w).Encode(v.Nearby)
	})
}

// GetOne Display a single data
func (h Handler) makeGetOne(v Injection) http.HandlerFunc {

	type APIDocs struct {
		ResolutionValue struct{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.Nearby[params["id"]]
		mtx.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	})
}

func (h Handler) makeUpdate(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)

		var t nearby.Model
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.Nearby[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("no item to update"))
		}

		if t.ContactTime != 0 {
			item.ContactTime = t.ContactTime
		}

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	})
}

// Create create a new item
func (h Handler) makeCreate(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var n nearby.Model
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.Nearby[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	})
}

// Delete Delete an item
func (h Handler) makeDelete(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		defer mtx.Unlock()

		_, deleted := v.Nearby[params["id"]]

		if deleted != true {
			json.NewEncoder(w).Encode(nil)
			return;
		}

		delete(v.Nearby, params["id"])
		json.NewEncoder(w).Encode(deleted)
	})
}
