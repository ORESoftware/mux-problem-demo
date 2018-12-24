package person

import (
	"encoding/json"
	"errors"
	"huru/models/person"
	"huru/mw"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	tc "huru/type-creator"
)

// Handler just what it says
type Handler struct{}

type Injection struct {
	PeopleById, PeopleByHandle, PeopleByEmail person.Map
}

func (i Injection) GetVal() interface{} {
	return i
}

var (
	mtx sync.Mutex
)

func (h Handler) Mount(router *mux.Router, v Injection) Handler {


	//add := helpers.MakeRouteAdder(router, v)
	//add("/person", []string{"GET"}, h.makeGetMany)
	//add("/person/{id}", []string{"GET"}, h.makeGetOne)
	//add("/person/{id}", []string{"POST"}, h.makeCreate)
	//add("/person/{id}", []string{"DELETE"}, h.makeDelete)
	//add("/person/{id}", []string{"PUT"}, h.makeUpdateByID)

	d := tc.Docs{}

	mwList := mw.List(mw.Error(), mw.Auth("x-huru-api-token"))
	router.HandleFunc("/person", mw.Middleware(mwList, h.makeGetMany(d, v))).Methods("GET")
	router.HandleFunc("/person/{id}", mw.Middleware(mwList, h.makeGetOne(d, v))).Methods("GET")
	router.HandleFunc("/person/{id}", mw.Middleware(mwList, h.makeCreate(d, v))).Methods("POST")
	router.HandleFunc("/person/{id}", mw.Middleware(mwList, h.makeDelete(d, v))).Methods("DELETE")
	router.HandleFunc("/person/{id}", mw.Middleware(mwList, h.makeUpdateByID(d, v))).Methods("PUT")

	return h
}

func (h Handler) makeGetMany(d tc.Docs, v Injection) http.HandlerFunc {

	type RespBody struct {
		TypeCreatorMeta string `tc_resp_body_type:"1" type:"bar"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `tc_req_body_type:"1" type:"star"`
		Handle          string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(v.PeopleById)
		}))
}

func (h Handler) makeGetOne(d tc.Docs, x interface{}) http.HandlerFunc {

	v, ok := x.(Injection);
	if !ok {
		log.Fatal("Could not convert to Injection.");
	}

	type RespBody struct {
		TypeCreatorMeta string `type:"bar" tc_resp_body_type:"1"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `type:"star" tc_req_body_type:"1"`
		Handle          string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			mtx.Lock()
			item, ok := v.PeopleById[params["id"]]
			mtx.Unlock()
			if ok {
				json.NewEncoder(w).Encode(item)
			} else {
				io.WriteString(w, "null")
			}
		}))
}

func (h Handler) makeCreate(d tc.Docs, x interface{}) http.HandlerFunc {

	v, ok := x.(Injection)

	if !ok {
		log.Fatal("Could not convert to Injection.");
	}

	type RespBody struct {
		TypeCreatorMeta string `type:"bar" tc_resp_body_type:"1"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `type:"star" tc_req_body_type:"1"`
		Handle          string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var n person.Model
			json.NewDecoder(r.Body).Decode(&n)
			mtx.Lock()
			v.PeopleById[strconv.Itoa(n.ID)] = n
			mtx.Unlock()
			json.NewEncoder(w).Encode(&n)
		}))
}

func (h Handler) makeDelete(d tc.Docs, x interface{}) http.HandlerFunc {

	v, ok := x.(Injection);
	if !ok {
		log.Fatal("Could not convert to Injection.");
	}

	type RespBody struct {
		TypeCreatorMeta string `type:"bar" tc_resp_body_type:"1"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `type:"star" tc_req_body_type:"1"`
		Handle          string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			mtx.Lock()
			_, isDeletable := v.PeopleById[params["id"]]
			delete(v.PeopleById, params["id"])
			mtx.Unlock()
			json.NewEncoder(w).Encode(isDeletable)
		}))
}

func (h Handler) makeUpdateByID(d tc.Docs, x interface{}) http.HandlerFunc {

	v, ok := x.(Injection);
	if !ok {
		log.Fatal();
	}

	type RespBody struct {
		TypeCreatorMeta string `type:"bar" tc_resp_body_type:"1"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `type:"star" tc_req_body_type:"1"`
		Handle          string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			params := mux.Vars(r)
			decoder := json.NewDecoder(r.Body)

			var t person.Model
			err := decoder.Decode(&t)

			if err != nil {
				panic(err)
			}

			mtx.Lock()
			defer mtx.Unlock()

			item, ok := v.PeopleById[params["id"]]

			if ok == false {
				panic(errors.New("no item to update"))
			}

			//by_handle := v.PeopleByHandle[item.Handle]

			if t.Handle != "" {
				item.Handle = t.Handle
			}

			if t.Work != "" {
				item.Work = t.Work
			}

			if t.Image != "" {
				item.Image = t.Image
			}

			if t.Firstname != "" {
				item.Firstname = t.Firstname
			}

			if t.Lastname != "" {
				item.Lastname = t.Lastname
			}

			if t.Email != "" {
				item.Email = t.Email
			}

			if ok {
				json.NewEncoder(w).Encode(item)
			} else {
				io.WriteString(w, "null")
			}
		}))
}
