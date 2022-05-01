package domain

import "weather-alerts-service/pkg/geojson"

type DomainShape struct {
	Geometry   geojson.MultiPolygon
	Properties map[string]any
}
