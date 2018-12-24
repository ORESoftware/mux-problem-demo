package helpers

import (
	"github.com/gorilla/mux"
	"net/http"
	tc "huru/type-creator"
)

type Injection interface {
	GetVal() interface{}
}

type RouteMaker = func(d tc.Docs, v interface{}) http.HandlerFunc
type RouteAdder = func(r string, methods []string, f RouteMaker)

func MakeRouteAdder(router *mux.Router, v interface{}) RouteAdder {
	return func(r string, methods []string, f RouteMaker) {
		d := tc.Docs{Methods:methods, Route:r}
		router.HandleFunc(r, f(d, v)).Methods(methods...)
	}
}
