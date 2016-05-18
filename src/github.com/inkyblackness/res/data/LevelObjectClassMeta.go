package data

import "github.com/inkyblackness/res"

// LevelObjectClassMeta contains meta information about a level object class
type LevelObjectClassMeta struct {
	EntrySize  int
	EntryCount int
}

var levelObjectClassMetaList = []LevelObjectClassMeta{
	{LevelWeaponEntrySize, LevelWeaponEntryCount},
	{LevelAmmoEntrySize, LevelAmmoEntryCount},
	{LevelProjectileEntrySize, LevelProjectileEntryCount},
	{LevelExplosiveEntrySize, LevelExplosiveEntryCount},
	{LevelPatchEntrySize, LevelPatchEntryCount},
	{LevelHardwareEntrySize, LevelHardwareEntryCount},
	{LevelSoftwareEntrySize, LevelSoftwareEntryCount},
	{LevelSceneryEntrySize, LevelSceneryEntryCount},
	{LevelItemEntrySize, LevelItemEntryCount},
	{LevelPanelEntrySize, LevelPanelEntryCount},
	{LevelBarrierEntrySize, LevelBarrierEntryCount},
	{LevelAnimationEntrySize, LevelAnimationEntryCount},
	{LevelMarkerEntrySize, LevelMarkerEntryCount},
	{LevelContainerEntrySize, LevelContainerEntryCount},
	{LevelCritterEntrySize, LevelCritterEntryCount}}

// LevelObjectClassMetaEntry returns the meta entry for the corresponding object class.
func LevelObjectClassMetaEntry(class res.ObjectClass) LevelObjectClassMeta {
	return levelObjectClassMetaList[int(class)]
}
