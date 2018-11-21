package photos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Album struct {
	client     *http.Client
	Id         string `json:"id"`
	Title      string `json:"title"`
	ProductUrl string `json:"productUrl"`
}

func NewAlbum(client *http.Client, title string) (*Album, error) {
	log.Printf("Creating album: \"%s\"", title)

	album := &Album{client, "", title, ""}

	jsonString := fmt.Sprintf(`{"album":{"title":"%s"}}`, album.Title)
	jsonBuffer := bytes.NewBufferString(jsonString)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/albums", APIbasePath, APIVersion), jsonBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := album.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBlob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBlob, album)
	if err != nil {
		return nil, err
	}

	return album, nil
}

func (a *Album) AddImage(filepath string, description string) error {
	token, err := Upload(a.client, filepath)
	if err != nil {
		return err
	}

	log.Printf("Adding media item to album")

	jsonString := fmt.Sprintf(`{
		"albumId":"%s",
		"newMediaItems":[{
			"description": "%s",
			"simpleMediaItem": {
				"uploadToken": "%s"
			}
		}]
	}`, a.Id, description, token)

	jsonBuffer := bytes.NewBufferString(jsonString)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/mediaItems:batchCreate", APIbasePath, APIVersion), jsonBuffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonBlob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	type Status struct {
		Message string `json:"message"`
	}

	type Item struct {
		Status *Status `json:"status"`
	}

	type Response struct {
		Items []Item `json:"newMediaItemResults"`
	}

	var items []Item
	response := &Response{items}

	err = json.Unmarshal(jsonBlob, response)
	if err != nil {
		return err
	}

	if len(response.Items) == 0 || response.Items[0].Status == nil || response.Items[0].Status.Message != "OK" {
		return errors.New(fmt.Sprintf("Bad response: %s", string(jsonBlob)))
	}

	return nil
}
