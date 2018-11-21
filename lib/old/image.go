package old

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wisdman/znak-photoslibrary-migration/lib/exif"
)

type Image struct {
	oid         uint32
	date        string
	author      string
	Description string
}

func (i *Image) OID() uint32 {
	return i.oid
}

func (i *Image) Url() string {
	return fmt.Sprintf("%s/photo/photo.php?jpg=%d", APIbasePath, i.oid)
}

func (i *Image) Filename() string {
	return fmt.Sprintf("img-%d.jpg", i.oid)
}

func (i *Image) Filepath() string {
	return fmt.Sprintf("/tmp/%s", i.Filename())
}

func (i *Image) Download() error {
	log.Printf("Downloading %s", i.Filename())

	res, err := http.Get(i.Url())
	if err != nil {
		return err
	}

	defer res.Body.Close()

	file, err := os.Create(i.Filepath())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) UpdateEXIF() error {
	args := []string{
		"-overwrite_original",
		fmt.Sprintf("-AllDates=%s", i.date),
		fmt.Sprintf("-Author=%s", i.author),
		fmt.Sprintf("-Comment=%s", i.Description),
		i.Filepath(),
	}

	return exif.Exec(args)
}

func (i *Image) Remove() error {
	log.Printf("Removing %s", i.Filename())
	if err := os.Remove(i.Filepath()); err != nil {
		return err
	}

	return nil
}

func (i *Image) AllDates(date string) error {
	return nil
}
