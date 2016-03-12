package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
)

// Index is the handler for the root URL
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to this example API\n")
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
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusOK, Text: "Tagged Passage"}); err != nil {
				log.Printf("%s: %s", "ERROR could not encode JSON response", err.Error())
			}
			return
		}
	}

	// If we didn't find it, 304
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotModified)

}

// myRetrieveHandler generates a handler accesing DBs or other datasources
func myRetrieveHandler(handler dbRetrieveHandler, dbb RetrieveInfoer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, dbb)
	}
}

// RetrieveTag is the handler that returns bible content (from DBP) based on tags
func RetrieveTag(w http.ResponseWriter, r *http.Request, dbb RetrieveInfoer) {

	configuration := ImportConfig()

	// read the JSON body into a string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s: %s", "ERROR Could not read request body", err.Error())
	}

	tagbook, tagverse, err := dbb.QueryTopTags(body, configuration.TagPostTable)
	if err != nil {
		log.Printf("%s: %s", "ERROR could retrieve top tags for hashtag", err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNoContent
	}

	// Get JSON response from DBP for tagbook, tagverse

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusOK, Text: "Cool"}); err != nil {
		log.Printf("%s: %s", "ERROR could not encode JSON response", err.Error())
	}
	return

	// If we didn't find it, 204
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)

}

func isPostValid(body []byte) bool {

	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		log.Printf("%s: %s", "ERROR Could not unmarshall JSON into interface", err.Error())
	}
	m := f.(map[string]interface{})

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	if !stringContains(keys, "tags") {
		return false
	}
	if !stringContains(keys, "book")  && !stringContains(keys, "chapter") {
		return false
	}
	if !stringContains(keys, "startVerse") && !stringContains(keys, "endVerse") {
		return false
	}

	return true

}

func stringContains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
        return true
    }
  }
  return false
}