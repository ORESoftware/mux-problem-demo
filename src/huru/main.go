package main

import (
	"encoding/json"
	"fmt"
	"huru/models"
	"huru/mw"
	"huru/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)


func main() {


	routerParent := mux.NewRouter()  // root router
	routerParent.Use(mw.ErrorHandler())
	routerParent.Use(mw.LoggingHandler())

	routerParent.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info("the logging handling middleware should have logged something before this.")
		panic("this should get trapped by error handler 1.")

	}))


	routerParent.HandleFunc("/dogs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info("the logging handling middleware should have logged something before this.")
		panic("this should get trapped by error handler 2.")

	}))

	//routerParent.Use(mw.Auth())
	//subRouter := mux.NewRouter()
	//subRouter.Use(loggingMiddleware)
	//subRouter.Use(errorMiddleware)
	//subRouter.Use(authMiddleware)

	router := routerParent.PathPrefix("/api/v1").Subrouter();
	//router.Use(loggingMiddleware)
	//router.Use(errorMiddleware)
	//router.Use(authMiddleware)

	router.Use(mw.ErrorHandler())  // we call this again
	router.Use(mw.LoggingHandler()) // we call this again


	router.HandleFunc("/dogs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info("the logging handling middleware should have logged something before this.")
		panic("this should get trapped by error handler 3.")

	}))


	{
		handler := routes.PersonDataHandler{}
		//subRouter := router.PathPrefix("/").Subrouter()
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

		subRouter.Use(mw.ErrorHandler())   // we call this again
		subRouter.Use(mw.LoggingHandler())  // we call this again

		subRouter.HandleFunc("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// *>*>* the middleware doesn't get hit for this route
			log.Info("the logging handling middleware should have logged something before this.")
			panic("this should get trapped by error handler 4.")

		}))

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
