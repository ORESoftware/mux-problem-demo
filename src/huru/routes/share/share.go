package share

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"huru/models/share"
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
	Share share.Map
}

var (
	mtx sync.Mutex
)

func (h Handler) Mount(router *mux.Router, v Injection) Handler {

	//router.Use(mw.Auth(), mw.Logging(), mw.Error());
	//router.Handle("/share", h.makeGetMany(v)).Methods("GET")
	//router.Handle("/share", getMiddleware(h.makeGetMany(v))).Methods("GET")
	//router.HandleFunc("/share", getMiddleware(h.makeGetMany(v))).Methods("GET")

	//logging := mw.Logging(struct{}{})
	//auth := mw.Auth("x-huru-api-token")
	//errorh := mw.Error()

	//router.HandleFunc("/share", mw.AdaptFuncs(h.makeGetMany(v), auth, errorh, logging)).Methods("GET")

	mwList := mw.List(mw.Error(), mw.Auth("x-huru-api-token"))

	router.HandleFunc("/share", mw.Middleware(mwList, h.makeGetMany(v))).Methods("GET")
	router.HandleFunc("/share/{id}", mw.Middleware(mwList,h.makeGetOne(v))).Methods("GET")
	router.HandleFunc("/share/{id}", mw.Middleware(mwList,h.makeCreate(v))).Methods("POST")
	router.HandleFunc("/share/{id}", mw.Middleware(mwList,h.makeDelete(v))).Methods("DELETE")
	router.HandleFunc("/share/{id}", mw.Middleware(mwList,h.makeUpdateByID(v))).Methods("PUT")

	return h
}


func (h Handler) makeGetMany(v Injection) http.HandlerFunc {
	return mw.Middleware(
		//mw.Error(),
		//mw.Auth("x-huru-api-token"),
		//mw.Allow("admin"),  /// HERE IS ACL / ROLES
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("now we are sending response.");
			json.NewEncoder(w).Encode(v.Share)
		}),
	)
}

// MakeGetOne Display a single data
func (h Handler) makeGetOne(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.Share[params["id"]]
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
		var n share.Model
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.Share[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	})
}

// MakeDelete Delete an item
func (h Handler) makeDelete(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, deleted := v.Share[params["id"]]
		delete(v.Share, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(deleted)
	})
}

func (h Handler) makeUpdateByID(v Injection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t share.Model
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.Share[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("no item to update"))
		}

		if t.FieldName != "" {
			if t.FieldName != item.FieldName {
				panic(utils.AppError{
					StatusCode: 409,
					Message:    utils.JoinArgs("FieldName does not match, see: ", t.FieldName, item.FieldName),
				})
			}
		}

		item.FieldValue = t.FieldValue

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	})
}
