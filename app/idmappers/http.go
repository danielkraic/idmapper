package idmappers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/sirupsen/logrus"
)

// NewHTTPIDMapper creates IDMapper that reads data from http
func NewHTTPIDMapper(log *logrus.Logger, url string, timeout time.Duration) (*idmapper.IDMapper, error) {
	if url == "" {
		return nil, fmt.Errorf("Empty url set for HTTP IDMapper")
	}
	return idmapper.NewIDMapper(&httpSource{
		url: url,
		log: log,
	})
}

type httpSource struct {
	log     *logrus.Logger
	url     string
	timeout time.Duration
}

type httpSourceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (source httpSource) Read() (idmapper.ValuesMap, error) {
	result := make(idmapper.ValuesMap)

	client := http.Client{
		Timeout:       source.timeout,
		CheckRedirect: nil,
	}

	httpResponse, err := client.Get(source.url)
	if err != nil {
		return result, fmt.Errorf("failed to get url %s: %s", source.url, err)
	}
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			source.log.Errorf("failed to close http body: %s", err)
		}
	}()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response from url %s: %s", source.url, err)
	}

	var responseData []httpSourceResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return result, fmt.Errorf("failed to decode json from url %s: %s", source.url, err)
	}

	for _, item := range responseData {
		result[item.ID] = item.Name
	}

	return result, nil
}
