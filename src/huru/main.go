package main

import (
	"encoding/json"
	"fmt"
	"huru/mw"
	"io"

	"huru/models"
	"huru/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// https://www.cyberciti.biz/faq/howto-add-postgresql-user-account/


func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

const (
	MyAPIKey = "MY_EXAMPLE_KEY"
)

// AuthMiddleware is an example of a middleware layer. It handles the request authorization
// by checking for a key in the url.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestKey := r.URL.Query().Get("key")
		if len(requestKey) == 0 || requestKey != MyAPIKey {
			// Report Unauthorized
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `{"error":"invalid_key"}`)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

func main() {


	routerParent := mux.NewRouter()
	routerParent.Use(mw.ErrorHandler())
	routerParent.Use(mw.LoggingHandler())
	//routerParent.Use(mw.Auth())

	//subRouter := mux.NewRouter()
	//subRouter.Use(loggingMiddleware)
	//subRouter.Use(errorMiddleware)
	//subRouter.Use(authMiddleware)

	router := routerParent.PathPrefix("/api/v1").Subrouter();
	//router.Use(loggingMiddleware)
	//router.Use(errorMiddleware)
	//router.Use(authMiddleware)


	{
		handler := routes.PersonDataHandler{}
		subRouter := router.PathPrefix("/").Subrouter()
		//subRouter.Use(mw.Auth)
		subRouter.Use(mw.AuthHandler())
		handler.Mount(subRouter, routes.PersonDataInjection{DataFields: models.DataFieldsInit()})
	}


	{
		// nearby
		handler := routes.NearbyHandler{}
		subRouter := router.PathPrefix("/").Subrouter()
		//subRouter.Use(middleware.AuthHandler())
		handler.Mount(subRouter, routes.NearbyInjection{Nearby: models.NearbyInit()})
	}

	{
		// share
		handler := routes.ShareHandler{}
		subRouter := router.PathPrefix("/").Subrouter()
		//subRouter.Use(mw.ErrorHandler());
		handler.Mount(subRouter, routes.ShareInjection{Share: models.ShareInit()})
	}


	handler404 := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(404)
		json.NewEncoder(w).Encode(struct {
			Message string
		}{"Could not find a matching handler on the API server."})

	}

	router.NotFoundHandler = http.HandlerFunc(handler404);
	routerParent.NotFoundHandler = http.HandlerFunc(handler404)


	host := os.Getenv("huru_api_host")
	port := os.Getenv("huru_api_port")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "80"
	}

	log.Info(fmt.Sprintf("Huru API server listening on port %s", port))
	path := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(path, routerParent))

}
