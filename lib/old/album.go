package old

import (
	"log"

	"github.com/jackc/pgx"
)

type Album struct {
	id       uint64
	title    string
	date     string
	author   string
	location string
	Images   []*Image
}

func GetAlbum(id uint64) (*Album, error) {
	var title string
	var date string
	var author string
	var location string

	if err := DB.QueryRow(`
    SELECT name as title, to_char(date,'YYYY:MM:DD HH12:MI:SS') as date, author, place as location FROM album WHERE id = $1
  `, id).Scan(&title, &date, &author, &location); err != nil {
		return nil, err
	}
	var images []*Image
	album := &Album{id, title, date, author, location, images}

	if err := album.GetImages(); err != nil {
		return nil, err
	}

	return album, nil
}

func NextAlbum() (*Album, error) {
	var id uint64

	err := DB.QueryRow(`
    select id FROM album WHERE NOT migrated LIMIT 1
  `).Scan(&id)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return GetAlbum(id)
}

func (a *Album) ID() uint64 {
	return a.id
}

func (a *Album) Title() string {
	return a.title
}

func (a *Album) Date() string {
	return a.date
}

func (a *Album) Author() string {
	if a.author == "" {
		return "Znak"
	}
	return a.author
}

func (a *Album) Location() string {
	if a.location == "" {
		return "Екатеринбург"
	}
	return a.location
}

func (a *Album) ImagesOIDS() (images []uint32) {
	for _, image := range a.Images {
		images = append(images, image.OID())
	}
	return
}

func (a *Album) GetImages() error {
	rows, err := DB.Query(`
    SELECT id, coalesce(array_to_string(json_array_castext(tags) || json_array_castext(people),', ',' '),'') as description FROM images WHERE album = $1
  `, a.id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var images []*Image

	for rows.Next() {
		var oid uint32
		var Description string
		if err := rows.Scan(&oid, &Description); err != nil {
			return err
		}

		images = append(images, &Image{oid, a.date, a.author, Description})
	}

	a.Images = images
	return nil
}

func (a *Album) PrepareImages() error {
	for i, image := range a.Images {
		log.Printf("Image %d (%d of %d)", image.OID(), i+1, len(a.Images))

		if err := image.Download(); err != nil {
			return err
		}

		if err := image.UpdateEXIF(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Album) Remove() error {
	for _, image := range a.Images {
		if err := image.Remove(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Album) MarkMigrated() error {
	log.Printf("Marking album %d as Migrated", a.id)

	_, err := DB.Query(`
    UPDATE album SET migrated = TRUE WHERE id = $1
  `, a.id)

	if err != nil {
		return err
	}

	return nil
}
