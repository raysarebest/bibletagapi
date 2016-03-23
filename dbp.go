package main

import (
  "log"
  "net/http"
  "fmt"
  "bufio"
  "os"
  "io/ioutil"
  "encoding/json"
)

// QueryDBP get bible content from DBP for a hashtag
func (dbb DBInfo) QueryDBP(tb TagBook, tc TagChapter, tv TagVerse) (interface{}, error) {

  book := tb.Group
  sverse := tv.Group[0]
  everse := tv.Group[1]
  chapter := tc.Group

  // determine if the reference is in the OT or NT
  var testament string
  otbooks, err := ReadLines("files/ot.csv")
  if err != nil {
    log.Printf("%s: %s", "ERROR Could not read in OT book IDs", err.Error())
  }
  if StringContains(otbooks, book) {
    testament = "O"
  }
  ntbooks, err := ReadLines("files/nt.csv")
  if err != nil {
    log.Printf("%s: %s", "ERROR Could not read in NT book IDs", err.Error())
  }
  if StringContains(ntbooks, book) {
    testament = "N"
  }

  apikey := ImportConfig().DBPAPIKey
  url := fmt.Sprintf(`http://dbt.io/text/verse?key=%s&dam_id=ENGESV%s2ET&book_id=%s&chapter_id=%d&verse_start=%d&verse_end=%d&v=2`,
    apikey, testament, book, int(chapter), int(sverse), int(everse))

  res, err := http.Get(url)
  if err != nil {
    log.Printf("%s: %s", "ERROR could not retrieve DBP content for hashtag", err.Error())
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Printf("%s: %s", "ERROR Could not read DBP response body", err.Error())
  }

  var jsonDataer interface{}
  json.Unmarshal(body, &jsonDataer)

  return jsonDataer, err

}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}