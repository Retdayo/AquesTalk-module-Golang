package aquestalk

import (
	"os"
	"path/filepath"
	"sort"
)

type Voice struct {
	Name string
	Path string
}

func Voices(baseDir string) ([]Voice, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var voices []Voice
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		name := e.Name()
		libPath := filepath.Join(baseDir, name, "libAquesTalk.so")

		if _, err := os.Stat(libPath); err == nil {
			voices = append(voices, Voice{
				Name: name,
				Path: libPath,
			})
		}
	}

	sort.Slice(voices, func(i, j int) bool {
		return voices[i].Name < voices[j].Name
	})

	return voices, nil
}

func ListVoices(baseDir string) ([]Voice, error) {
	return Voices(baseDir)
}
