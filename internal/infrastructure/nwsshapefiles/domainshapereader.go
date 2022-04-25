package nwsshapefiles

import (
	"fmt"
	"github.com/jonas-p/go-shp"
	"weather-alerts-service/pkg/geojson"
)

type DomainShapeReader[T any] struct {
	ZipFilePath   string
	ShapeFileName string

	ShapeAdapter DomainShapeMapper[T]
}

type DomainShapeMapper[T any] interface {
	MapFromCurrentShape(shapeFile *shp.ZipReader, polygon *geojson.MultiPolygon) T
}

// GetAllCounties will read all counties from the shape file and return them as a slice of counties
func (dep *DomainShapeReader[T]) GetAll() ([]T, error) {

	shapeFile, err := shp.OpenShapeFromZip(dep.ZipFilePath, dep.ShapeFileName)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer shapeFile.Close()

	ret := make([]T, 0)

	for shapeFile.Next() {
		_, shape := shapeFile.Shape()
		polygon, ok := shape.(*shp.Polygon)
		if !ok {
			return nil, fmt.Errorf("not a valid shape")
		}

		// Get the geojson polygon
		geoJsonPolygon := toGeoJsonMultipolygon(polygon)
		domainType := dep.ShapeAdapter.MapFromCurrentShape(shapeFile, &geoJsonPolygon)

		ret = append(ret, domainType)
	}

	return ret, nil
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
