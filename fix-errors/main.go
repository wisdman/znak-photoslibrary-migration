package main

import (
  "context"
  "os"
  "log"
  "fmt"
  "encoding/json"

  "github.com/wisdman/znak-photoslibrary-migration/lib/oauth"
  photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const OAuthScope = "https://www.googleapis.com/auth/photoslibrary"

type Album struct {
  Id         string `json:"id"`
  Title      string `json:"title"`
  ProductUrl string `json:"productUrl"`
  ItemsCount int    `json:"itemsCount"`
  Error      bool   `json:"error"`
}

func main() {
  albumsFile, err := os.OpenFile("albums.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
    log.Fatalf("error opening file: %v", err)
  }
  defer albumsFile.Close()

  ctx := context.Background()
  client, err := oauth.New(ctx, []string{OAuthScope})
  if err != nil {
    log.Fatal(err)
  }

  photos, err := photoslibrary.New(client)
  if err != nil {
    log.Fatal(err)
  }

  var i uint64
  var nextPageToken string = ""

  for {
    albumsReq := photos.Albums.List().PageSize(50)
    if nextPageToken != "" {
      albumsReq = albumsReq.PageToken(nextPageToken)
    }

    albums, err := albumsReq.Do()
    if err != nil {
      log.Printf("\nALBUM LIST REQUEST ERROR: %v\n", err)
      continue
    }

    for _, album := range albums.Albums {
      i++
      fmt.Printf("\rAlbum: %d", i)

      jAlbum := &Album{Id: album.Id, Title: album.Title, ProductUrl: album.ProductUrl, ItemsCount:0, Error:false}
      var pageToken string = ""
      for {
        r := &photoslibrary.SearchMediaItemsRequest{AlbumId:album.Id, PageSize:100, PageToken: pageToken}

        items, err := photos.MediaItems.Search(r).Do()
        if err != nil {
          jAlbum.Error = true
          break
        }

        jAlbum.ItemsCount = jAlbum.ItemsCount + len(items.MediaItems)
        pageToken = items.NextPageToken
        if pageToken == "" {
          break
        }
      }

      buf, err := json.Marshal(jAlbum)
      if err != nil {
        log.Printf("\nALBUM MARSHAL ERROR: %v\n", err)
        continue
      }

      if _, err := albumsFile.Write(append(buf, ",\n"...)); err != nil {
        log.Fatal(err)
      }
    }

    nextPageToken = albums.NextPageToken
    if nextPageToken == "" {
      break
    }
  }

  log.Printf("\n=== Complite ===")
}
