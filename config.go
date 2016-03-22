package main

import "os"

// Configuration - config information to pass to various functions
type Configuration struct {
	Dbname            string
	Dbaddress         string
	TagPostTable      string
	DBPAPIKey					string
}

// ImportConfig instantiates the above structure info
func ImportConfig() Configuration {

	// initialize default configuration
	configuration := Configuration{
		Dbname:        "bibletagapi",
		Dbaddress:     "localhost:28015",
		TagPostTable:  "tags",
		DBPAPIKey:		 "",
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
	if os.Getenv("BIBLETAGAPI_DBP_API_KEY") != "" {
		configuration.DBPAPIKey = os.Getenv("BIBLETAGAPI_DBP_API_KEY")
	}

	return configuration
}