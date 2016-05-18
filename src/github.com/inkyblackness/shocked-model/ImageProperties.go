package model

// ImageProperties contain extra information about images - if applicable.
type ImageProperties struct {
	HotspotLeft   int `json:"hotspotLeft"`
	HotspotTop    int `json:"hotspotTop"`
	HotspotRight  int `json:"hotspotRight"`
	HotspotBottom int `json:"hotspotBottom"`
}
