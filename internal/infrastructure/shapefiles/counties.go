package shapefiles

import (
	"fmt"
	"weather-alerts-service/internal/domain"
)

func MapDomainShapeToCounty(shape domain.DomainShape) domain.County {

	county := domain.County{}
	county.Geometry = shape.Geometry
	for fieldName, value := range shape.Properties {
		switch fieldName {
		case "STATE":
			county.State = fmt.Sprintf("%v", value)
		case "FIPS":
			county.CountyFips = fmt.Sprintf("%v", value)
		case "COUNTYNAME":
			county.CountyName = fmt.Sprintf("%v", value)
		case "TIME_ZONE":
			county.TimeZone = fmt.Sprintf("%v", value)
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
