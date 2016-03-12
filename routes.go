package main

import "net/http"

// Route - used to pass information about a particular route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - used to pass information about multiple routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	// route to processed new tags for verses
	Route{
		"PostTag",
		"POST",
		"/tag",
		http.HandlerFunc(myTagHandler(PostTag, DBInfo{})),
	},
	// route to return verses based on tags
	Route{
		"RetrieveTag",
		"PUT",
		"/tag",
		http.HandlerFunc(myRetrieveHandler(RetrieveTag, DBInfo{})),
	},
}
