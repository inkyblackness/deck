package styling

import "fmt"

type nullStyle struct {
	nullFunc StyleFunc
}

// NullStyle returns a Style that does no styling at all
func NullStyle() Style {
	style := &nullStyle{nullFunc: fmt.Sprint}

	return style
}

func (style *nullStyle) Prompt() StyleFunc {
	return style.nullFunc
}

func (style *nullStyle) Error() StyleFunc {
	return style.nullFunc
}

func (style *nullStyle) Status() StyleFunc {
	return style.nullFunc
}

func (style *nullStyle) Added() StyleFunc {
	return style.nullFunc
}

func (style *nullStyle) Removed() StyleFunc {
	return style.nullFunc
}
