package main

import (
  "context"
  "os"
  "log"
  "fmt"

  "github.com/wisdman/znak-photoslibrary-migration/lib/oauth"
  photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const OAuthScope = "https://www.googleapis.com/auth/photoslibrary"

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

      author, err := GetAuthor(album.Title)
      if err != nil {
        log.Printf("\nGetting author error: %v\n", err)
        continue
      }

      if author == "" {
        log.Printf("\nAuthor not found: %v\n", err)
        continue
      }

      request := &photoslibrary.AddEnrichmentToAlbumRequest{
        AlbumPosition: &photoslibrary.AlbumPosition{ Position: "FIRST_IN_ALBUM" },
        NewEnrichmentItem: &photoslibrary.NewEnrichmentItem{
          TextEnrichment: &photoslibrary.TextEnrichment{ Text: author },
        },
      }

      _, err = photos.Albums.AddEnrichment(album.Id, request).Do()
      if err != nil {
        log.Printf("\nALBUM AddEnrichment REQUEST ERROR: %v\n", err)
        continue
      }
    }

    nextPageToken = albums.NextPageToken
    if nextPageToken == "" {
      break
    }
  }

  log.Printf("\n=== Complite ===")
}
