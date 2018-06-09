package convert

type Text struct {
	Entries []TextEntry `xml:"Entry"`
}

type TextEntry struct {
	Block *int   `xml:"block,attr,omitempty"`
	CData string `xml:",cdata"`
}
