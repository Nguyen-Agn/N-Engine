//go:build debug

package enginetype

import "fmt"

// LogError prints an error message to the console only if the engine is compiled with the 'debug' tag.
func LogError(format string, a ...any) {
	fmt.Printf("[N-Engine Error] "+format+"\n", a...)
}
