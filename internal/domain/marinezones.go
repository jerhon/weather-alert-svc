package domain

import "weather-alerts-service/pkg/geojson"

type MarineZone struct {
	ID     string
	WFO    string
	GL_WFO string
	Name   string

	Geometry geojson.MultiPolygon
}
