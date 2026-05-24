package naudio

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

type IAudioSystem = domain.IAudioSystem
type AudioData = domain.AudioData
type IAudioLW = domain.IAudioLW

var (
	Audio = enginetype.Audio
)
