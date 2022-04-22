package geojson_test

import (
	"encoding/json"
	"testing"
	"weather-alerts-service/internal/testutils"
	"weather-alerts-service/pkg/geojson"
	"weather-alerts-service/pkg/nwsclient"
)

func TestFeatureCollection_Unmarshall(t *testing.T) {
	testData, _ := testutils.ReadTestFile("active_alerts.json")
	featureCollection := new(geojson.FeatureCollection[nwsclient.AlertProperties])
	err := json.Unmarshal([]byte(testData), featureCollection)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(featureCollection.Features), featureCollection.Features)
}
