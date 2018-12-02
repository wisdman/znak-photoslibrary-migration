package main

import (
  "log"
  "os"
  "encoding/json"
  "io/ioutil"
  "github.com/jackc/pgx"
  "fmt"
)

var DB *pgx.ConnPool

type Album struct {
  Id         string `json:"id"`
  Title      string `json:"title"`
  ProductUrl string `json:"productUrl"`
  ItemsCount int    `json:"itemsCount"`
  Error      bool   `json:"error"`
}


func main() {
  config, err := pgx.ParseEnvLibpq()
  if err != nil {
    log.Fatal(err)
  }

  DB, err = pgx.NewConnPool(pgx.ConnPoolConfig{
    ConnConfig:     config,
    MaxConnections: 10,
  })
  if err != nil {
    log.Fatal(err)
  }

  jsonFile, err := os.Open("albums.json")
  if err != nil {
    log.Fatal(err)
  }
  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    log.Fatal(err)
  }

  albums := make([]Album,0)

  err = json.Unmarshal(byteValue, &albums)
  if err != nil {
    log.Fatal(err)
  }

  for i, album := range albums {
    fmt.Printf("\rAlbum: %d", i)
    _, err := DB.Exec("INSERT INTO google (id, titel, productUrl, itemsCount, error) VALUES ($1, $2, $3, $4, $5)",
                      album.Id, album.Title, album.ProductUrl, album.ItemsCount, album.Error)
    if err != nil {
      log.Fatal(err)
    }
  }

}
