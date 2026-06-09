package naudio

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
)

type IAudioSystem = domain.IAudioSystem
type AudioData = domain.AudioData
type IAudioLW = domain.IAudioLW

var (
	Audio = enginetype.Audio
)
