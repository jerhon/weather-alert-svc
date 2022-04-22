package geojson

import (
	"encoding/json"
	"fmt"
)

type Geometry struct {
	Type            string           `json:"type"`
	Point           *Point           `json:"-"`
	Polygon         *Polygon         `json:"-"`
	LineString      *LineString      `json:"-"`
	MultiPoint      *MultiPoint      `json:"-"`
	MultiPolygon    *MultiPolygon    `json:"-"`
	MultiLineString *MultiLineString `json:"-"`
}

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Polygon struct {
	Type        string
	Coordinates [][][]float64
}

type LineString struct {
	Type        string
	Coordinates [][]float64
}

type MultiPoint struct {
	Type        string
	Coordinates [][]float64
}

type MultiPolygon struct {
	Type        string
	Coordinates [][][][]float64
}

type MultiLineString struct {
	Type        string
	Coordinates [][][]float64
}

func (geo *Geometry) GetShape() interface{} {
	if geo.Point != nil {
		return geo.Point
	} else if geo.MultiPoint != nil {
		return geo.MultiPoint
	} else if geo.LineString != nil {
		return geo.LineString
	} else if geo.Polygon != nil {
		return geo.Polygon
	} else if geo.MultiLineString != nil {
		return geo.MultiLineString
	} else if geo.MultiPolygon != nil {
		return geo.MultiPolygon
	}
	return nil
}

func (point *Point) GetRawCoordinates() interface{} {
	return point.Coordinates
}

func (geo *Geometry) GetMultiPolygon() *MultiPolygon {
	if geo.Polygon != nil {
		return &MultiPolygon{
			Type:        "MultiPolygon",
			Coordinates: [][][][]float64{geo.Polygon.Coordinates},
		}
	}
	if geo.MultiPolygon != nil {
		return geo.MultiPolygon
	}
	return nil
}

func (geo *Geometry) UnmarshalJSON(data []byte) (err error) {
	var raw map[string]json.RawMessage
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	var geoType string
	err = json.Unmarshal(raw["type"], &geoType)
	if err != nil {
		return err
	}
	geo.Type = geoType
	switch geoType {
	case "Point":
		point := new(Point)
		err = json.Unmarshal(data, &point)
		geo.Point = point
	case "Polygon":
		polygon := new(Polygon)
		err = json.Unmarshal(data, &polygon)
		geo.Polygon = polygon
	case "LineString":
		lineString := new(LineString)
		err = json.Unmarshal(data, &lineString)
		geo.LineString = lineString
	case "MultiPoint":
		multiPoint := new(MultiPoint)
		err = json.Unmarshal(data, &multiPoint)
		geo.MultiPoint = multiPoint
	case "MultiPolygon":
		multipolygon := new(MultiPolygon)
		err = json.Unmarshal(data, &multipolygon)
		geo.MultiPolygon = multipolygon
	case "MultiLineString":
		multiLineString := new(MultiLineString)
		err = json.Unmarshal(data, &multiLineString)
		geo.MultiLineString = multiLineString
	default:
		return fmt.Errorf("unknown geojson type = %s", geoType)
	}
	return nil
}

func (geo *Geometry) MarshalJSON() ([]byte, error) {
	var geoProperties map[string]interface{}
	geoProperties["type"] = geo.Type
	shape := geo.GetShape()
	marshalledJson, err := json.Marshal(shape)
	return marshalledJson, err
}
