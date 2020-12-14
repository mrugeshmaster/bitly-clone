package main

import (
	"encoding/json"
	"math"
	"net/http"
	"time"
	"os"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// GCP
var gateway = os.Getenv("gateway") 
// var gateway = "34.215.253.104:8000/lr"

func ControlPanelServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes_CP(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes_CP(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler_CP(formatter)).Methods("GET")
	mx.HandleFunc("/createshortlink", createNewShortLinkHandler(formatter)).Methods("POST")
}

// API Ping Handler
func pingHandler_CP(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Control Panel API version 1.0 alive!"})
	}
}

// API Create New Short Link Handler
func createNewShortLinkHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var jsonData map[string]string
		json.NewDecoder(req.Body).Decode(&jsonData)
		longurl := jsonData["url"]
		shortlink_code := generateShortLinkCode()
		
		queueBody := "{ \"shortlink_code\": \"" + shortlink_code + "\", \"longurl\": \"" + longurl + "\"}"
		shortLinkCreatequeue_send(queueBody)

		shorturl := "http://" + gateway + "/" + shortlink_code
		respBody := "{ shortlink_code: " + shorturl + ", longurl: " + longurl + "}"
		formatter.JSON(w, http.StatusOK, respBody)
	}
}

func generateShortLinkCode() string {
	timestamp := time.Now().UTC().UnixNano()
	n := int(timestamp)

	characterSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234556789"
	base := 62

	b := make([]byte, 0, 7)
	for n > 0 {
		r := math.Mod(float64(n), float64(base))
		n /= base
		if len(b) < 8 {
			b = append([]byte{characterSet[int(r)]}, b...)
		}
	}

	return string(b)
}
