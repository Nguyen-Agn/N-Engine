package napi

import "github.com/Nguyen-Agn/N-Engine/modules/enginetype"

// LogError prints an error message to the console only if the engine is compiled with the 'debug' tag.
// Usage: go run -tags debug main.go
func LogError(format string, a ...any) {
	enginetype.LogError(format, a...)
}
