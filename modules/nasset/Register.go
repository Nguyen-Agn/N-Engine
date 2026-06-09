package nasset

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/domain/bridge"
)

type ISpriteLoader = domain.ISpriteLoader
type IAudioLoader = domain.IAudioLoader
type IManifestLoader = domain.IManifestLoader

type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal

type SpriteLW = bridge.SpriteLW
