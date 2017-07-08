package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseScenery = interpreters.New()

var displayScenery = baseScenery.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 4)).
	With("LoopType", 2, 2).As(interpreters.EnumValue(map[uint32]string{0: "Forward", 1: "Forward/Backward", 2: "Backward", 3: "Forward/Backward"})).
	With("AlternationType", 4, 2).As(interpreters.EnumValue(map[uint32]string{0: "Don't Alternate", 3: "Alternate Randomly"})).
	With("PictureSource", 6, 2).As(interpreters.RangedValue(0, 0x01FF)).
	With("AlternateSource", 8, 2).As(interpreters.RangedValue(0, 0x01FF))

var displayControlPedestal = baseScenery.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 4)).
	With("TriggerObjectIndex", 2, 2).As(interpreters.ObjectIndex()).
	With("AlternationType", 4, 2).As(interpreters.EnumValue(map[uint32]string{0: "Don't Alternate", 3: "Alternate Randomly"})).
	With("PictureSource", 6, 2).As(interpreters.RangedValue(0, 0x01FF)).
	With("AlternateSource", 8, 2).As(interpreters.RangedValue(0, 0x01FF))

var cabinetFurniture = baseScenery.
	With("Object1Index", 2, 2).As(interpreters.ObjectIndex()).
	With("Object2Index", 4, 2).As(interpreters.ObjectIndex())

var texturableFurniture = baseScenery.
	With("TextureIndex", 6, 2).As(interpreters.RangedValue(0, 500))

var wordScenery = baseScenery.
	With("TextIndex", 0, 2).As(interpreters.RangedValue(0, 511)).
	With("Font", 2, 1).As(interpreters.Bitfield(map[uint32]string{
	0x0F: "Face",
	0xF0: "Size"})).
	With("Color", 4, 1).As(interpreters.RangedValue(0, 255))

var textureMapScenery = baseScenery.
	With("TextureIndex", 6, 2).As(interpreters.SpecialValue("LevelTexture"))

var buttonControlPedestal = baseScenery.
	With("TriggerObjectIndex", 2, 2).As(interpreters.ObjectIndex())

var surgicalMachine = baseScenery.
	With("BrokenState", 2, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "OK", 0xE7: "Broken"})).
	With("BrokenMessageIndex", 5, 1)

var securityCamera = baseScenery.
	With("PanningSwitch", 2, 1).As(interpreters.EnumValue(map[uint32]string{0: "Stationary", 1: "Panning"}))

var solidBridge = baseScenery.
	With("Size", 2, 1).
	With("Height", 3, 1).
	With("TopBottomTexture", 4, 1).
	With("SideTexture", 5, 1)

var forceBridge = baseScenery.
	With("Size", 2, 1).
	With("Height", 3, 1).
	With("Color", 6, 1)

func initScenery() interpreterRetriever {
	displays := newInterpreterLeaf(displayScenery)
	textureable := newInterpreterLeaf(texturableFurniture)

	electronics := newInterpreterEntry(baseScenery)
	electronics.set(6, displays)
	electronics.set(7, displays)

	furniture := newInterpreterEntry(baseScenery)
	furniture.set(2, newInterpreterLeaf(cabinetFurniture))
	furniture.set(5, textureable)
	furniture.set(7, textureable)
	furniture.set(8, textureable)

	surfaces := newInterpreterEntry(baseScenery)
	surfaces.set(3, newInterpreterLeaf(wordScenery))
	surfaces.set(6, displays)
	surfaces.set(7, newInterpreterLeaf(textureMapScenery))
	surfaces.set(8, displays)
	surfaces.set(9, displays)

	lighting := newInterpreterEntry(baseScenery)

	medicalEquipment := newInterpreterEntry(baseScenery)
	medicalEquipment.set(0, newInterpreterLeaf(buttonControlPedestal))
	medicalEquipment.set(3, newInterpreterLeaf(surgicalMachine))
	medicalEquipment.set(5, textureable)

	scienceSecurityEquipment := newInterpreterEntry(baseScenery)
	scienceSecurityEquipment.set(4, newInterpreterLeaf(securityCamera))
	scienceSecurityEquipment.set(6, newInterpreterLeaf(displayControlPedestal))

	gardenScenery := newInterpreterEntry(baseScenery)

	bridges := newInterpreterEntry(baseScenery)
	solidBridges := newInterpreterLeaf(solidBridge)
	forceBridges := newInterpreterLeaf(forceBridge)
	bridges.set(0, solidBridges)
	bridges.set(1, solidBridges)
	bridges.set(7, forceBridges)
	bridges.set(8, forceBridges)
	bridges.set(9, forceBridges)

	class := newInterpreterEntry(baseScenery)
	class.set(0, electronics)
	class.set(1, furniture)
	class.set(2, surfaces)
	class.set(3, lighting)
	class.set(4, medicalEquipment)
	class.set(5, scienceSecurityEquipment)
	class.set(6, gardenScenery)
	class.set(7, bridges)

	return class
}
