package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jamesk14022/barcrawler/cache"
	"github.com/jamesk14022/barcrawler/handlers"

	gorillaHandlers "github.com/gorilla/handlers"
)

var staticDir = os.Getenv("WEB_DIR")

const port = ":8080"

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFun7c(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type routeHandler func(http.ResponseWriter, *http.Request) error

// ServeHTTP allows our routeHandler type to satisfy http.Handler interface.
func (fn routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {

	cache.InitCache()

	router := mux.NewRouter()

	router.Handle("/cities", routeHandler(handlers.GetCityCoordinates)).Methods("GET")
	router.Handle("/pubs", routeHandler(handlers.GetRandomCrawl)).Methods("GET")
	router.Handle("/citypoints", routeHandler(handlers.GetAllCityPoints)).Methods("GET")
	router.Handle("/crawl", routeHandler(handlers.PostCrawl)).Methods("POST")
	router.Handle("/photo", routeHandler(handlers.GetPhoto)).Methods("GET")

	router.
		PathPrefix("/").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	corsRouter := enableCORS(router)
	loggedRouter := gorillaHandlers.LoggingHandler(os.Stdout, corsRouter)

	server := &http.Server{
		Addr:    port,
		Handler: loggedRouter,
	}

	log.Println("Starting Barcrawler Server on ", port)
	log.Fatal(server.ListenAndServe())
}
