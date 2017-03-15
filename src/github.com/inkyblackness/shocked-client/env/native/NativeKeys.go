package native

import (
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/inkyblackness/shocked-client/env/keys"
)

var keyMap = map[glfw.Key]keys.Key{

	glfw.KeyEnter:     keys.KeyEnter,
	glfw.KeyKPEnter:   keys.KeyEnter,
	glfw.KeyEscape:    keys.KeyEscape,
	glfw.KeyBackspace: keys.KeyBackspace,
	glfw.KeyTab:       keys.KeyTab,

	glfw.KeyDown:  keys.KeyDown,
	glfw.KeyLeft:  keys.KeyLeft,
	glfw.KeyRight: keys.KeyRight,
	glfw.KeyUp:    keys.KeyUp,

	glfw.KeyDelete:   keys.KeyDelete,
	glfw.KeyEnd:      keys.KeyEnd,
	glfw.KeyHome:     keys.KeyHome,
	glfw.KeyInsert:   keys.KeyInsert,
	glfw.KeyPageDown: keys.KeyPageDown,
	glfw.KeyPageUp:   keys.KeyPageUp,

	glfw.KeyLeftAlt:      keys.KeyAlt,
	glfw.KeyLeftControl:  keys.KeyControl,
	glfw.KeyLeftShift:    keys.KeyShift,
	glfw.KeyLeftSuper:    keys.KeySuper,
	glfw.KeyRightAlt:     keys.KeyAlt,
	glfw.KeyRightControl: keys.KeyControl,
	glfw.KeyRightShift:   keys.KeyShift,
	glfw.KeyRightSuper:   keys.KeySuper,

	glfw.KeyPause:       keys.KeyPause,
	glfw.KeyPrintScreen: keys.KeyPrintScreen,

	glfw.KeyCapsLock:   keys.KeyCapsLock,
	glfw.KeyScrollLock: keys.KeyScrollLock,

	glfw.KeyF1:  keys.KeyF1,
	glfw.KeyF10: keys.KeyF10,
	glfw.KeyF11: keys.KeyF11,
	glfw.KeyF12: keys.KeyF12,
	glfw.KeyF13: keys.KeyF13,
	glfw.KeyF14: keys.KeyF14,
	glfw.KeyF15: keys.KeyF15,
	glfw.KeyF16: keys.KeyF16,
	glfw.KeyF17: keys.KeyF17,
	glfw.KeyF18: keys.KeyF18,
	glfw.KeyF19: keys.KeyF19,
	glfw.KeyF2:  keys.KeyF2,
	glfw.KeyF20: keys.KeyF20,
	glfw.KeyF21: keys.KeyF21,
	glfw.KeyF22: keys.KeyF22,
	glfw.KeyF23: keys.KeyF23,
	glfw.KeyF24: keys.KeyF24,
	glfw.KeyF25: keys.KeyF25,
	glfw.KeyF3:  keys.KeyF3,
	glfw.KeyF4:  keys.KeyF4,
	glfw.KeyF5:  keys.KeyF5,
	glfw.KeyF6:  keys.KeyF6,
	glfw.KeyF7:  keys.KeyF7,
	glfw.KeyF8:  keys.KeyF8,
	glfw.KeyF9:  keys.KeyF9}
