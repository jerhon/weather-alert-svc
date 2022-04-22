package nwsshapefiles

import (
	"bytes"
	"fmt"
	"github.com/jonas-p/go-shp"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/geojson"
)

type CountyReader struct {
	ZipFilePath   string
	ShapeFileName string
}

// GetAllCounties will read all counties from the shape file and return them as a slice of counties
func (dep *CountyReader) GetAllCounties() ([]domain.County, error) {

	shapeFile, err := shp.OpenShapeFromZip(dep.ZipFilePath, dep.ShapeFileName)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer shapeFile.Close()

	ret := make([]domain.County, 0)

	for shapeFile.Next() {
		_, shape := shapeFile.Shape()
		polygon, ok := shape.(*shp.Polygon)
		if !ok {
			return nil, fmt.Errorf("not a valid shape")
		}

		county := domain.County{}

		// Get the geojson polygon
		geoJsonPolygon := toGeoJsonMultipolygon(polygon)
		county.Geometry = geoJsonPolygon
		setCountyAttributes(shapeFile, &county)

		ret = append(ret, county)
	}

	return ret, nil
}

func setCountyAttributes(shapeFile *shp.ZipReader, county *domain.County) {
	fields := shapeFile.Fields()
	for fieldIdx, field := range fields {
		fieldName := string(bytes.Split(field.Name[:], []byte{0})[0])
		value := shapeFile.Attribute(fieldIdx)
		switch fieldName {
		case "STATE":
			county.State = value
		case "FIPS":
			county.CountyFips = value
		case "COUNTYNAME":
			county.CountyName = value
		case "TIME_ZONE":
			county.TimeZone = value
		}

		// From this URL https://www.weather.gov/gis/Counties
		// STATE (fips code)
		// CWA (County Warning Area)
		// COUNTYNAME County Name
		// FIPS Fips code the countt
		// TIME_ZONE
		//  //Two letters appear for the nine (9) counties (10 records total) which are divided by a time zone boundary, which are located in the states of FL (Gulf), ID (Idaho), ND (McKenzie, Dunn, and Sioux), NE (Cherry), OR (Malheur), SD (Stanley), and TX (Culberson).
		// FE_AREA
		// LON Central
		// LAT Central point

	}

	// SAME and UGC are codes by defined by
	if county.CountyFips != "" {
		county.SAME = "0" + county.CountyFips
	}
	if county.CountyFips != "" && county.State != "" {
		county.UGC = county.State + "C" + county.CountyFips[3:]
	}
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

	return geojson.MultiPolygon{Type: "MultiPolygon", Coordinates: coordinates}
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
