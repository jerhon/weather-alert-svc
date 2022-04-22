package geojson

type FeatureCollection[T any] struct {
	Type     string       `json:"type"`
	Features []Feature[T] `json:"features"`
}

type Feature[T any] struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Geometry   *Geometry `json:"geometry"`
	Properties *T        `json:"properties"`
}
