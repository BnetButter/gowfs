package main;

// FeatureCollection represents a GeoJSON FeatureCollection
type GeoJSONFeatureCollection struct {
	Type     string    `json:"type"` // always "FeatureCollection"
	Features []GeoJSONFeature `json:"features"`
}

// Feature represents a single GeoJSON Feature
type GeoJSONFeature struct {
	Type       string                 `json:"type"` // always "Feature"
	ID         string                 `json:"id"`
	Geometry   GeoJSONGeometry               `json:"geometry"`
	Properties map[string]string `json:"properties"`
}

// Geometry is a generic GeoJSON geometry object
type GeoJSONGeometry struct {
	Type        string      `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

