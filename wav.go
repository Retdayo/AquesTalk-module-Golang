package aquestalk

import (
	"encoding/binary"
	"errors"
)

func ApplyPlaybackSpeed(wav []byte, playbackSpeed int) ([]byte, error) {
	return ChangePlaybackSpeed(wav, playbackSpeed)
}

func ChangePlaybackSpeed(wav []byte, playbackSpeed int) ([]byte, error) {
	if len(wav) < 44 {
		return nil, errors.New("invalid wav: too short")
	}
	if playbackSpeed <= 0 {
		return nil, errors.New("invalid playbackSpeed")
	}

	ratio := float64(playbackSpeed) / 100.0

	sampleRate := binary.LittleEndian.Uint32(wav[24:28])
	byteRate := binary.LittleEndian.Uint32(wav[28:32])

	newSampleRate := uint32(float64(sampleRate) * ratio)
	newByteRate := uint32(float64(byteRate) * ratio)

	binary.LittleEndian.PutUint32(wav[24:28], newSampleRate)
	binary.LittleEndian.PutUint32(wav[28:32], newByteRate)

	return wav, nil
}
