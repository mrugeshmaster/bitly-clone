package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func LinkRedirectServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes_LR(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes_LR(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/{shortlinkcode:[\\w]{8}$}", redirectLinkHandler(formatter)).Methods("GET")
	mx.HandleFunc("/linkstats", statsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/linkstats", statsHandler(formatter)).Queries("metric_unit", "{metric_unit}", "metric_value", "{metric_value}").Methods("GET")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Link Redirect API version 1.0 alive!"})
	}
}

// API Redirect Short URL to Long URL Handler
func redirectLinkHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortlink_code := vars["shortlinkcode"]
		body := "{ \"shortlink_code\": \"" + shortlink_code + "\" }"
		// check in cache, return if available
		var longURL string
		if longURL = getLongURLFromCache(shortlink_code); longURL == "" {
			// fetch from database
			log.Println("Cache Miss")
			longURL = getLongURLFromDB(shortlink_code)
		}

		log.Println(longURL)

		redirectLinkqueue_send(body)

		http.Redirect(w, r, longURL, 302)
	}
}

func statsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		metric_unit := req.FormValue("metric_unit")
		metric_value := req.FormValue("metric_value")

		if metric_unit != "" && metric_value != "" {
			metric_value_int, _ := strconv.Atoi(metric_value)
			getStatistics(metric_unit, metric_value_int)
		}
		fmt.Println("Fetch Statistics")
		statistics := getStatistics()
		formatter.JSON(w, http.StatusOK, statistics)
	}
}
