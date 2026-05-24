package nasset

import (
	"autoworld/domain"
	"autoworld/domain/bridge"
)

type ISpriteLoader = domain.ISpriteLoader
type IAudioLoader = domain.IAudioLoader
type IManifestLoader = domain.IManifestLoader

type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal

type SpriteLW = bridge.SpriteLW
