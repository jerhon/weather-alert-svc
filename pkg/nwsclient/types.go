package nwsclient

import "weather-alerts-service/pkg/geojson"

type AlertResult struct {
	LastModified string
	Alerts       *geojson.FeatureCollection[AlertProperties]
}

type Geocode struct {
	SAME []string `json:"SAME"`
	UGC  []string `json:"UGC"`
}

type AlertReference struct {
	Identifier string `json:"identifier"`
	Sender     string `json:"sender"`
	Sent       string `json:"sent"`
}

type AlertProperties struct {
	Id            string           `json:"id"`
	AreaDesc      string           `json:"areaDesc"`
	Geocode       Geocode          `json:"geocode"`
	AffectedZones []string         `json:"affectedZones"`
	References    []AlertReference `json:"references"`
	Sent          string           `json:"sent"`
	Effective     string           `json:"effective"`
	Onset         string           `json:"onset"`
	Expires       string           `json:"expires"`
	Ends          string           `json:"ends"`
	Status        string           `json:"status"`
	MessageType   string           `json:"messageType"`
	Category      string           `json:"category"`
	Severity      string           `json:"severity"`
	Certainty     string           `json:"certainty"`
	Urgency       string           `json:"urgency"`
	Event         string           `json:"event"`
	Sender        string           `json:"sender"`
	SenderName    string           `json:"senderName"`
	Headline      string           `json:"headline"`
	Description   string           `json:"description"`
	Instruction   string           `json:"instruction"`
	Response      string           `json:"response"`

	// TODO: Parameters
}

/*
       "parameters": {
           "AWIPSidentifier": [
               "SPSLCH"
           ],
           "WMOidentifier": [
               "WWUS84 KLCH 132345"
           ],
           "NWSheadline": [
               "Strong thunderstorms will impact portions of southeastern Rapides and Avoyelles Parishes through 715 PM CDT"
           ],
           "eventMotionDescription": [
               "2022-04-13T23:45:00-00:00...storm...275DEG...38KT...31.32,-91.96 31.05,-92.34"
           ],
           "maxWindGust": [
               "50 MPH"
           ],
           "maxHailSize": [
               "0.25"
           ],
           "BLOCKCHANNEL": [
               "EAS",
               "NWEM",
               "CMAS"
           ],
           "EAS-ORG": [
               "WXR"
           ]
       }
   }
*/
