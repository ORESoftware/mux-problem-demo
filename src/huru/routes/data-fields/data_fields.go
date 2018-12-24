package data_fields

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"huru/models/person-data-fields"
	"huru/mw"
	"huru/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Handler - ShareHandler just what it says
type Handler struct{}

type Injection struct {
	DataFields person_data_fields.Map
}

var (
	mtx sync.Mutex
)

func (h Handler) Mount(router *mux.Router, v Injection) Handler {

	mwList := mw.List(mw.Error(), mw.Auth("x-huru-api-token"))

	router.HandleFunc("/person_data_field", mw.Middleware(mwList, h.makeGetMany(v))).Methods("GET")
	router.HandleFunc("/person_data_field/{id}", mw.Middleware(mwList, h.makeGetOne(v))).Methods("GET")
	router.HandleFunc("/person_data_field/{id}", mw.Middleware(mwList, h.makeCreate(v))).Methods("POST")
	router.HandleFunc("/person_data_field/{id}", mw.Middleware(mwList, h.makeDelete(v))).Methods("DELETE")
	router.HandleFunc("/person_data_field/{id}", mw.Middleware(mwList, h.makeUpdateByID(v))).Methods("PUT")

	return h
}

func (h Handler) makeGetMany(v Injection) http.HandlerFunc {
	return mw.Middleware(
		//mw.Error(),
		//mw.Auth("x-huru-api-token"),
		mw.Allow("admin"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("now we are sending response.");
			json.NewEncoder(w).Encode(v.DataFields)
		}),
	)
}

// MakeGetOne Display a single data
func (h Handler) makeGetOne(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.DataFields[params["id"]]
		mtx.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	})
}

// MakeCreate create a new item
func (h Handler) makeCreate(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var n person_data_fields.Model
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.DataFields[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	})
}

// MakeDelete Delete an item
func (h Handler) makeDelete(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, deleted := v.DataFields[params["id"]]
		delete(v.DataFields, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(deleted)
	})
}

func (h Handler) makeUpdateByID(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t person_data_fields.Model
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.DataFields[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("no item to update"))
		}

		if t.Key != "" {
			if t.Key != item.Key {
				panic(utils.AppError{
					StatusCode: 409,
					Message:    utils.JoinArgs("FieldName does not match, see: ", t.Key, item.Value),
				})
			}
		}

		item.Value = t.Value

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	})
}
