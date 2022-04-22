package domain

import (
	"time"
	"weather-alerts-service/pkg/geojson"
)

type Alert struct {
	OriginId string `bson:"originid"`
	Type     string
	AreaDesc string

	SAME       []string
	UGC        []string
	References []string

	Sent        time.Time
	Effective   time.Time
	Onset       time.Time
	Expires     time.Time
	Ends        time.Time
	Status      string
	MessageType string
	Category    string
	Severity    string
	Certainty   string
	Urgency     string
	Event       string
	Sender      string
	SenderName  string
	Headline    string
	Description string
	Instruction string
	Response    string

	Geometry *geojson.MultiPolygon
}

func (alert *Alert) GetKey() string {
	return alert.OriginId
}
