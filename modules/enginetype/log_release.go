//go:build !debug

package enginetype

// LogError is a no-op in release builds, fully optimized away by the Go compiler.
func LogError(format string, a ...any) {}
