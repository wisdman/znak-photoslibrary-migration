package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wisdman/znak-photoslibrary-migration/lib/oauth"
	"github.com/wisdman/znak-photoslibrary-migration/lib/old"
	"github.com/wisdman/znak-photoslibrary-migration/lib/photos"
)

func main() {

	logFile, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	ctx := context.Background()
	client, err := oauth.New(ctx, []string{photos.OAuthScope})
	if err != nil {
		log.Fatal(err)
	}

	for {
		oldAlbum, err := old.NextAlbum()
		if err != nil {
			log.Fatal(err)
		}

		if oldAlbum == nil {
			break
		}

		log.Printf("=== Processing album %d ===", oldAlbum.ID())
		log.Printf("Title: %s", oldAlbum.Title())
		log.Printf("Date: %s", oldAlbum.Date())
		log.Printf("Author: %s", oldAlbum.Author())
		log.Printf("Location: %s", oldAlbum.Location())
		log.Printf("Images: %v", oldAlbum.ImagesOIDS())

		err = oldAlbum.PrepareImages()
		if err != nil {
			log.Fatal(err)
		}

		newAlbum, err := photos.NewAlbum(client, oldAlbum.Title())
		if err != nil {
			log.Fatal(err)
		}

		for i, oldImage := range oldAlbum.Images {
			log.Printf("Image %d (%d of %d)", oldImage.OID(), i+1, len(oldAlbum.Images))
			err = newAlbum.AddImage(oldImage.Filepath(), oldImage.Description)
			if err != nil {
				logStr := fmt.Sprintf("[%s] %s: %d\n", oldAlbum.ID(), oldAlbum.Title(), oldImage.OID())
				if _, err := logFile.Write([]byte(logStr)); err != nil {
					log.Fatal(err)
				}

				log.Printf("IMAGE ERROR: %s", logStr)
			}
		}

		err = oldAlbum.Remove()
		if err != nil {
			log.Fatal(err)
		}

		err = oldAlbum.MarkMigrated()
		if err != nil {
			log.Fatal(err)
		}

		old.Progress()
		log.Printf("=== POOL ===\n%v", old.DB.Stat())
	}

	log.Printf("=== Complite ===")
}
