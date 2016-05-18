package convert

type Text struct {
	Entries []TextEntry `xml:"Entry"`
}

type TextEntry struct {
	Block *uint16 `xml:"block,attr,omitempty"`
	CData string  `xml:",cdata"`
}
