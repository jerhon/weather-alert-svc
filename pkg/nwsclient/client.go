package nwsclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"weather-alerts-service/pkg/geojson"
)

type NwsClient struct {
	BaseUrl           string
	ApplicationString string
}

// NewNwsClient Creates a new client to retrieve data from the NWS weather service.
// applicationString is required to identify your application. You should include an e-mail as requested by the NWS to identify your application.
// See https://www.weather.gov/documentation/services-web-api for more information on the service.
func NewNwsClient(applicationString string) *NwsClient {
	client := new(NwsClient)
	client.BaseUrl = "https://api.weather.gov"
	client.ApplicationString = applicationString
	return client
}

func (client *NwsClient) GetActiveAlerts(ifModifiedSince *string) (*AlertResult, error) {
	requestData := []byte("")
	requestBuffer := bytes.NewBuffer(requestData)
	request, err := http.NewRequest("GET", client.BaseUrl+"/alerts/active", requestBuffer)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Go/1.18 ("+client.ApplicationString+")")
	if ifModifiedSince != nil {
		request.Header.Add("If-Modified-Since", *ifModifiedSince)
	}
	httpClient := &http.Client{}
	response, _ := httpClient.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	lastModifier := response.Header.Get("Last-Modifier")
	featureCollection := new(geojson.FeatureCollection[AlertProperties])
	err = json.Unmarshal(body, featureCollection)
	if err != nil {
		return nil, err
	}
	result := &AlertResult{LastModified: lastModifier, Alerts: featureCollection}
	return result, nil
}
