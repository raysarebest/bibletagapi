package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"

	"github.com/gorilla/mux"
)

// Index is the handler for the root URL
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to this example API\n")
}

func ReturnOptions(w http.ResponseWriter, r *http.Request) { 

	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Methods", "POST") 
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization") 

	// Stop here if its Preflighted OPTIONS request 
	if r.Method == "OPTIONS" { 
		return 
	}
}

// DBInfo provides functionality for handlers accessing DBs or other datasources
type DBInfo struct{}

type dbTagHandler func(w http.ResponseWriter, r *http.Request, dbb TagInfoer)

// myTagHandler generates a handler accesing DBs or other datasources
func myTagHandler(handler dbTagHandler, dbb TagInfoer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, dbb)
	}
}

// PostTag is the handler that processes new tags for Bible verses
// and stores the tag in RethinkDB
func PostTag(w http.ResponseWriter, r *http.Request, dbb TagInfoer) {

	configuration := ImportConfig()

	// read the JSON body into a string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s: %s", "ERROR Could not read request body", err.Error())
	}

	// verify if the JSON body includes the minimum requirements for a tag post
	isvalid := isPostValid(body)

	// if a valid post, push to RethinkDB
	if isvalid {
		rethinkerr := dbb.PostRethink(body, configuration.TagPostTable)
		if rethinkerr == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusOK, Text: "Tagged Passage"}); err != nil {
				log.Printf("%s: %s", "ERROR could not encode JSON response", err.Error())
			}
			return
		}
	}

	// If we didn't find it, 304
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNotModified)
	return

}

type dbRetrieveHandler func(w http.ResponseWriter, r *http.Request, dbb RetrieveInfoer)

// myRetrieveHandler generates a handler accesing DBs or other datasources
func myRetrieveHandler(handler dbRetrieveHandler, dbb RetrieveInfoer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, dbb)
	}
}

// RetrieveTag is the handler that returns bible content (from DBP) based on tags
func RetrieveTag(w http.ResponseWriter, r *http.Request, dbb RetrieveInfoer) {

	configuration := ImportConfig()

	// parsed variables from the router
	vars := mux.Vars(r)
	currenttag := vars["currenttag"]

	// Get the top tagged book and verses
	tagbook, tagchapter, tagverse, err := dbb.QueryTopTags(currenttag, configuration.TagPostTable)
	if err != nil {
		log.Printf("%s: %s", "ERROR could retrieve top tags for hashtag", err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var dbpmsg interface{}
	if len(tagverse.Group) > 0 {

		// Get JSON response from DBP for tagbook, tagverse
		dbpmsg, err = dbb.QueryDBP(tagbook, tagchapter, tagverse)
		if err != nil {
			log.Printf("%s: %s", "ERROR could retrieve DBP content for hashtag", err.Error())
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(dbpmsg); err != nil {
			log.Printf("%s: %s", "ERROR could not encode JSON response", err.Error())
		}
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
	return

}

func isPostValid(body []byte) bool {

	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		log.Printf("%s: %s", "ERROR Could not unmarshal JSON into interface", err.Error())
	}
	m := f.(map[string]interface{})

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	if !StringContains(keys, "tag") {
		return false
	}
	if !StringContains(keys, "book")  && !StringContains(keys, "chapter") {
		return false
	}
	if !StringContains(keys, "startVerse") && !StringContains(keys, "endVerse") {
		return false
	}

	if _, ok := m["chapter"].(string); ok {
        return false
    }
    if _, ok := m["startVerse"].(string); ok {
        return false
    }
    if _, ok := m["endVerse"].(string); ok {
        return false
    }

	return true

}

func StringContains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
        return true
    }
  }
  return false
}