package shapefiles

import (
	"bytes"
	"fmt"
	"github.com/jonas-p/go-shp"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/geojson"
)

type DomainShapeReader struct {
	ZipFilePath   string
	ShapeFileName string
}

// GetAllShapes will read all shapes from a shape file.  It is assuming the shapes are multipolygons
func (dep DomainShapeReader) GetAllShapes() ([]domain.DomainShape, error) {

	shapeFile, err := shp.OpenShapeFromZip(dep.ZipFilePath, dep.ShapeFileName)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer shapeFile.Close()

	ret := make([]domain.DomainShape, 0)

	for shapeFile.Next() {
		_, shape := shapeFile.Shape()
		polygon, ok := shape.(*shp.Polygon)
		if !ok {
			return nil, fmt.Errorf("not a valid shape")
		}

		// Get the geojson polygon
		geoJsonPolygon := toGeoJsonMultipolygon(polygon)
		props := getShapeProperties(shapeFile)

		domainShape := domain.DomainShape{
			Geometry:   geoJsonPolygon,
			Properties: props,
		}

		ret = append(ret, domainShape)
	}

	return ret, nil
}

func getShapeProperties(shapeFile *shp.ZipReader) map[string]any {
	props := make(map[string]any)
	fields := shapeFile.Fields()
	for fieldIdx, field := range fields {
		fieldName := string(bytes.Split(field.Name[:], []byte{0})[0])
		value := shapeFile.Attribute(fieldIdx)
		props[fieldName] = value
	}
	return props
}

func toGeoJsonMultipolygon(polygon *shp.Polygon) geojson.MultiPolygon {
	coordinates := make([][][][]float64, 0)
	lastIdx := 0
	for _, partIdx := range polygon.Parts {
		if partIdx == 0 {
			continue
		}
		jsonPolygon := toGeoJsonPolygon(polygon.Points[lastIdx:int(partIdx)])
		coordinates = append(coordinates, jsonPolygon.Coordinates)
		lastIdx = int(partIdx)
	}

	jsonPolygon := toGeoJsonPolygon(polygon.Points[lastIdx:])
	coordinates = append(coordinates, jsonPolygon.Coordinates)

	multipolygon := geojson.MultiPolygon{Type: "MultiPolygon", Coordinates: coordinates}
	return multipolygon
}

func toGeoJsonPolygon(points []shp.Point) geojson.Polygon {
	geoPoints := make([][]float64, 0)
	for _, point := range points {
		geoCoordinate := toGeoJsonCoordinate(point)
		geoPoints = append(geoPoints, geoCoordinate)
	}
	geoPolyPoints := [][][]float64{geoPoints}
	geoPoly := geojson.Polygon{Type: "Polygon", Coordinates: geoPolyPoints}
	return geoPoly
}

func toGeoJsonCoordinate(point shp.Point) []float64 {
	return []float64{point.X, point.Y}
}
