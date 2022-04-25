package nwsshapefiles

import (
	"bytes"
	"github.com/jonas-p/go-shp"
	"weather-alerts-service/internal/domain"
	"weather-alerts-service/pkg/geojson"
)

type CountyShapeAdapter struct {
}

func (dep *CountyShapeAdapter) MapFromCurrentShape(shapeFile *shp.ZipReader, geometry *geojson.MultiPolygon) domain.County {

	county := domain.County{}
	county.Geometry = *geometry
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

	return county
}
