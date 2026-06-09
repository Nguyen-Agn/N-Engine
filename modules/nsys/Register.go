package nsys

import (
	"sync"

	"github.com/Nguyen-Agn/N-Engine/domain"
)

type IGlobalConfig = domain.IGlobalConfig

type globalStore struct {
	sprites       map[string]domain.ISpriteLW
	sounds        map[string]domain.IAudioLW
	globalObjects map[string]domain.IObject

	variable map[string]any
	constant map[string]any

	mu sync.RWMutex
}

var (
	store *globalStore
	//config *IGlobalConfig
	once sync.Once
)
