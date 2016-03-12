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
  QueryTopTags([]byte, string) (TagBook, TagChapter, TagVerse, error)
  QueryDBP(TagBook, TagChapter, TagVerse) (interface{}, error)
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

// TagChapter encodes tagged bible chapters for a sepecific hashtag
type TagChapter struct {
  Group      float64  `json:"group"`
  Reduction  float64    `json:"reduction"`
}

// TagChapters is a slice of TagChapter
type TagChapters []TagChapter

// QueryTopTags queries Rethink DB to get top tagged verses for a hashtag
func (dbb DBInfo) QueryTopTags(msg []byte, table string) (TagBook, TagChapter, TagVerse, error) {

  configuration := ImportConfig()

  // parse the message
  var jsonDataer interface{}
  jsonerr := json.Unmarshal(msg, &jsonDataer)
  if jsonerr != nil {
    log.Printf("%s: %s", "ERROR could not parse JSON message", jsonerr)
    return TagBook{}, TagChapter{}, TagVerse{}, jsonerr
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
    return TagBook{}, TagChapter{}, TagVerse{}, nil
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
    }).Group("chapter").Count().Run(session)
    defer res.Close()

    var row interface{}
    var tcs TagChapters
    for res.Next(&row) {
      rowMap := row.(map[string]interface{})
      tc := TagChapter{
        Group: rowMap["group"].(float64),
        Reduction: rowMap["reduction"].(float64),
      }
      tcs = append(tcs, tc)
    }
    sort.Sort(tcs)

    var tagchapter TagChapter
    if len(tcs) > 0 {

      tagchapter = tcs[0]

      // Get the Top Tagged Verse(s)
      res, err = r.DB(configuration.Dbname).Table(table).GetAllByIndex("tag", tag).Filter(map[string]interface{}{
        "book": tagbook.Group,
      }).Filter(map[string]interface{}{
        "chapter": tagchapter.Group,
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

      return tagbook, tagchapter, tagverse, nil

    }

  }

  return TagBook{}, TagChapter{}, TagVerse{}, nil

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

func (slice TagChapters) Len() int {
    return len(slice)
}

func (slice TagChapters) Less(i, j int) bool {
    return slice[i].Reduction < slice[j].Reduction;
}

func (slice TagChapters) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}