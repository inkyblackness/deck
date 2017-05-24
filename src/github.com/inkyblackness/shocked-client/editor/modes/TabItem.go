package modes

type tabPage interface {
	SetVisible(visible bool)
}

type tabItem struct {
	page        tabPage
	displayName string
}

func (item *tabItem) String() string {
	return item.displayName
}
