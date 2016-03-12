package main

import (
  "log"
  "encoding/json"

  r "github.com/dancannon/gorethink"
)

var session *r.Session

// TagInfoer is a one method interface agent
type TagInfoer interface {
  PostRethink([]byte, string) error
}


// PostRethink posts JSON data to rethinkdb
func (dbb DBInfo) PostRethink(msg []byte, table string) error {

  configuration := ImportConfig()

  // parse the message
  var jsonDataer interface{}
	jsonerr := json.Unmarshal(msg, &jsonDataer)
  if jsonerr != nil {
    log.Printf("%s: %s", "ERROR could not parse JSON message", jsonerr)
    return jsonerr
  }
  m := jsonDataer.(map[string]interface{})

  // connect to Rethink
  session, err := r.Connect(r.ConnectOpts{
      Address: configuration.Dbaddress,
  })
  if err != nil {
    log.Printf("%s: %s", "ERROR could not connect to RethinkDB", err)
    return err
  }

  // push to Rethink
  query := r.DB(configuration.Dbname).Table(table).Insert(m)
	_, err = query.Run(session)
  if err != nil {
    log.Printf("%s: %s", "ERROR could not push data to RethinkDB", err)
    return err
  }

  return nil

}