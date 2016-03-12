package main

import (
  "log"
  "encoding/json"
  "sort"

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

// RetrieveInfoer is a one method interface agent
type RetrieveInfoer interface {
  QueryTopTags([]byte, string) (TagBook, TagVerse, error)
}

// TagBook encodes tagged bible books for a sepecific hashtag
type TagBook struct {
  Group      string  `json:"group"`
  Reduction  float64  `json:"reduction"`
}

// TagBooks is a slice of TagBook
type TagBooks []TagBook

// TagVerse encodes tagged bible verses for a sepecific hashtag
type TagVerse struct {
  Group      []float64  `json:"group"`
  Reduction  float64    `json:"reduction"`
}

// TagVerses is a slice of TagVerse
type TagVerses []TagVerse

// QueryTopTags queries Rethink DB to get top tagged verses for a hashtag
func (dbb DBInfo) QueryTopTags(msg []byte, table string) (TagBook, TagVerse, error) {

  configuration := ImportConfig()

  // parse the message
  var jsonDataer interface{}
  jsonerr := json.Unmarshal(msg, &jsonDataer)
  if jsonerr != nil {
    log.Printf("%s: %s", "ERROR could not parse JSON message", jsonerr)
    return TagBook{}, TagVerse{}, jsonerr
  }
  m := jsonDataer.(map[string]interface{})

  // parse the requested tags
  tag := m["tag"].(string)

  // connect to Rethink
  session, err := r.Connect(r.ConnectOpts{
      Address: configuration.Dbaddress,
  })
  if err != nil {
    log.Printf("%s: %s", "ERROR could not connect to RethinkDB", err)
    return TagBook{}, TagVerse{}, nil
  }

  // Get the Top Tagged Book
  res, err := r.DB(configuration.Dbname).Table(table).GetAllByIndex("tag", tag).Group("book").Count().Run(session)
  defer res.Close()

  var row interface{}
  var tbs TagBooks
  for res.Next(&row) {
    rowMap := row.(map[string]interface{})
    tb := TagBook{
      Group: rowMap["group"].(string),
      Reduction: rowMap["reduction"].(float64),
    }
    tbs = append(tbs, tb)
  }
  sort.Sort(tbs)
  if len(tbs) > 0 {

    tagbook := tbs[0]

    // Get the Top Tagged Verse(s)
    res, err = r.DB(configuration.Dbname).Table(table).GetAllByIndex("tag", tag).Filter(map[string]interface{}{
      "book": tagbook.Group,
    }).Group("startVerse", "endVerse").Count().Run(session)
    defer res.Close()

    var row interface{}
    var tvs TagVerses
    for res.Next(&row) {
      rowMap := row.(map[string]interface{})
      verses := rowMap["group"].([]interface{})
      var vs []float64
      for _, v := range verses {
        vs = append(vs, v.(float64))
      }
      tv := TagVerse{
        Group: vs,
        Reduction: rowMap["reduction"].(float64),
      }
      tvs = append(tvs, tv)
    }
    sort.Sort(tvs)

    var tagverse TagVerse
    if len(tvs) > 0 {
      tagverse = tvs[0]
    }

    return tagbook, tagverse, nil

  }

  return TagBook{}, TagVerse{}, nil

}

func (slice TagBooks) Len() int {
    return len(slice)
}

func (slice TagBooks) Less(i, j int) bool {
    return slice[i].Reduction < slice[j].Reduction;
}

func (slice TagBooks) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func (slice TagVerses) Len() int {
    return len(slice)
}

func (slice TagVerses) Less(i, j int) bool {
    return slice[i].Reduction < slice[j].Reduction;
}

func (slice TagVerses) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}