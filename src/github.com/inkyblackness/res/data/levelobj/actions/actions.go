package actions

import (
	"github.com/inkyblackness/res/data/interpreters"
)

func forType(typeID int) func(*interpreters.Instance) bool {
	return func(inst *interpreters.Instance) bool {
		return inst.Get("Type") == uint32(typeID)
	}
}

var transportHackerDetails = interpreters.New().
	With("TargetX", 0, 4).
	With("TargetY", 4, 4)

var changeHealthDetails = interpreters.New().
	With("HealthDelta", 4, 2).
	With("HealthChangeFlag", 6, 2).
	With("PowerDelta", 8, 2).
	With("PowerChangeFlag", 10, 2)

var cloneMoveObjectDetails = interpreters.New().
	With("ObjectIndex", 0, 2).
	With("MoveFlag", 2, 2).
	With("TargetX", 4, 4).
	With("TargetY", 8, 4).
	With("TargetHeight", 12, 4)

var setGameVariableDetails = interpreters.New().
	With("VariableKey", 0, 4).
	With("Value", 4, 2).
	With("Operation", 6, 2).
	With("Message1", 8, 4).
	With("Message2", 12, 4)

var showCutsceneDetails = interpreters.New().
	With("CutsceneIndex", 0, 4).
	With("EndGameFlag", 4, 4)

var triggerOtherObjectsDetails = interpreters.New().
	With("Object1Index", 0, 2).
	With("Object1Delay", 2, 2).
	With("Object2Index", 4, 2).
	With("Object2Delay", 6, 2).
	With("Object3Index", 8, 2).
	With("Object3Delay", 10, 2).
	With("Object4Index", 12, 2).
	With("Object4Delay", 14, 2)

var changeLightingDetails = interpreters.New().
	With("ReferenceObjectIndex", 2, 2).
	With("TransitionType", 4, 2).
	With("LightSurface", 10, 2)

var effectDetails = interpreters.New().
	With("SoundIndex", 0, 2).
	With("SoundPlayCount", 2, 2).
	With("VisualEffect", 4, 2).
	With("AdditionalVisualEffect", 8, 2)

var changeTileHeightsDetails = interpreters.New().
	With("TileX", 0, 4).
	With("TileY", 4, 4).
	With("TargetFloorHeight", 8, 2).
	With("TargetCeilingHeight", 10, 2)

var randomTimerDetails = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("TimeLimit", 4, 4).
	With("ActivationValue", 8, 4)

var cycleObjectsDetails = interpreters.New().
	With("ObjectIndex1", 0, 4).
	With("ObjectIndex2", 4, 4).
	With("ObjectIndex3", 8, 4).
	With("NextObject", 12, 4)

var deleteObjectsDetails = interpreters.New().
	With("ObjectIndex1", 0, 2).
	With("ObjectIndex2", 4, 2).
	With("ObjectIndex3", 8, 2).
	With("MessageIndex", 12, 2)

var receiveEmailDetails = interpreters.New().
	With("EmailIndex", 0, 2)

var changeEffectDetails = interpreters.New().
	With("DeltaValue", 0, 2).
	With("EffectChangeFlag", 2, 2).
	With("EffectType", 4, 4)

var setObjectParameterDetails = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("Value1", 4, 4).
	With("Value2", 8, 4).
	With("Value3", 12, 4)

var setScreenPictureDetails = interpreters.New().
	With("ScreenObjectIndex1", 0, 2).
	With("ScreenObjectIndex2", 2, 2).
	With("SingleSequenceSource", 4, 4).
	With("LoopSequenceSource", 8, 4)

var trapMessageDetails = interpreters.New().
	With("BackgroundImageIndex", 0, 4).
	With("MessageIndex", 4, 4).
	With("TextColor", 8, 4).
	With("MfdSuppressionFlag", 12, 4)

var spawnObjectsDetails = interpreters.New().
	With("ObjectClass", 2, 1).
	With("ObjectSubclass", 1, 1).
	With("ObjectType", 0, 1).
	With("ReferenceObject1Index", 4, 2).
	With("ReferenceObject2Index", 6, 2).
	With("NumberOfObjects", 8, 4)

var changeObjectTypeDetails = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("NewType", 4, 2)

// Change state block

var toggleRepulsorChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("OffTextureIndex", 4, 1).
	With("OnTextureIndex", 5, 1)

var showGameCodeDigitChange = interpreters.New().
	With("ScreenObjectIndex", 0, 4).
	With("DigitNumber", 4, 4)

var setParameterFromVariableChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("ParameterNumber", 4, 4).
	With("VariableIndex", 8, 4)

var setButtonStateChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("NewState", 4, 4)

var doorControlChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("ControlValue", 4, 4)

var setConditionChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("Condition", 4, 4)

var makeItemRadioactiveChange = interpreters.New().
	With("ObjectIndex", 0, 4).
	With("WatchedObjectIndex", 4, 2).
	With("WatchedObjectTriggerState", 6, 2)

var orientedTriggerObjectChange = interpreters.New().
	With("HorizontalDirection", 0, 2).
	With("ObjectIndex", 4, 2)

var closeDataMfdChange = interpreters.New().
	With("ObjectIndex", 0, 4)

var changeStateDetails = interpreters.New().
	With("Type", 0, 4).
	Refining("ToggleRepulsor", 4, 12, toggleRepulsorChange, forType(1)).
	Refining("ShowGameCodeDigit", 4, 12, showGameCodeDigitChange, forType(2)).
	Refining("SetParameterFromVariable", 4, 12, setParameterFromVariableChange, forType(3)).
	Refining("SetButtonState", 4, 12, setButtonStateChange, forType(4)).
	Refining("DoorControl", 4, 12, doorControlChange, forType(5)).
	Refining("ReturnToMenu", 4, 12, interpreters.New(), forType(6)).
	// 7: undefined
	// 8: undefined
	Refining("ShodanPixelation", 4, 12, interpreters.New(), forType(9)).
	Refining("SetCondition", 4, 12, setConditionChange, forType(10)).
	Refining("ShowSystemAnalyzer", 4, 12, interpreters.New(), forType(11)).
	Refining("MakeItemRadioactive", 4, 12, makeItemRadioactiveChange, forType(12)).
	Refining("OrientedTriggerObject", 4, 12, orientedTriggerObjectChange, forType(13)).
	Refining("CloseDataMfd", 4, 12, closeDataMfdChange, forType(14)).
	Refining("EarthDestructionByLaser", 4, 12, interpreters.New(), forType(15))

var unconditionalAction = interpreters.New().
	With("Type", 0, 1).
	With("UsageQuota", 1, 1).
	Refining("TransportHacker", 6, 16, transportHackerDetails, forType(1)).
	Refining("ChangeHealth", 6, 16, changeHealthDetails, forType(2)).
	Refining("CloneMoveObject", 6, 16, cloneMoveObjectDetails, forType(3)).
	Refining("SetGameVariable", 6, 16, setGameVariableDetails, forType(4)).
	Refining("ShowCutscene", 6, 16, showCutsceneDetails, forType(5)).
	Refining("TriggerOtherObjects", 6, 16, triggerOtherObjectsDetails, forType(6)).
	Refining("ChangeLighting", 6, 16, changeLightingDetails, forType(7)).
	Refining("Effect", 6, 16, effectDetails, forType(8)).
	Refining("ChangeTileHeights", 6, 16, changeTileHeightsDetails, forType(9)).
	// 10 unknown
	Refining("RandomTimer", 6, 16, randomTimerDetails, forType(11)).
	Refining("CycleObjects", 6, 16, cycleObjectsDetails, forType(12)).
	Refining("DeleteObjects", 6, 16, deleteObjectsDetails, forType(13)).
	// 14 unknown
	Refining("ReceiveEmail", 6, 16, receiveEmailDetails, forType(15)).
	Refining("ChangeEffect", 6, 16, changeEffectDetails, forType(16)).
	Refining("SetObjectParameter", 6, 16, setObjectParameterDetails, forType(17)).
	Refining("SetScreenPicture", 6, 16, setScreenPictureDetails, forType(18)).
	Refining("ChangeState", 6, 16, changeStateDetails, forType(19)).
	// 20, 21 unknown
	Refining("TrapMessage", 6, 16, trapMessageDetails, forType(22)).
	Refining("SpawnObjects", 6, 16, spawnObjectsDetails, forType(23)).
	Refining("ChangeObjectType", 6, 16, changeObjectTypeDetails, forType(24))
