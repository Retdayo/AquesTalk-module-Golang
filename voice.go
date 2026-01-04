package aquestalk

import (
	"os"
	"path/filepath"
	"sort"
)

// 声種情報
type Voice struct {
	Name string // フォルダ名（f1, m1 など）
	Path string // libAquesTalk.so のフルパス
}

// lib64/ 以下から声種を列挙
func ListVoices(baseDir string) ([]Voice, error) {
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
