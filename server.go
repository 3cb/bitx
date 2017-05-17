package main

import (
	"fmt"
	"log"
	"net/http"

	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	// create Gorilla mux router
	r := mux.NewRouter()

	// list routes here
	r.Handle("/api/getHistory/{currency}", historyHandler())

	// use negroni to serve static files
	n := negroni.New(negroni.NewStatic(http.Dir("./static/")))
	n.UseHandler(r)

	// start server
	log.Fatal(http.ListenAndServe(":3000", n))
}

func historyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		currency := vars["currency"]
		api := fmt.Sprintf("https://api.gemini.com/v1/trades/%v?limit_trades=500", currency)
		resp, err := http.Get(api)
		if err != nil {
			log.Printf("Error during http request to %v: %v", api, err)
		}
		data, err2 := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err2 != nil {
			log.Printf("Error parsing response body: %v", err2)
		}
		w.Write(data)
	})
}
