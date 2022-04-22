package domain

import "weather-alerts-service/pkg/geojson"

/*


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
*/

type County struct {
	UGC        string
	SAME       string
	State      string
	CWA        string
	CountyName string
	CountyFips string
	TimeZone   string

	Geometry geojson.MultiPolygon
}
