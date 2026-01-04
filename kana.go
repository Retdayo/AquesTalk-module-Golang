package aquestalk

import (
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/ikawaha/kagome-dict/ipa"
)

var tk *tokenizer.Tokenizer

func init() {
	dict := ipa.Dict()
	t, err := tokenizer.New(dict)
	if err != nil {
		panic(err)
	}
	tk = t
}

// カタカナ → ひらがな
func kataToHira(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 'ァ' && r <= 'ヶ' {
			b.WriteRune(r - 0x60)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// 日本語 → ひらがな（Goライブラリのみ）
func Normalize(text string) string {
	tokens := tk.Tokenize(text)

	var out strings.Builder
	for _, token := range tokens {
		// EOS(DUMMY) は除外
		if token.Class == tokenizer.DUMMY {
			continue
		}

		feat := token.Features()

		// 記号は Surface をそのまま使う
		if len(feat) > 0 && feat[0] == "記号" {
			out.WriteString(token.Surface)
			continue
		}

		// features[7] = 読み（カタカナ）
		if len(feat) > 7 && feat[7] != "*" {
			out.WriteString(kataToHira(feat[7]))
		} else {
			out.WriteString(token.Surface)
		}
	}

	kana := out.String()

	// AquesTalk向け正規化
	kana = strings.ReplaceAll(kana, "ー", "う")
	kana = strings.ReplaceAll(kana, "！", "。")
	kana = strings.ReplaceAll(kana, "!", "。")
	kana = strings.ReplaceAll(kana, "？", "。")
	kana = strings.ReplaceAll(kana, "?", "。")

	return kana
}
