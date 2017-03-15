package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseScenery = interpreters.New()

var displayScenery = baseScenery.
	With("FrameCount", 0, 2).
	With("LoopType", 2, 2).
	With("AlternationType", 4, 2).
	With("PictureSource", 6, 2).
	With("AlternateSource", 8, 2)

var displayControlPedestal = baseScenery.
	With("FrameCount", 0, 2).
	With("TriggerObjectIndex", 2, 2).
	With("AlternationType", 4, 2).
	With("PictureSource", 6, 2).
	With("AlternateSource", 8, 2)

var cabinetFurniture = baseScenery.
	With("Object1Index", 2, 2).
	With("Object2Index", 4, 2)

var texturableFurniture = baseScenery.
	With("TextureIndex", 6, 2)

var wordScenery = baseScenery.
	With("TextIndex", 0, 2).
	With("FontAndSize", 2, 1).
	With("Color", 4, 1)

var textureMapScenery = baseScenery.
	With("TextureIndex", 6, 2)

var buttonControlPedestal = baseScenery.
	With("TriggerObjectIndex", 2, 2)

var surgicalMachine = baseScenery.
	With("BrokenState", 2, 1).
	With("BrokenMessageIndex", 5, 1)

var securityCamera = baseScenery.
	With("PanningSwitch", 2, 1)

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
