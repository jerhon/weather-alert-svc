package geojson

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"testing"
	"weather-alerts-service/internal/testutils"
)

func TestGeometry_UnmarshalJSON(t *testing.T) {
	content, _ := testutils.ReadTestFile("point.json")
	var g Geometry
	err := json.Unmarshal([]byte(content), &g)

	if err != nil {
		t.Error("Unexpected error: ", err)
	}

	if g.Point == nil {
		t.Error("Expected a point.")
	}

	if len(g.Point.Coordinates) != 2 {
		t.Error("Expected two points.")
	}

	expected1 := 125.6
	expected2 := 10.1
	if g.Point.Coordinates[0] != expected1 {
		t.Errorf("Point 1 doesn't match.")
	}
	if g.Point.Coordinates[1] != expected2 {
		t.Errorf("Point 2 doesn't match.")
	}
}

func TestGeometry_MarshalJSON(t *testing.T) {
	// TODO should set this at a higher level in the application.
	decimal.MarshalJSONWithoutQuotes = true

	coordinates := []float64{1.234, 5.6789}
	p := Point{Type: "Point", Coordinates: coordinates}

	data, _ := json.Marshal(p)
	t.Log(string(data))
}
