package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"aquestalk"
)

func readIntOrDefault(scanner *bufio.Scanner, prompt string, def int) int {
	fmt.Printf("%s [%d]> ", prompt, def)
	if !scanner.Scan() {
		return def
	}
	s := strings.TrimSpace(scanner.Text())
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// コンソールで選択
func selectVoice(voices []aquestalk.Voice) (aquestalk.Voice, error) {
	if len(voices) == 0 {
		return aquestalk.Voice{}, fmt.Errorf("no voices found")
	}

	fmt.Println("使用する声種を選択してください:")
	for i, v := range voices {
		fmt.Printf(" [%d] %s\n", i, v.Name)
	}

	for {
		fmt.Print("番号> ")
		var idx int
		_, err := fmt.Scan(&idx)
		if err == nil && idx >= 0 && idx < len(voices) {
			return voices[idx], nil
		}
		fmt.Println("無効な入力です")
	}
}

func main() {
	voices, err := aquestalk.ListVoices(AQUESTALK_BASE_DIR)
	if err != nil {
		panic(err)
	}

	voice, err := selectVoice(voices)
	if err != nil {
		panic(err)
	}

	fmt.Println("選択された声種:", voice.Name)

	aq, err := aquestalk.LoadAquesTalk(voice.Path, DEV_KEY)
	if err != nil {
		panic(err)
	}

	fmt.Println("AquesTalk Linux (Go)")
	fmt.Println("読み上げ速度 / 再生速度（音程） 分離版")
	fmt.Println("exit で終了")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("テキスト入力> ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if strings.EqualFold(text, "exit") {
			break
		}

		speechSpeed := readIntOrDefault(scanner, "読み上げ速度", DEFAULT_SPEECH_SPEED)
		playbackSpeed := readIntOrDefault(scanner, "再生速度(音程)", DEFAULT_PLAYBACK_SPEED)

		normalized := aquestalk.Normalize(text)
		fmt.Println("変換後テキスト:", normalized)

		wav, err := aq.Synthe(normalized, speechSpeed)
		if err != nil {
			fmt.Println("音声生成失敗:", err)
			continue
		}

		wav, err = aquestalk.ApplyPlaybackSpeed(wav, playbackSpeed)
		if err != nil {
			fmt.Println("音程変更失敗:", err)
			continue
		}

		if err := os.WriteFile("out.wav", wav, 0644); err != nil {
			fmt.Println("保存失敗:", err)
			continue
		}

		fmt.Println("WAV生成完了: out.wav")
	}
}
