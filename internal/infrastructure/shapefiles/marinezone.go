package shapefiles

import (
	"fmt"
	"weather-alerts-service/internal/domain"
)

/*

Field Name	Type	width,dec	Description
ID	character	6	Marine Zone Identifier
WFO	character	3	Assigned WFO (Office Identifier)
GL_WFO	character	3	Great lakes WFO responsible for Open Lake Forecastts
NAME	character	250	Name of Marine Zone (In the offshore zone file, this attribute is "Name")
AJOIN0	character	5	Not Used
AJOIN1	character	5	Not Used
LON	numeric	10,5	Longitude of Centroid [decimal degrees]
LAT	numeric	9,5	Latitude of Centroid [decimal degrees]

*/

func MapDomainShapeToMarineZone(shape domain.DomainShape) domain.MarineZone {

	id := fmt.Sprintf("%v", shape.Properties["ID"])
	wfo := fmt.Sprintf("%v", shape.Properties["WFO"])
	glWfo := fmt.Sprintf("%v", shape.Properties["GL_WFO"])
	name := fmt.Sprintf("%v", shape.Properties["NAME"])

	return domain.MarineZone{
		ID:       id,
		WFO:      wfo,
		GL_WFO:   glWfo,
		Name:     name,
		Geometry: shape.Geometry,
	}
}
