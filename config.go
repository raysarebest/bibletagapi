package main

import "os"

// Configuration - config information to pass to various functions
type Configuration struct {
	Dbname            string
	Dbaddress         string
	TagPostTable      string
}

// ImportConfig instantiates the above structure info
func ImportConfig() Configuration {

	// initialize default configuration
	configuration := Configuration{
		Dbname:        "bibletagapi",
		Dbaddress:     "localhost:28015",
		TagPostTable:  "tags",
	}

	// override defaults if environmental vars available
	if os.Getenv("BIBLETAGAPI_DBNAME") != "" {
		configuration.Dbname = os.Getenv("BIBLETAGAPI_DBNAME")
	}
	if os.Getenv("BIBLETAGAPI_DBADDRESS") != "" {
		configuration.Dbaddress = os.Getenv("BIBLETAGAPI_DBADDRESS")
	}
	if os.Getenv("BIBLETAGAPI_TAGPOSTTABLE") != "" {
		configuration.TagPostTable = os.Getenv("BIBLETAGAPI_TAGPOSTTABLE")
	}

	return configuration
}