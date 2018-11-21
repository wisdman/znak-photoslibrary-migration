package photos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Upload(client *http.Client, filepath string) (string, error) {
	log.Printf("Uploading %s", filepath)

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", APIbasePath, APIVersion), file)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Goog-Upload-File-Name", file.Name())

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
