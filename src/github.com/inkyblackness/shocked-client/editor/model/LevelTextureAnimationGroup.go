package model

import (
	"github.com/inkyblackness/shocked-model"
)

// LevelTextureAnimationGroup describes a group of textures.
type LevelTextureAnimationGroup struct {
	id int

	properties model.TextureAnimation
}

func newLevelTextureAnimationGroup(id int) *LevelTextureAnimationGroup {
	return &LevelTextureAnimationGroup{id: id}
}

// FrameTime is for one frame, in milliseconds.
func (group *LevelTextureAnimationGroup) FrameTime() int {
	return *group.properties.FrameTime
}

// FrameCount returns the amount of frames in the animation.
func (group *LevelTextureAnimationGroup) FrameCount() int {
	return *group.properties.FrameCount
}

// LoopType returns how the animation shall behave.
func (group *LevelTextureAnimationGroup) LoopType() int {
	return *group.properties.LoopType
}
