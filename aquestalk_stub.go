//go:build !linux || !cgo

package aquestalk

import "errors"

// AquesTalk is only supported on Linux with cgo enabled.
// This stub exists so the module can be fetched and documented on other platforms.
type AquesTalk struct{}

func LoadAquesTalk(libPath, devKey string) (*AquesTalk, error) {
	return nil, errors.New("aquestalk: LoadAquesTalk is supported on linux with cgo")
}

func (a *AquesTalk) Synthe(text string, speechSpeed int) ([]byte, error) {
	return nil, errors.New("aquestalk: Synthe is supported on linux with cgo")
}
