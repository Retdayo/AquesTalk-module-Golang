//go:build !linux || !cgo

package aquestalk

import "errors"

type AquesTalk struct{}

func Open(libPath, devKey string) (*AquesTalk, error) {
	return nil, errors.New("aquestalk: LoadAquesTalk is supported on linux with cgo")
}

func LoadAquesTalk(libPath, devKey string) (*AquesTalk, error) {
	return Open(libPath, devKey)
}

func (a *AquesTalk) Speak(text string, speechSpeed int) ([]byte, error) {
	return nil, errors.New("aquestalk: Synthe is supported on linux with cgo")
}

func (a *AquesTalk) Synthe(text string, speechSpeed int) ([]byte, error) {
	return a.Speak(text, speechSpeed)
}

func (a *AquesTalk) SpeakWithPlayback(text string, speechSpeed int, playbackSpeed int) ([]byte, error) {
	return nil, errors.New("aquestalk: SpeakWithPlayback is supported on linux with cgo")
}

func (a *AquesTalk) SyntheWithPlaybackSpeed(text string, speechSpeed int, playbackSpeed int) ([]byte, error) {
	return a.SpeakWithPlayback(text, speechSpeed, playbackSpeed)
}
