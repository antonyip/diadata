package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// DownloadResource is a simple utility that downloads a resource
// from @url and stores it into @filepath.
func DownloadResource(filepath, url string) error {

	log.Printf("Downloading data")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// GetRequest performs a get request on @url and returns the response body
// as a slice of byte data.
func GetRequest(url string) ([]byte, error) {

	// Get url
	response, err := http.Get(url)

	// Check, whether the request was successful
	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	// Close response body after function
	defer response.Body.Close()

	// Check the status code for a 200 so we know we have received a
	// proper response.
	if response.StatusCode != 200 {
		log.Error("HTTP Response Error: ", response.StatusCode)
		return []byte{}, fmt.Errorf("HTTP Response Error %d\n", response.StatusCode)
	}

	// Read the response body
	XMLdata, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	return XMLdata, err
}

// PostRequest performs a POST request on @url and returns the response body
// as a slice of byte data.
func PostRequest(url string, body io.Reader) ([]byte, error) {

	// Get url
	response, err := http.Post(url, "", body)

	// Check, whether the request was successful
	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	// Close response body after function
	defer response.Body.Close()

	// Check the status code for a 200 so we know we have received a
	// proper response.
	if response.StatusCode != 200 {
		log.Error("HTTP Response Error: ", response.StatusCode)
		return []byte{}, fmt.Errorf("HTTP Response Error %d\n", response.StatusCode)
	}

	// Read the response body
	XMLdata, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	return XMLdata, err
}

// GetCoinPrice Gets the price in USD of coin through our API
func GetCoinPrice(coin string) (float64, error) {
	type Quotation struct {
		Symbol             string
		Name               string
		Price              float64
		PriceYesterday     *float64
		VolumeYesterdayUSD *float64
		Source             string
		Time               time.Time
		ITIN               string
	}
	url := "https://api.diadata.org/v1/quotation/" + coin
	data, err := GetRequest(url)
	if err != nil {
		return 0, err
	}
	Quot := Quotation{}
	err = json.Unmarshal(data, &Quot)
	if err != nil {
		return 0, err
	}
	return Quot.Price, nil
}
