package main

import (
	"context"
	"log"

	"github.com/wisdman/znak-photoslibrary-migration/lib/oauth"
	"github.com/wisdman/znak-photoslibrary-migration/lib/old"
	"github.com/wisdman/znak-photoslibrary-migration/lib/photos"
)

func main() {

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
				log.Fatal(err)
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
	}

	log.Printf("=== Complite ===")
}
