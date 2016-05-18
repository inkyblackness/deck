package model

// Referable describes anything that can be addressed. It comes with a "href" property.
type Referable struct {
	// Href contains the path to reach the referred object.
	Href string `json:"href"`
}
