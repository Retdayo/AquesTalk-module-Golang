package aquestalk

import (
	"encoding/binary"
	"errors"
)

// 再生速度変更（= 音程 / WAVヘッダ）
// playbackSpeed: 100=等速, 130=1.3倍（高音・短く）
func ApplyPlaybackSpeed(wav []byte, playbackSpeed int) ([]byte, error) {
	if len(wav) < 44 {
		return nil, errors.New("invalid wav: too short")
	}
	if playbackSpeed <= 0 {
		return nil, errors.New("invalid playbackSpeed")
	}

	ratio := float64(playbackSpeed) / 100.0

	// WAVヘッダ: 24..27 sampleRate, 28..31 byteRate
	sampleRate := binary.LittleEndian.Uint32(wav[24:28])
	byteRate := binary.LittleEndian.Uint32(wav[28:32])

	newSampleRate := uint32(float64(sampleRate) * ratio)
	newByteRate := uint32(float64(byteRate) * ratio)

	binary.LittleEndian.PutUint32(wav[24:28], newSampleRate)
	binary.LittleEndian.PutUint32(wav[28:32], newByteRate)

	return wav, nil
}
