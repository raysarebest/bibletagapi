package main

// custom error type for bases that can't be found
type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
