package modes

var classNames = []string{
	"Weapons 0",
	"AmmoClips 1",
	"Projectiles 2",
	"Explosives 3",
	"Patches 4",
	"Hardware 5",
	"Software 6",
	"Scenery 7",
	"Items 8",
	"Panels 9",
	"Barriers 10",
	"Animations 11",
	"Markers 12",
	"Containers 13",
	"Critters 14"}

var maxObjectsPerClass = []int{16, 32, 32, 32, 32, 8, 16, 176, 128, 64, 64, 32, 160, 64, 64}

type objectClassItem struct {
	class int
}

func (item *objectClassItem) String() string {
	return classNames[item.class]
}
